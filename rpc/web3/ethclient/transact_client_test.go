// +build ethereum,integration

package ethclient

import (
	"context"
	"fmt"
	"testing"

	"github.com/klye-dev/hsc-main/execution/solidity"
	"github.com/klye-dev/hsc-main/tests/web3/web3test"
	"github.com/klye-dev/hsc-main/txs/payload"
	"github.com/stretchr/testify/require"
)

func TestEthTransactClient_CallTxSync(t *testing.T) {
	pk := web3test.GetPrivateKey(t)
	cli := NewTransactClient(NewEthClient(web3test.GetChainRPCClient())).WithAccounts(pk)
	input := pk.GetAddress()
	gasPrice, err := cli.GetGasPrice()
	require.NoError(t, err)
	nonce, err := cli.GetTransactionCount(input)
	require.NoError(t, err)
	txe, err := cli.CallTxSync(context.Background(), &payload.CallTx{
		Input: &payload.TxInput{
			Address:  input,
			Sequence: nonce,
		},
		GasPrice: gasPrice,
		GasLimit: BasicGasLimit * 10,
		Data:     solidity.Bytecode_EventEmitter,
	})
	require.NoError(t, err)
	fmt.Println(txe)
}
