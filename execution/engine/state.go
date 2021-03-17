package engine

import (
	"github.com/klyed/hivesmartchain/execution/exec"
)

type State struct {
	*CallFrame
	Blockchain
	exec.EventSink
}
