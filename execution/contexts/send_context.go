package contexts

import (
	"fmt"

	"github.com/klyed/hivesmartchain/acm/acmstate"
	"github.com/klyed/hivesmartchain/execution/errors"
	"github.com/klyed/hivesmartchain/execution/exec"
	"github.com/klyed/hivesmartchain/logging"
	"github.com/klyed/hivesmartchain/permission"
	"github.com/klyed/hivesmartchain/txs/payload"
)

type SendContext struct {
	State  acmstate.ReaderWriter
	Logger *logging.Logger
	tx     *payload.SendTx
}

func (ctx *SendContext) Execute(txe *exec.TxExecution, p payload.Payload) error {
	var ok bool
	ctx.tx, ok = p.(*payload.SendTx)
	if !ok {
		return fmt.Errorf("payload must be SendTx, but is: %v", txe.Envelope.Tx.Payload)
	}
	accounts, inTotal, err := getInputs(ctx.State, ctx.tx.Inputs)
	if err != nil {
		return err
	}

	// ensure all inputs have send permissions
	err = allHavePermission(ctx.State, permission.Send, accounts, ctx.Logger)
	if err != nil {
		return errors.Wrap(err, "at least one input lacks permission for SendTx")
	}

	// add outputs to accounts map
	// if any outputs don't exist, all inputs must have CreateAccount perm
	accounts, err = getOrMakeOutputs(ctx.State, accounts, ctx.tx.Outputs, ctx.Logger)
	if err != nil {
		return err
	}

	outTotal, err := validateOutputs(ctx.tx.Outputs)
	if err != nil {
		return err
	}
	if outTotal > inTotal {
		return errors.Codes.InsufficientFunds
	}
	if outTotal < inTotal {
		return errors.Codes.Overpayment
	}
	if outTotal == 0 {
		return errors.Codes.ZeroPayment
	}

	// Good! Adjust accounts
	err = adjustByInputs(accounts, ctx.tx.Inputs)
	if err != nil {
		return err
	}

	err = adjustByOutputs(accounts, ctx.tx.Outputs)
	if err != nil {
		return err
	}

	for _, acc := range accounts {
		err = ctx.State.UpdateAccount(acc)
		if err != nil {
			return err
		}
	}

	for _, i := range ctx.tx.Inputs {
		txe.Input(i.Address, nil)
	}

	for _, o := range ctx.tx.Outputs {
		txe.Output(o.Address, nil)
	}

	return nil
}
