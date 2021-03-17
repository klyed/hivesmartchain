package engine

import (
	"github.com/klyed/hivesmartchain/acm"
	"github.com/klyed/hivesmartchain/crypto"
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
