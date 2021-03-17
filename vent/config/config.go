package config

import (
	"time"

	"github.com/klye-dev/hsc-main/crypto"
	"github.com/klye-dev/hsc-main/vent/chain"
	"github.com/klye-dev/hsc-main/vent/sqlsol"
	"github.com/klye-dev/hsc-main/vent/types"
)

const DefaultPostgresDBURL = "postgres://postgres@localhost:5432/postgres?sslmode=disable"

// VentConfig is a set of configuration parameters
type VentConfig struct {
	DBAdapter           string
	DBURL               string
	DBSchema            string
	ChainAddress        string
	HTTPListenAddress   string
	BlockConsumerConfig chain.BlockConsumerConfig
	// Global contracts to watch specified as hex
	WatchAddresses []crypto.Address
	MinimumHeight  uint64
	SpecFileOrDirs []string
	AbiFileOrDirs  []string
	SpecOpt        sqlsol.SpecOpt
	// Announce status every AnnouncePeriod
	AnnounceEvery time.Duration
}

// DefaultFlags returns a configuration with default values
func DefaultVentConfig() *VentConfig {
	return &VentConfig{
		DBAdapter:         types.PostgresDB,
		DBURL:             DefaultPostgresDBURL,
		DBSchema:          "vent",
		ChainAddress:      "localhost:10997",
		HTTPListenAddress: "0.0.0.0:8080",
		SpecOpt:           sqlsol.None,
		AnnounceEvery:     time.Second * 5,
	}
}
