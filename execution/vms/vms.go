package vms

import (
	"github.com/klyed/hivesmartchain/execution/defaults"
	"github.com/klyed/hivesmartchain/execution/engine"
	"github.com/klyed/hivesmartchain/execution/evm"
	"github.com/klyed/hivesmartchain/execution/wasm"
)

type VirtualMachines struct {
	*evm.EVM
	*wasm.WVM
}

func NewConnectedVirtualMachines(options engine.Options) *VirtualMachines {
	options = defaults.CompleteOptions(options)
	evm := evm.New(options)
	wvm := wasm.New(options)
	// Allow the virtual machines to call each other
	engine.Connect(evm, wvm)
	return &VirtualMachines{
		EVM: evm,
		WVM: wvm,
	}
}
