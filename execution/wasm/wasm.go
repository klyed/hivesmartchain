package wasm

import (
	"github.com/klye-dev/hivesmartchain/acm"
	"github.com/klye-dev/hivesmartchain/acm/acmstate"
	"github.com/klye-dev/hivesmartchain/execution/defaults"
	"github.com/klye-dev/hivesmartchain/execution/engine"
	"github.com/klye-dev/hivesmartchain/execution/errors"
	"github.com/klye-dev/hivesmartchain/execution/exec"
	"github.com/klye-dev/hivesmartchain/execution/native"
	lifeExec "github.com/perlin-network/life/exec"
)

// Implements ewasm, see https://github.com/ewasm/design
// WASM
var DefaultVMConfig = lifeExec.VMConfig{
	DisableFloatingPoint: true,
	MaxMemoryPages:       16,
	DefaultMemoryPages:   16,
}

type WVM struct {
	engine.Externals
	options            engine.Options
	vmConfig           lifeExec.VMConfig
	externalDispatcher engine.Dispatcher
}

func New(options engine.Options) *WVM {
	vm := &WVM{
		options:  defaults.CompleteOptions(options),
		vmConfig: DefaultVMConfig,
	}
	vm.externalDispatcher = engine.Dispatchers{&vm.Externals, options.Natives, vm}
	return vm
}

func Default() *WVM {
	return New(engine.Options{})
}

// RunWASM creates a WASM VM, and executes the given WASM contract code
func (vm *WVM) Execute(st acmstate.ReaderWriter, blockchain engine.Blockchain, eventSink exec.EventSink,
	params engine.CallParams, code []byte) (output []byte, cerr error) {
	defer func() {
		if r := recover(); r != nil {
			cerr = errors.Codes.ExecutionAborted
		}
	}()

	st = native.NewState(vm.options.Natives, st)

	state := engine.State{
		CallFrame:  engine.NewCallFrame(st).WithMaxCallStackDepth(vm.options.CallStackMaxDepth),
		Blockchain: blockchain,
		EventSink:  eventSink,
	}

	output, err := vm.Contract(code).Call(state, params)

	if err == nil {
		// Only sync back when there was no exception
		err = state.CallFrame.Sync()
	}
	// Always return output - we may have a reverted exception for which the return is meaningful
	return output, err
}

func (vm *WVM) Dispatch(acc *acm.Account) engine.Callable {
	if len(acc.WASMCode) == 0 {
		return nil
	}
	return vm.Contract(acc.WASMCode)
}

func (vm *WVM) Contract(code []byte) *Contract {
	return &Contract{
		vm:   vm,
		code: code,
	}
}
