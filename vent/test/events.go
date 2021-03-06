package test

import (
	"context"
	"testing"

	"github.com/klyed/hivesmartchain/encoding"
	"github.com/klyed/hivesmartchain/rpc/web3/ethclient"
	"google.golang.org/grpc"

	"github.com/klyed/hivesmartchain/crypto"
	"github.com/klyed/hivesmartchain/execution/evm/abi"
	"github.com/klyed/hivesmartchain/execution/exec"
	"github.com/klyed/hivesmartchain/rpc/rpctransact"
	"github.com/klyed/hivesmartchain/txs/payload"
	"github.com/stretchr/testify/require"
)

const gasLimit = ethclient.BasicGasLimit * 100

type TransactClient interface {
	CallTxSync(ctx context.Context, in *payload.CallTx, opts ...grpc.CallOption) (*exec.TxExecution, error)
}

func NewHiveSmartChainTransactClient(t testing.TB, listenAddress string) rpctransact.TransactClient {
	t.Helper()

	conn, err := encoding.GRPCDial(listenAddress)
	require.NoError(t, err)
	return rpctransact.NewTransactClient(conn)
}

func CreateContract(t testing.TB, cli TransactClient, inputAddress crypto.Address) *exec.TxExecution {
	t.Helper()

	txe, err := cli.CallTxSync(context.Background(), &payload.CallTx{
		Input: &payload.TxInput{
			Address: inputAddress,
		},
		Address:  nil,
		Data:     Bytecode_EventsTest,
		GasLimit: gasLimit,
	})
	require.NoError(t, err)

	if txe.Exception != nil {
		t.Fatalf("call should not generate exception but returned: %v", txe.Exception.Error())
	}

	return txe
}

func CallRemoveEvent(t testing.TB, cli TransactClient, inputAddress, contractAddress crypto.Address,
	name string) *exec.TxExecution {
	return Call(t, cli, inputAddress, contractAddress, "removeThing", name)
}

func CallRemoveEvents(t testing.TB, cli TransactClient, inputAddress, contractAddress crypto.Address,
	name string) *exec.TxExecution {
	return Call(t, cli, inputAddress, contractAddress, "removeThings", name)
}

func CallAddEvent(t testing.TB, cli TransactClient, inputAddress, contractAddress crypto.Address,
	name, description string) *exec.TxExecution {
	return Call(t, cli, inputAddress, contractAddress, "addThing", name, description)
}

func CallAddEvents(t testing.TB, cli TransactClient, inputAddress, contractAddress crypto.Address,
	name, description string) *exec.TxExecution {
	return Call(t, cli, inputAddress, contractAddress, "addThings", name, description)
}

func Call(t testing.TB, cli TransactClient, inputAddress, contractAddress crypto.Address,
	functionName string, args ...interface{}) *exec.TxExecution {
	t.Helper()

	spec, err := abi.ReadSpec(Abi_EventsTest)
	require.NoError(t, err)

	data, _, err := spec.Pack(functionName, args...)
	require.NoError(t, err)

	txe, err := cli.CallTxSync(context.Background(), &payload.CallTx{
		Input: &payload.TxInput{
			Address: inputAddress,
		},
		Address:  &contractAddress,
		Data:     data,
		GasLimit: gasLimit,
	})
	require.NoError(t, err)

	if txe.Exception != nil {
		t.Fatalf("call should not generate exception but returned: %v", txe.Exception.Error())
	}

	return txe
}
