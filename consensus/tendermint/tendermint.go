package tendermint

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/tendermint/tendermint/version"

	"github.com/klyed/hivesmartchain/binary"
	"github.com/klyed/hivesmartchain/consensus/abci"
	"github.com/klyed/hivesmartchain/crypto"
	"github.com/klyed/hivesmartchain/genesis"
	"github.com/klyed/hivesmartchain/logging"
	"github.com/klyed/hivesmartchain/logging/structure"
	//"github.com/klyed/hivesmartchain/bridges"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/proxy"
	tmTypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

func init() {
	// Tendermint now sets this dynamically in it's build... we could also automate setting it
	version.TMCoreSemVer = "0.0.4"
}

// Serves as a wrapper around the Tendermint node's closeable resources (database connections)
type Node struct {
	*node.Node
	closers []interface {
		Close() error
	}
}

func DBProvider(ID string, backendType dbm.BackendType, dbDir string) (dbm.DB, error) {
	return dbm.NewDB(ID, backendType, dbDir)
}

// Since Tendermint doesn't close its DB connections
func (n *Node) DBProvider(ctx *node.DBContext) (dbm.DB, error) {
	db, err := DBProvider(ctx.ID, dbm.BackendType(ctx.Config.DBBackend), ctx.Config.DBDir())
	if err != nil {
		return nil, err
	}
	n.closers = append(n.closers, db)
	return db, nil
}

func (n *Node) Close() {
	for _, closer := range n.closers {
		closer.Close()
	}
}

func NewNode(conf *config.Config, privValidator tmTypes.PrivValidator, genesisDoc *tmTypes.GenesisDoc,
	app *abci.App, metricsProvider node.MetricsProvider, logger *logging.Logger) (*Node, error) {

	var err error
	// disable Tendermint's RPC
	conf.RPC.ListenAddress = ""

	nodeKey, err := EnsureNodeKey(conf.NodeKeyFile())
	if err != nil {
		return nil, err
	}

	nde := &Node{}
	nde.Node, err = node.NewNode(conf, privValidator,
		nodeKey, proxy.NewLocalClientCreator(app),
		func() (*tmTypes.GenesisDoc, error) {
			return genesisDoc, nil
		},
		nde.DBProvider,
		metricsProvider,
		NewLogger(logger.WithPrefix(structure.ComponentKey, structure.Tendermint).
			With(structure.ScopeKey, "tendermint.NewNode")))
	if err != nil {
		return nil, err
	}
	app.SetMempoolLocker(nde.Mempool())
	return nde, nil
}

func DeriveGenesisDoc(hscGenesisDoc *genesis.GenesisDoc, appHash []byte) *tmTypes.GenesisDoc {
	validators := make([]tmTypes.GenesisValidator, len(hscGenesisDoc.Validators))
	for i, validator := range hscGenesisDoc.Validators {
		validators[i] = tmTypes.GenesisValidator{
			Address: validator.PublicKey.TendermintAddress(),
			PubKey:  validator.PublicKey.TendermintPubKey(),
			Name:    validator.Name,
			Power:   int64(validator.Amount),
		}
	}
	consensusParams := tmTypes.DefaultConsensusParams()
	// This is the smallest increment we can use to get a strictly increasing sequence
	// of block time - we set it low to avoid skew
	// if the BlockTimeIota is longer than the average block time
	consensusParams.Block.TimeIotaMs = 1

	return &tmTypes.GenesisDoc{
		ChainID:         hscGenesisDoc.GetChainID(),
		GenesisTime:     hscGenesisDoc.GenesisTime,
		Validators:      validators,
		AppHash:         appHash,
		ConsensusParams: consensusParams,
		InitialHeight:   1,
	}
}

func NewNodeInfo(ni p2p.DefaultNodeInfo) *NodeInfo {
	address, _ := crypto.AddressFromHexString(string(ni.ID()))
	return &NodeInfo{
		ID:            address,
		Moniker:       ni.Moniker,
		ListenAddress: ni.ListenAddr,
		Version:       ni.Version,
		Channels:      binary.HexBytes(ni.Channels),
		Network:       ni.Network,
		RPCAddress:    ni.Other.RPCAddress,
		TxIndex:       ni.Other.TxIndex,
	}
}

func NewNodeKey() *p2p.NodeKey {
	privKey := ed25519.GenPrivKey()
	return &p2p.NodeKey{
		PrivKey: privKey,
	}
}

func WriteNodeKey(nodeKeyFile string, key json.RawMessage) error {
	err := os.MkdirAll(path.Dir(nodeKeyFile), 0777)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(nodeKeyFile, key, 0600)
}

func EnsureNodeKey(nodeKeyFile string) (*p2p.NodeKey, error) {
	err := os.MkdirAll(path.Dir(nodeKeyFile), 0777)
	if err != nil {
		return nil, err
	}

	return p2p.LoadOrGenNodeKey(nodeKeyFile)
}
