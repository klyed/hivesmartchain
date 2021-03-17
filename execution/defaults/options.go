package defaults

import (
	"github.com/klyed/hivesmartchain/execution/engine"
	"github.com/klyed/hivesmartchain/execution/native"
	"github.com/klyed/hivesmartchain/logging"
)

func CompleteOptions(options engine.Options) engine.Options {
	// Set defaults
	if options.MemoryProvider == nil {
		options.MemoryProvider = engine.DefaultDynamicMemoryProvider
	}
	if options.Logger == nil {
		options.Logger = logging.NewNoopLogger()
	}
	if options.Natives == nil {
		options.Natives = native.MustDefaultNatives()
	}
	return options
}
