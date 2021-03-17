// +build integration,ethereum

package service_test

import (
	"testing"

	"github.com/klyed/hivesmartchain/rpc/web3/ethclient"
	"github.com/klyed/hivesmartchain/tests/web3/web3test"
	"github.com/klyed/hivesmartchain/vent/test"
	"github.com/stretchr/testify/require"
)

func TestEthereumConsumer(t *testing.T) {
	pk := web3test.GetPrivateKey(t)
	tcli := ethclient.NewTransactClient(ethclient.NewEthClient(web3test.GetChainRPCClient()))
	chainID, err := tcli.GetChainID()
	require.NoError(t, err)
	testConsumer(t, chainID, test.PostgresVentConfig(web3test.GetChainRemote()), tcli, pk.GetAddress())
}
