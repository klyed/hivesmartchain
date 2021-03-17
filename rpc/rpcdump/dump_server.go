package rpcdump

import (
	"github.com/klye-dev/hsc-main/bcm"
	"github.com/klye-dev/hsc-main/dump"
	"github.com/klye-dev/hsc-main/execution/state"
	"github.com/klye-dev/hsc-main/logging"
)

type dumpServer struct {
	UnimplementedDumpServer
	dumper *dump.Dumper
}

var _ DumpServer = &dumpServer{}

func NewDumpServer(state *state.State, blockchain bcm.BlockchainInfo, logger *logging.Logger) *dumpServer {
	return &dumpServer{
		dumper: dump.NewDumper(state, blockchain).WithLogger(logger),
	}
}

func (ds *dumpServer) GetDump(param *GetDumpParam, stream Dump_GetDumpServer) error {
	return ds.dumper.Transmit(stream, 0, param.Height, dump.All)
}
