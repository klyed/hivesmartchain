package execution

import (
	"github.com/klyed/hivesmartchain/acm"
	"github.com/klyed/hivesmartchain/acm/acmstate"
	"github.com/klyed/hivesmartchain/bcm"
	"github.com/klyed/hivesmartchain/crypto"
	"github.com/klyed/hivesmartchain/execution/contexts"
	"github.com/klyed/hivesmartchain/execution/engine"
	"github.com/klyed/hivesmartchain/execution/exec"
	"github.com/klyed/hivesmartchain/execution/vms"
	"github.com/klyed/hivesmartchain/logging"
	"github.com/klyed/hivesmartchain/txs"
	"github.com/klyed/hivesmartchain/txs/payload"
)

// Run a contract's code on an isolated and unpersisted state
// Cannot be used to create new contracts
func CallSim(reader acmstate.Reader, blockchain bcm.BlockchainInfo, fromAddress, address crypto.Address, data []byte,
	logger *logging.Logger) (*exec.TxExecution, error) {

	cache := acmstate.NewCache(reader)
	exe := contexts.CallContext{
		VMS:           vms.NewConnectedVirtualMachines(engine.Options{}),
		RunCall:       true,
		State:         cache,
		MetadataState: acmstate.NewMemoryState(),
		Blockchain:    blockchain,
		Logger:        logger,
	}

	txe := exec.NewTxExecution(txs.Enclose(blockchain.ChainID(), &payload.CallTx{
		Input: &payload.TxInput{
			Address: fromAddress,
		},
		Address:  &address,
		Data:     data,
		GasLimit: contexts.GasLimit,
	}))

	// Set height for downstream synchronisation purposes
	txe.Height = blockchain.LastBlockHeight()
	err := exe.Execute(txe, txe.Envelope.Tx.Payload)
	if err != nil {
		return nil, err
	}
	return txe, nil
}

// Run the given code on an isolated and unpersisted state
// Cannot be used to create new contracts.
func CallCodeSim(reader acmstate.Reader, blockchain bcm.BlockchainInfo, fromAddress, address crypto.Address, code, data []byte,
	logger *logging.Logger) (*exec.TxExecution, error) {

	// Attach code to target account (overwriting target)
	cache := acmstate.NewCache(reader)
	err := cache.UpdateAccount(&acm.Account{
		Address: address,
		EVMCode: code,
	})

	if err != nil {
		return nil, err
	}
	return CallSim(cache, blockchain, fromAddress, address, data, logger)
}
