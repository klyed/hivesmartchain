package config

import (
	"fmt"

	"github.com/klye-dev/hsc-main/config/source"
	"github.com/klye-dev/hsc-main/consensus/tendermint"
	"github.com/klye-dev/hsc-main/crypto"
	"github.com/klye-dev/hsc-main/execution"
	"github.com/klye-dev/hsc-main/genesis"
	"github.com/klye-dev/hsc-main/keys"
	"github.com/klye-dev/hsc-main/logging/logconfig"
	"github.com/klye-dev/hsc-main/rpc"
	tmConfig "github.com/tendermint/tendermint/config"
)

const DefaultBurrowConfigTOMLFileName = "burrow.toml"
const DefaultBurrowConfigEnvironmentVariable = "BURROW_CONFIG_JSON"
const DefaultGenesisDocJSONFileName = "genesis.json"

type BurrowConfig struct {
	// Set on startup
	ValidatorAddress *crypto.Address `json:",omitempty" toml:",omitempty"`
	Passphrase       *string         `json:",omitempty" toml:",omitempty"`
	// From config file
	HscDir  string
	GenesisDoc *genesis.GenesisDoc                `json:",omitempty" toml:",omitempty"`
	Tendermint *tendermint.BurrowTendermintConfig `json:",omitempty" toml:",omitempty"`
	Execution  *execution.ExecutionConfig         `json:",omitempty" toml:",omitempty"`
	Keys       *keys.KeysConfig                   `json:",omitempty" toml:",omitempty"`
	RPC        *rpc.RPCConfig                     `json:",omitempty" toml:",omitempty"`
	Logging    *logconfig.LoggingConfig           `json:",omitempty" toml:",omitempty"`
}

func DefaultBurrowConfig() *BurrowConfig {
	return &BurrowConfig{
		HscDir:  ".burrow",
		Tendermint: tendermint.DefaultBurrowTendermintConfig(),
		Keys:       keys.DefaultKeysConfig(),
		RPC:        rpc.DefaultRPCConfig(),
		Execution:  execution.DefaultExecutionConfig(),
		Logging:    logconfig.DefaultNodeLoggingConfig(),
	}
}

func (conf *BurrowConfig) Verify() error {
	if conf.ValidatorAddress == nil {
		return fmt.Errorf("could not finalise address - please provide one in config or via --account-address")
	}
	return nil
}

func (conf *BurrowConfig) TendermintConfig() (*tmConfig.Config, error) {
	return conf.Tendermint.Config(conf.HscDir, conf.Execution.TimeoutFactor)
}

func (conf *BurrowConfig) JSONString() string {
	return source.JSONString(conf)
}

func (conf *BurrowConfig) TOMLString() string {
	return source.TOMLString(conf)
}
