package execution

import (
	"github.com/klye-dev/hsc-main/acm"
	"github.com/klye-dev/hsc-main/acm/acmstate"
	"github.com/klye-dev/hsc-main/bcm"
	"github.com/klye-dev/hsc-main/crypto"
	"github.com/klye-dev/hsc-main/execution/contexts"
	"github.com/klye-dev/hsc-main/execution/engine"
	"github.com/klye-dev/hsc-main/execution/exec"
	"github.com/klye-dev/hsc-main/execution/vms"
	"github.com/klye-dev/hsc-main/logging"
	"github.com/klye-dev/hsc-main/txs"
	"github.com/klye-dev/hsc-main/txs/payload"
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
