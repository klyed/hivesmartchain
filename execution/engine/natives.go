package engine

import (
	"github.com/klye-dev/hsc-main/acm"
	"github.com/klye-dev/hsc-main/crypto"
)

type Native interface {
	Callable
	SetExternals(externals Dispatcher)
	ContractMeta() []*acm.ContractMeta
	FullName() string
	Address() crypto.Address
}

type Natives interface {
	ExternalDispatcher
	GetByAddress(address crypto.Address) Native
}
