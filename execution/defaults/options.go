package defaults

import (
	"github.com/KLYE-Dev/HSC-MAIN/execution/engine"
	"github.com/KLYE-Dev/HSC-MAIN/execution/native"
	"github.com/KLYE-Dev/HSC-MAIN/logging"
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
