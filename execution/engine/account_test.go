package engine

import (
	"testing"

	"github.com/klyed/hivesmartchain/acm"
	"github.com/klyed/hivesmartchain/acm/acmstate"
	"github.com/klyed/hivesmartchain/crypto"
	"github.com/klyed/hivesmartchain/execution/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestState_CreateAccount(t *testing.T) {
	st := acmstate.NewMemoryState()
	address := AddressFromName("frogs")
	err := CreateAccount(st, address)
	require.NoError(t, err)
	err = CreateAccount(st, address)
	require.Error(t, err)
	require.Equal(t, errors.Codes.DuplicateAddress, errors.GetCode(err))

	st = acmstate.NewMemoryState()
	err = CreateAccount(st, address)
	require.NoError(t, err)
	err = InitEVMCode(st, address, []byte{1, 2, 3})
	require.NoError(t, err)
}

func TestState_Sync(t *testing.T) {
	backend := acmstate.NewCache(acmstate.NewMemoryState())
	st := NewCallFrame(backend)
	address := AddressFromName("frogs")

	err := CreateAccount(st, address)
	require.NoError(t, err)
	amt := uint64(1232)
	addToBalance(t, st, address, amt)
	err = st.Sync()
	require.NoError(t, err)

	acc, err := backend.GetAccount(address)
	require.NoError(t, err)
	assert.Equal(t, acc.Balance, amt)
}

func TestState_NewCache(t *testing.T) {
	st := NewCallFrame(acmstate.NewMemoryState())
	address := AddressFromName("frogs")

	cache, err := st.NewFrame()
	require.NoError(t, err)
	err = CreateAccount(cache, address)
	require.NoError(t, err)
	amt := uint64(1232)
	addToBalance(t, cache, address, amt)

	acc, err := st.GetAccount(address)
	require.NoError(t, err)
	require.Nil(t, acc)

	// Sync through to cache
	err = cache.Sync()
	require.NoError(t, err)

	acc, err = MustAccount(cache, address)
	require.NoError(t, err)
	assert.Equal(t, amt, acc.Balance)

	cache, err = cache.NewFrame(acmstate.ReadOnly)
	require.NoError(t, err)
	cache, err = cache.NewFrame()
	require.NoError(t, err)
	err = UpdateAccount(cache, address, func(account *acm.Account) error {
		return account.AddToBalance(amt)
	})
	require.Error(t, err)
	require.Equal(t, errors.Codes.IllegalWrite, errors.GetCode(err))
}

func addToBalance(t testing.TB, st acmstate.ReaderWriter, address crypto.Address, amt uint64) {
	err := UpdateAccount(st, address, func(account *acm.Account) error {
		return account.AddToBalance(amt)
	})
	require.NoError(t, err)
}
