package forensics

import (
	"fmt"
	"testing"
	"time"

	"github.com/klyed/hivesmartchain/bcm"
	"github.com/klyed/hivesmartchain/event"
	"github.com/klyed/hivesmartchain/execution"
	"github.com/klyed/hivesmartchain/logging"
	"github.com/klyed/hivesmartchain/txs"

	"github.com/klyed/hivesmartchain/txs/payload"

	"github.com/klyed/hivesmartchain/consensus/tendermint"

	"github.com/klyed/hivesmartchain/acm"
	"github.com/klyed/hivesmartchain/execution/state"
	"github.com/klyed/hivesmartchain/genesis"
	"github.com/stretchr/testify/require"
	sm "github.com/klyed/tendermint/state"
	"github.com/klyed/tendermint/store"
	"github.com/klyed/tendermint/types"
	dbm "github.com/klyed/tm-db"
)

// This serves as a testbed for looking at non-deterministic hsc instances capture from the wild
// Put the path to 'good' and 'bad' hsc directories here (containing the config files and .hsc dir)

func TestStateComp(t *testing.T) {
	st1 := state.NewState(dbm.NewMemDB())
	_, _, err := st1.Update(func(ws state.Updatable) error {
		return ws.UpdateAccount(acm.NewAccountFromSecret("1"))
	})
	require.NoError(t, err)
	_, _, err = st1.Update(func(ws state.Updatable) error {
		return ws.UpdateAccount(acm.NewAccountFromSecret("2"))
	})
	require.NoError(t, err)

	db2 := dbm.NewMemDB()
	st2, err := st1.Copy(db2)
	require.NoError(t, err)
	err = CompareStateAtHeight(st2, st1, 0)
	require.Error(t, err)

	_, _, err = st2.Update(func(ws state.Updatable) error {
		return ws.UpdateAccount(acm.NewAccountFromSecret("3"))
	})
	require.NoError(t, err)

	err = CompareStateAtHeight(st2, st1, 1)
	require.Error(t, err)
}

func TestReplay(t *testing.T) {
	var height uint64 = 10
	genesisDoc, tmDB, hscDB := makeChain(t, height)

	src := NewSource(hscDB, tmDB, genesisDoc)
	dst := NewSourceFromGenesis(genesisDoc)
	re := NewReplay(src, dst)

	rc, err := re.Blocks(1, height)
	require.NoError(t, err)
	require.Len(t, rc, int(height-1))
}

func makeChain(t *testing.T, max uint64) (*genesis.GenesisDoc, dbm.DB, dbm.DB) {
	genesisDoc, _, validators := genesis.NewDeterministicGenesis(0).GenesisDoc(0, 1)

	tmDB := dbm.NewMemDB()
	bs := store.NewBlockStore(tmDB)
	gd := tendermint.DeriveGenesisDoc(genesisDoc, nil)
	st, err := sm.MakeGenesisState(&types.GenesisDoc{
		ChainID:    gd.ChainID,
		Validators: gd.Validators,
		AppHash:    gd.AppHash,
	})
	require.NoError(t, err)

	hscDB, hscState, hscChain, err := initHiveSmartChain(genesisDoc)
	require.NoError(t, err)

	committer, err := execution.NewBatchCommitter(hscState, execution.ParamsFromGenesis(genesisDoc),
		hscChain, event.NewEmitter(), logging.NewNoopLogger())
	require.NoError(t, err)

	var stateHash []byte
	for i := uint64(1); i < max; i++ {
		makeBlock(t, st, bs, func(block *types.Block) {

			decoder := txs.NewProtobufCodec()
			err = bcm.NewBlock(decoder, block).Transactions(func(txEnv *txs.Envelope) error {
				_, err := committer.Execute(txEnv)
				require.NoError(t, err)
				return nil
			})
			// empty if height == 1
			block.AppHash = stateHash
			// we need app hash in the abci header
			abciHeader := types.TM2PB.Header(&block.Header)
			stateHash, err = committer.Commit(&abciHeader)
			require.NoError(t, err)

		}, validators[0])
		require.Equal(t, int64(i), bs.Height())
	}
	return genesisDoc, tmDB, hscDB
}

func makeBlock(t *testing.T, st sm.State, bs *store.BlockStore, commit func(*types.Block), val *acm.PrivateAccount) {
	height := bs.Height() + 1
	tx := makeTx(t, st.ChainID, height, val)
	block, _ := st.MakeBlock(height, []types.Tx{tx}, new(types.Commit), nil,
		st.Validators.GetProposer().Address)

	commit(block)
	partSet := block.MakePartSet(2)
	commitSigs := []types.CommitSig{{Timestamp: time.Time{}}}
	seenCommit := types.NewCommit(height, 0, types.BlockID{
		Hash:          block.Hash(),
		PartSetHeader: partSet.Header(),
	}, commitSigs)
	bs.SaveBlock(block, partSet, seenCommit)
}

func makeTx(t *testing.T, chainID string, height int64, val *acm.PrivateAccount) (tx types.Tx) {
	sendTx := payload.NewSendTx()
	amount := uint64(height)
	acc := acm.NewAccountFromSecret(fmt.Sprintf("%d", height))
	sendTx.AddInputWithSequence(val.GetPublicKey(), amount, uint64(height))
	sendTx.AddOutput(acc.GetAddress(), amount)
	txEnv := txs.Enclose(chainID, sendTx)
	err := txEnv.Sign(val)
	require.NoError(t, err)

	data, err := txs.NewProtobufCodec().EncodeTx(txEnv)
	require.NoError(t, err)
	return types.Tx(data)
}
