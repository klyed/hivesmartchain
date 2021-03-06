package config

import (
	"fmt"

	"github.com/klyed/hivesmartchain/config/source"
	"github.com/klyed/hivesmartchain/consensus/tendermint"
	"github.com/klyed/hivesmartchain/crypto"
	"github.com/klyed/hivesmartchain/execution"
	"github.com/klyed/hivesmartchain/genesis"
	"github.com/klyed/hivesmartchain/keys"
	"github.com/klyed/hivesmartchain/logging/logconfig"
	"github.com/klyed/hivesmartchain/rpc"
	tmConfig "github.com/tendermint/tendermint/config"
)

const DefaultBurrowConfigTOMLFileName = "hsc.toml"
const DefaultBurrowConfigEnvironmentVariable = "HSC_CONFIG_JSON"
const DefaultGenesisDocJSONFileName = "genesis.json"

type BurrowConfig struct {
	// Set on startup
	ValidatorAddress *crypto.Address `json:",omitempty" toml:",omitempty"`
	Passphrase       *string         `json:",omitempty" toml:",omitempty"`
	// From config file
	HscDir     string
	GenesisDoc *genesis.GenesisDoc        `json:",omitempty" toml:",omitempty"`
	Tendermint *tendermint.BurrowConfig   `json:",omitempty" toml:",omitempty"`
	Execution  *execution.ExecutionConfig `json:",omitempty" toml:",omitempty"`
	Keys       *keys.KeysConfig           `json:",omitempty" toml:",omitempty"`
	RPC        *rpc.RPCConfig             `json:",omitempty" toml:",omitempty"`
	Logging    *logconfig.LoggingConfig   `json:",omitempty" toml:",omitempty"`
}

func DefaultBurrowConfig() *BurrowConfig {
	return &BurrowConfig{
		HscDir:     ".hivesmartchain",
		Tendermint: tendermint.DefaultBurrowConfig(),
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
