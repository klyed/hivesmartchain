package engine

import (
	"github.com/klye-dev/hsc-main/execution/exec"
)

type State struct {
	*CallFrame
	Blockchain
	exec.EventSink
}
