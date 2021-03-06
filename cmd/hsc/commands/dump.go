package commands

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/klyed/hivesmartchain/encoding"

	cli "github.com/jawher/mow.cli"
	"github.com/klyed/hivesmartchain/core"
	"github.com/klyed/hivesmartchain/dump"
	"github.com/klyed/hivesmartchain/logging/logconfig"
	"github.com/klyed/hivesmartchain/rpc/rpcdump"
	"github.com/klyed/hivesmartchain/rpc/rpcquery"
)

type dumpOptions struct {
	height            *int
	filename          *string
	useBinaryEncoding *bool
}

func maybeOutput(verbose *bool, output Output, format string, args ...interface{}) {
	if verbose != nil && *verbose {
		output.Logf(format, args)
	}
}

func addDumpOptions(cmd *cli.Cmd, specOptions ...string) *dumpOptions {
	cmd.Spec += "[--height=<state height to dump at>] [--binary]"
	for _, spec := range specOptions {
		cmd.Spec += " " + spec
	}
	cmd.Spec += "[FILE]"
	return &dumpOptions{
		height:            cmd.IntOpt("h height", 0, "Block height to dump to, defaults to latest block height"),
		useBinaryEncoding: cmd.BoolOpt("b binary", false, "Output in binary encoding (default is JSON)"),
		filename:          cmd.StringArg("FILE", "", "Location to output dump, if no argument is given then this streams to STDOUT"),
	}
}

// Dump saves the state from a remote chain
func Dump(output Output) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		verbose := cmd.BoolOpt("v verbose", false, "Verbose debug information")
		cmd.Spec += "[--verbose]"

		cmd.Command("local", "create a dump from local Hive Smart Chain directory", func(cmd *cli.Cmd) {
			maybeOutput(verbose, output, "dumping from local Hive Smart Chain dir")
			configFileOpt := cmd.String(configFileOption)
			genesisFileOpt := cmd.String(genesisFileOption)

			dumpOpts := addDumpOptions(cmd, configFileSpec, genesisFileSpec)

			cmd.Action = func() {
				conf, err := obtainDefaultConfig(*configFileOpt, *genesisFileOpt)
				if err != nil {
					output.Fatalf("could not obtain config: %v", err)
				}

				kern, err := core.NewKernel(conf.HscDir)
				if err != nil {
					output.Fatalf("could not create Hive Smart Chain kernel: %v", err)
				}

				err = kern.LoadState(conf.GenesisDoc)
				if err != nil {
					output.Fatalf("could not load Hive Smart Chain state: %v", err)
				}

				// Include all logging by default
				logger, err := logconfig.New().Logger()
				if err != nil {
					output.Fatalf("could not make logger: %v", err)
				}

				source := dump.NewDumper(kern.State, kern.Blockchain).WithLogger(logger).
					Source(0, uint64(*dumpOpts.height), dump.All)

				err = dumpToFile(*dumpOpts.filename, source, *dumpOpts.useBinaryEncoding)
				if err != nil {
					output.Fatalf("could not dump to file %s': %v", *dumpOpts.filename, err)
				}
				maybeOutput(verbose, output, "dump successfully written to '%s'", *dumpOpts.filename)
			}
		})

		cmd.Command("remote", "pull a dump from a remote Hive Smart Chain node", func(cmd *cli.Cmd) {
			chainURLOpt := cmd.StringOpt("c chain", "127.0.0.1:10997", "chain to be used in IP:PORT format")
			timeoutOpt := cmd.IntOpt("t timeout", 0, "Timeout in seconds")
			dumpOpts := addDumpOptions(cmd, "[--chain=<chain GRPC address>]", "[--timeout=<GRPC timeout seconds>]")

			cmd.Action = func() {
				maybeOutput(verbose, output, "dumping from remote chain at %s", *chainURLOpt)

				ctx, cancel := context.WithCancel(context.Background())
				if *timeoutOpt != 0 {
					timeout := time.Duration(*timeoutOpt) * time.Second
					ctx, cancel = context.WithTimeout(context.Background(), timeout)
				}
				defer cancel()

				conn, err := encoding.GRPCDialContext(ctx, *chainURLOpt)
				if err != nil {
					output.Fatalf("failed to connect: %v", err)
				}

				qCli := rpcquery.NewQueryClient(conn)
				chainStatus, err := qCli.Status(context.Background(), &rpcquery.StatusParam{})
				if err != nil {
					output.Logf("could not get chain status: %v", err)
				}
				stat, err := json.Marshal(chainStatus)
				if err != nil {
					output.Logf("failed to marshal: %v", err)
				}
				maybeOutput(verbose, output, "dumping from chain: %s", string(stat))

				dc := rpcdump.NewDumpClient(conn)
				receiver, err := dc.GetDump(ctx, &rpcdump.GetDumpParam{Height: uint64(*dumpOpts.height)})
				if err != nil {
					output.Fatalf("failed to retrieve dump: %v", err)
				}

				err = dumpToFile(*dumpOpts.filename, receiver, *dumpOpts.useBinaryEncoding)
				if err != nil {
					output.Fatalf("could not dump to file %s': %v", *dumpOpts.filename, err)
				}

				maybeOutput(verbose, output, "dump successfully written to '%s'", *dumpOpts.filename)
			}
		})
	}
}

func dumpToFile(filename string, source dump.Source, useBinaryEncoding bool) error {
	var file *os.File
	var err error
	if filename == "" {
		file = os.Stdout
	} else {
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
	}

	// Receive
	err = dump.Write(file, source, useBinaryEncoding, dump.All)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}
