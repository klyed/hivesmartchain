package engine

import (
	"github.com/KLYE-Dev/HSC-MAIN/execution/exec"
)

type State struct {
	*CallFrame
	Blockchain
	exec.EventSink
}
