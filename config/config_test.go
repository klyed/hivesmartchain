package config

import (
	"fmt"
	"testing"

	"github.com/klyed/hivesmartchain/genesis"
)

func TestHiveSmartChainConfigSerialise(t *testing.T) {
	conf := &HiveSmartChainConfig{
		GenesisDoc: &genesis.GenesisDoc{
			ChainName: "Foo",
		},
	}
	fmt.Println(conf.JSONString())
}
