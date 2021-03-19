package config

import (
	"fmt"
	"testing"

	"github.com/klyed/hivesmartchain/genesis"
)

func TestBurrowTendermintConfigSerialise(t *testing.T) {
	conf := &BurrowTendermintConfig{
		GenesisDoc: &genesis.GenesisDoc{
			ChainName: "Foo",
		},
	}
	fmt.Println(conf.JSONString())
}
