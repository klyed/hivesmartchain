package engine

import (
	"github.com/klye-dev/hivesmartchain/execution/exec"
)

type State struct {
	*CallFrame
	Blockchain
	exec.EventSink
}
