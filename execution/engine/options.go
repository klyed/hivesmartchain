package engine

import (
	"github.com/klye-dev/hivesmartchain/execution/errors"
	"github.com/klye-dev/hivesmartchain/logging"
)

// Options are parameters that are generally stable across a burrow configuration.
// Defaults will be used for any zero values.
type Options struct {
	MemoryProvider           func(errors.Sink) Memory
	Natives                  Natives
	Nonce                    []byte
	DebugOpcodes             bool
	DumpTokens               bool
	CallStackMaxDepth        uint64
	DataStackInitialCapacity uint64
	DataStackMaxDepth        uint64
	Logger                   *logging.Logger
}
