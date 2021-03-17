package commands

import (
	"github.com/KLYE-Dev/HSC-MAIN/core"
	cli "github.com/jawher/mow.cli"
)

// Start launches the burrow daemon
func Start(output Output) func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {
		configOpts := addConfigOptions(cmd)

		cmd.Action = func() {
			conf, err := configOpts.obtainBurrowConfig()
			if err != nil {
				output.Fatalf("could not set up config: %v", err)
			}

			if err := conf.Verify(); err != nil {
				output.Fatalf("cannot continue with config: %v", err)
			}

			output.Logf("Using validator address: %s", *conf.ValidatorAddress)

			kern, err := core.LoadKernelFromConfig(conf)
			if err != nil {
				output.Fatalf("could not configure Hive Smart Chain kernel: %v", err)
			}

			if err = kern.Boot(); err != nil {
				output.Fatalf("could not boot Hive Smart Chain kernel: %v", err)
			}

			kern.WaitForShutdown()
		}
	}
}
