package config

import (
	"fmt"
	"testing"

	"github.com/klye-dev/hsc-main/genesis"
)

func TestBurrowConfigSerialise(t *testing.T) {
	conf := &BurrowConfig{
		GenesisDoc: &genesis.GenesisDoc{
			ChainName: "Foo",
		},
	}
	fmt.Println(conf.JSONString())
}
