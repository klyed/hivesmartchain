package vms

import (
	"github.com/klye-dev/hsc-main/execution/defaults"
	"github.com/klye-dev/hsc-main/execution/engine"
	"github.com/klye-dev/hsc-main/execution/evm"
	"github.com/klye-dev/hsc-main/execution/wasm"
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
