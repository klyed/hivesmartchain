package contexts

import (
	"fmt"
	"math/big"

	"github.com/klyed/hivesmartchain/acm/acmstate"
	"github.com/klyed/hivesmartchain/acm/validator"
	"github.com/klyed/hivesmartchain/execution/exec"
	"github.com/klyed/hivesmartchain/logging"
	"github.com/klyed/hivesmartchain/txs/payload"
)

type UnbondContext struct {
	State        acmstate.ReaderWriter
	ValidatorSet validator.ReaderWriter
	Logger       *logging.Logger
	tx           *payload.UnbondTx
}

// Execute an UnbondTx to remove a validator
func (ctx *UnbondContext) Execute(txe *exec.TxExecution, p payload.Payload) error {
	var ok bool
	ctx.tx, ok = p.(*payload.UnbondTx)
	if !ok {
		return fmt.Errorf("payload must be UnbondTx, but is: %v", txe.Envelope.Tx.Payload)
	}

	if ctx.tx.Input.Address != ctx.tx.Output.Address {
		return fmt.Errorf("input and output address must match")
	}

	power := new(big.Int).SetUint64(ctx.tx.Output.GetAmount())
	account, err := ctx.State.GetAccount(ctx.tx.Input.Address)
	if err != nil {
		return err
	}

	err = account.AddToBalance(power.Uint64())
	if err != nil {
		return err
	}

	err = validator.SubtractPower(ctx.ValidatorSet, account.PublicKey, power)
	if err != nil {
		return err
	}

	return ctx.State.UpdateAccount(account)
}
