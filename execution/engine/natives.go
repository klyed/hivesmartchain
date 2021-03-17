package engine

import (
	"github.com/KLYE-Dev/HSC-MAIN/acm"
	"github.com/KLYE-Dev/HSC-MAIN/crypto"
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
