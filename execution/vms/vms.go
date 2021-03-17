package vms

import (
	"github.com/KLYE-Dev/HSC-MAIN/execution/defaults"
	"github.com/KLYE-Dev/HSC-MAIN/execution/engine"
	"github.com/KLYE-Dev/HSC-MAIN/execution/evm"
	"github.com/KLYE-Dev/HSC-MAIN/execution/wasm"
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
