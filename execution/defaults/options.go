package defaults

import (
	"github.com/klye-dev/hsc-main/execution/engine"
	"github.com/klye-dev/hsc-main/execution/native"
	"github.com/klye-dev/hsc-main/logging"
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
