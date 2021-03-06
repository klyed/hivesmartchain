// Copyright Monax Industries Limited
// SPDX-License-Identifier: Apache-2.0

package native

import (
	"fmt"

	"github.com/klyed/hivesmartchain/acm"
	"github.com/klyed/hivesmartchain/crypto"
	"github.com/klyed/hivesmartchain/execution/engine"
	"github.com/klyed/hivesmartchain/logging"
	"github.com/klyed/hivesmartchain/permission"
)

type Natives struct {
	callableByAddress map[crypto.Address]engine.Native
	callableByName    map[string]engine.Native
	logger            *logging.Logger
}

func New() *Natives {
	return &Natives{
		callableByAddress: make(map[crypto.Address]engine.Native),
		callableByName:    make(map[string]engine.Native),
		logger:            logging.NewNoopLogger(),
	}
}

func Merge(nss ...*Natives) (*Natives, error) {
	n := New()
	for _, ns := range nss {
		for _, contract := range ns.callableByName {
			err := n.register(contract)
			if err != nil {
				return nil, err
			}
		}
	}
	return n, nil
}

func (ns *Natives) WithLogger(logger *logging.Logger) *Natives {
	ns.logger = logger
	return ns
}

func (ns *Natives) Dispatch(acc *acm.Account) engine.Callable {
	return ns.GetByAddress(acc.Address)
}

func (ns *Natives) SetExternals(externals engine.Dispatcher) {
	for _, c := range ns.callableByAddress {
		c.SetExternals(externals)
	}
}

func (ns *Natives) Callables() []engine.Callable {
	callables := make([]engine.Callable, 0, len(ns.callableByAddress))
	for _, c := range ns.callableByAddress {
		callables = append(callables, c)
	}
	return callables
}

func (ns *Natives) GetByName(name string) engine.Native {
	return ns.callableByName[name]
}

func (ns *Natives) GetContract(name string) *Contract {
	c, _ := ns.GetByName(name).(*Contract)
	return c
}

func (ns *Natives) GetFunction(name string) *Function {
	f, _ := ns.GetByName(name).(*Function)
	return f
}

func (ns *Natives) GetByAddress(address crypto.Address) engine.Native {
	return ns.callableByAddress[address]
}

func (ns *Natives) IsRegistered(address crypto.Address) bool {
	_, ok := ns.callableByAddress[address]
	return ok
}

func (ns *Natives) MustContract(name, comment string, functions ...Function) *Natives {
	ns, err := ns.Contract(name, comment, functions...)
	if err != nil {
		panic(err)
	}
	return ns
}

func (ns *Natives) Contract(name, comment string, functions ...Function) (*Natives, error) {
	contract, err := NewContract(name, comment, ns.logger, functions...)
	if err != nil {
		return nil, err
	}
	err = ns.register(contract)
	if err != nil {
		return nil, err
	}
	return ns, nil
}

func (ns *Natives) MustFunction(comment string, address crypto.Address, permFlag permission.PermFlag, f interface{}) *Natives {
	ns, err := ns.Function(comment, address, permFlag, f)
	if err != nil {
		panic(err)
	}
	return ns
}

func (ns *Natives) Function(comment string, address crypto.Address, permFlag permission.PermFlag, f interface{}) (*Natives, error) {
	function, err := NewFunction(comment, address, permFlag, f)
	if err != nil {
		return nil, err
	}
	err = ns.register(function)
	if err != nil {
		return nil, err
	}
	return ns, nil
}

func (ns *Natives) register(callable engine.Native) error {
	name := callable.FullName()
	address := callable.Address()
	_, ok := ns.callableByName[name]
	if ok {
		return fmt.Errorf("cannot redeclare contract with name %s", name)
	}
	_, ok = ns.callableByAddress[address]
	if ok {
		return fmt.Errorf("cannot redeclare contract with address %v", address)
	}
	ns.callableByName[name] = callable
	ns.callableByAddress[address] = callable
	return nil
}
