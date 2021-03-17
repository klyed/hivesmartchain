package burrow

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/klye-dev/hivesmartchain/event"
	"github.com/klye-dev/hivesmartchain/event/query"
	"github.com/klye-dev/hivesmartchain/vent/chain"

	"github.com/klye-dev/hivesmartchain/binary"
	"github.com/klye-dev/hivesmartchain/crypto"
	"github.com/klye-dev/hivesmartchain/execution/errors"
	"github.com/klye-dev/hivesmartchain/execution/exec"
	"github.com/klye-dev/hivesmartchain/rpc/rpcevents"
	"github.com/klye-dev/hivesmartchain/rpc/rpcquery"
	"github.com/klye-dev/hivesmartchain/vent/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Chain struct {
	conn       *grpc.ClientConn
	filter     query.Query
	query      rpcquery.QueryClient
	exec       rpcevents.ExecutionEventsClient
	chainID    string
	version    string
	continuity exec.ContinuityOpt
}

var _ chain.Chain = (*Chain)(nil)

func New(conn *grpc.ClientConn, filter *chain.Filter) (*Chain, error) {
	client := rpcquery.NewQueryClient(conn)
	status, err := client.Status(context.Background(), &rpcquery.StatusParam{})
	if err != nil {
		return nil, fmt.Errorf("could not get initial status from Burrow: %w", err)
	}
	filterQuery, err := queryFromFilter(filter)
	if err != nil {
		return nil, fmt.Errorf("could not build Vent filter query: %w", err)
	}
	continuity := exec.Continuous
	if !query.IsEmpty(filterQuery) {
		// Since we may skip some events
		continuity = exec.NonConsecutiveEvents
	}
	return &Chain{
		conn:       conn,
		query:      client,
		filter:     filterQuery,
		exec:       rpcevents.NewExecutionEventsClient(conn),
		chainID:    status.ChainID,
		version:    status.HscVersion,
		continuity: continuity,
	}, nil
}

func (b *Chain) GetChainID() string {
	return b.chainID
}

func (b *Chain) GetVersion() string {
	return b.version
}

func (b *Chain) StatusMessage(ctx context.Context, lastProcessedHeight uint64) []interface{} {
	var catchUpRatio float64
	status, err := b.query.Status(ctx, &rpcquery.StatusParam{})
	if err != nil {
		err = fmt.Errorf("could not get Hive Smart Chain chain status: %w", err)
		return []interface{}{
			"msg", "status",
			"error", err.Error(),
		}
	}
	if status.SyncInfo.LatestBlockHeight > 0 {
		catchUpRatio = float64(lastProcessedHeight) / float64(status.SyncInfo.LatestBlockHeight)
	}
	return []interface{}{
		"msg", "status",
		"chain_type", "HSC",
		"last_processed_height", lastProcessedHeight,
		"fraction_caught_up", catchUpRatio,
		"hsc_latest_block_height", status.SyncInfo.LatestBlockHeight,
		"hsc_latest_block_duration", status.SyncInfo.LatestBlockDuration,
		"hsc_latest_block_hash", status.SyncInfo.LatestBlockHash,
		"hsc_latest_app_hash", status.SyncInfo.LatestAppHash,
		"hsc_latest_block_time", status.SyncInfo.LatestBlockTime,
		"hsc_latest_block_seen_time", status.SyncInfo.LatestBlockSeenTime,
		"hsc_node_info", status.NodeInfo,
		"hsc_catching_up", status.CatchingUp,
	}
}

func (b *Chain) ConsumeBlocks(ctx context.Context, in *rpcevents.BlockRange, consumer func(chain.Block) error) error {
	stream, err := b.exec.Stream(ctx, &rpcevents.BlocksRequest{
		BlockRange: in,
		Query:      b.filter.String(),
	})
	if err != nil {
		return fmt.Errorf("could not connect to block stream: %w", err)
	}

	return rpcevents.ConsumeBlockExecutions(stream, func(blockExecution *exec.BlockExecution) error {
		return consumer((*Block)(blockExecution))
	}, exec.Continuous)
}

func (b *Chain) Connectivity() connectivity.State {
	return b.conn.GetState()
}

func (b *Chain) GetABI(ctx context.Context, address crypto.Address) (string, error) {
	result, err := b.query.GetMetadata(ctx, &rpcquery.GetMetadataParam{
		Address: &address,
	})
	if err != nil {
		return "", err
	}
	return result.Metadata, nil
}

func (b *Chain) Close() error {
	return b.conn.Close()
}

type Block exec.BlockExecution

func NewBurrowBlock(block *exec.BlockExecution) *Block {
	return (*Block)(block)
}

func (b *Block) GetMetadata(columns types.SQLColumnNames) (map[string]interface{}, error) {
	blockHeader, err := json.Marshal(b.Header)
	if err != nil {
		return nil, fmt.Errorf("could not marshal block header: %w", err)
	}

	return map[string]interface{}{
		columns.Height:      strconv.FormatUint(b.Height, 10),
		columns.TimeStamp:   b.Header.GetTime(),
		columns.BlockHeader: string(blockHeader),
	}, nil
}

var _ chain.Block = (*Block)(nil)

func (b *Block) GetHeight() uint64 {
	return b.Height
}

func (b *Block) GetTxs() []chain.Transaction {
	txs := make([]chain.Transaction, len(b.TxExecutions))
	for i, tx := range b.TxExecutions {
		txs[i] = (*Transaction)(tx)
	}
	return txs
}

type Transaction exec.TxExecution

var _ chain.Transaction = (*Transaction)(nil)

func (tx *Transaction) GetOrigin() *chain.Origin {
	origin := (*exec.TxExecution)(tx).GetOrigin()
	if origin == nil {
		return nil
	}
	return &chain.Origin{
		ChainID: origin.ChainID,
		Height:  origin.Height,
		Index:   origin.Index,
	}
}

func (tx *Transaction) GetException() *errors.Exception {
	return tx.Exception
}

func (tx *Transaction) GetMetadata(columns types.SQLColumnNames) (map[string]interface{}, error) {
	// transaction raw data
	envelope, err := json.Marshal(tx.Envelope)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal envelope in tx %v: %v", tx, err)
	}

	events, err := json.Marshal(tx.Events)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal events in tx %v: %v", tx, err)
	}

	result, err := json.Marshal(tx.Result)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal result in tx %v: %v", tx, err)
	}

	receipt, err := json.Marshal(tx.Receipt)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal receipt in tx %v: %v", tx, err)
	}

	exception, err := json.Marshal(tx.Exception)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal exception in tx %v: %v", tx, err)
	}

	origin, err := json.Marshal(tx.Origin)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal origin in tx %v: %v", tx, err)
	}

	return map[string]interface{}{
		columns.Height:    tx.Height,
		columns.TxHash:    tx.TxHash.String(),
		columns.TxIndex:   tx.Index,
		columns.TxType:    tx.TxType.String(),
		columns.Envelope:  string(envelope),
		columns.Events:    string(events),
		columns.Result:    string(result),
		columns.Receipt:   string(receipt),
		columns.Origin:    string(origin),
		columns.Exception: string(exception),
	}, nil
}

func (tx *Transaction) GetHash() binary.HexBytes {
	return tx.TxHash
}

func (tx *Transaction) GetEvents() []chain.Event {
	// All txs have events, but not all have LogEvents
	var events []chain.Event
	for _, ev := range tx.Events {
		if ev.Log != nil {
			events = append(events, (*Event)(ev))
		}
	}
	return events
}

type Event exec.Event

var _ chain.Event = (*Event)(nil)

func (ev *Event) GetTransactionHash() binary.HexBytes {
	return ev.Header.TxHash
}

func (ev *Event) GetIndex() uint64 {
	return ev.Header.Index
}

func (ev *Event) GetTopics() []binary.Word256 {
	return ev.Log.Topics
}

func (ev *Event) GetData() []byte {
	return ev.Log.Data
}

func (ev *Event) GetAddress() crypto.Address {
	return ev.Log.Address
}

// Tags
func (ev *Event) Get(key string) (value interface{}, ok bool) {
	return (*exec.Event)(ev).Get(key)
}

func queryFromFilter(filter *chain.Filter) (query.Query, error) {
	if filter == nil || (len(filter.Topics) == 0 && len(filter.Addresses) == 0) {
		return new(query.Empty), nil
	}
	matchesFilter := query.NewBuilder()
	for _, address := range filter.Addresses {
		matchesFilter = matchesFilter.AndEquals("Address", address)
	}
	for i, topic := range filter.Topics {
		matchesFilter = matchesFilter.AndEquals(exec.LogNKey(i), topic)
	}
	// Note label vent's own EventTypeLabel has different casing!
	notLog := query.NewBuilder().AndNotEquals(event.EventTypeKey, exec.TypeLog)
	return matchesFilter.Or(notLog).Query()
}
