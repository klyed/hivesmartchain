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

const DefaultBurrowTendermintConfigTOMLFileName = "hsc.toml"
const DefaultBurrowTendermintConfigEnvironmentVariable = "HSC_CONFIG_JSON"
const DefaultGenesisDocJSONFileName = "genesis.json"

type BurrowTendermintConfig struct {
	// Set on startup
	ValidatorAddress *crypto.Address `json:",omitempty" toml:",omitempty"`
	Passphrase       *string         `json:",omitempty" toml:",omitempty"`
	// From config file
	HscDir     string
	GenesisDoc *genesis.GenesisDoc                `json:",omitempty" toml:",omitempty"`
	Tendermint *tendermint.BurrowTendermintConfig `json:",omitempty" toml:",omitempty"`
	Execution  *execution.ExecutionConfig         `json:",omitempty" toml:",omitempty"`
	Keys       *keys.KeysConfig                   `json:",omitempty" toml:",omitempty"`
	RPC        *rpc.RPCConfig                     `json:",omitempty" toml:",omitempty"`
	Logging    *logconfig.LoggingConfig           `json:",omitempty" toml:",omitempty"`
}

func DefaultBurrowTendermintConfig() *BurrowTendermintConfig {
	return &BurrowTendermintConfig{
		HscDir:     ".hivesmartchain",
		Tendermint: tendermint.DefaultBurrowTendermintConfig(),
		Keys:       keys.DefaultKeysConfig(),
		RPC:        rpc.DefaultRPCConfig(),
		Execution:  execution.DefaultExecutionConfig(),
		Logging:    logconfig.DefaultNodeLoggingConfig(),
	}
}

func (conf *BurrowTendermintConfig) Verify() error {
	conf.P2P.AddrBookStrict = false
	if conf.ValidatorAddress == nil {
		return fmt.Errorf("could not finalise address - please provide one in config or via --account-address")
	}
	return nil
}

func (conf *BurrowTendermintConfig) TendermintConfig() (*tmConfig.Config, error) {
	return conf.Tendermint.Config(conf.HscDir, conf.Execution.TimeoutFactor)
}

func (conf *BurrowTendermintConfig) JSONString() string {
	return source.JSONString(conf)
}

func (conf *BurrowTendermintConfig) TOMLString() string {
	return source.TOMLString(conf)
}
