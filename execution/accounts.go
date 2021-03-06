package execution

import (
	"sync"

	"github.com/klyed/hivesmartchain/acm"
	"github.com/klyed/hivesmartchain/acm/acmstate"
	"github.com/klyed/hivesmartchain/crypto"
	"github.com/klyed/hivesmartchain/keys"
	hsc_sync "github.com/klyed/hivesmartchain/sync"
)

// Accounts pairs an underlying state.Reader with a KeyClient to provide a signing variant of an account
// it also maintains a lock over addresses to provide a linearisation of signing events using SequentialSigningAccount
type Accounts struct {
	hsc_sync.RingMutex
	acmstate.Reader
	keyClient keys.KeyClient
}

type SigningAccount struct {
	*acm.Account
	crypto.Signer
}

type SequentialSigningAccount struct {
	Address       crypto.Address
	accountLocker sync.Locker
	getter        func() (*SigningAccount, error)
}

func NewAccounts(reader acmstate.Reader, keyClient keys.KeyClient, mutexCount int) *Accounts {
	return &Accounts{
		RingMutex: *hsc_sync.NewRingMutexNoHash(mutexCount),
		Reader:    reader,
		keyClient: keyClient,
	}
}
func (accs *Accounts) SigningAccount(address crypto.Address) (*SigningAccount, error) {
	signer, err := keys.AddressableSigner(accs.keyClient, address)
	if err != nil {
		return nil, err
	}
	account, err := accs.GetAccount(address)
	if err != nil {
		return nil, err
	}
	// If the account is unknown to us return a zeroed account
	if account == nil {
		account = &acm.Account{
			Address: address,
		}
	}
	pubKey, err := accs.keyClient.PublicKey(address)
	if err != nil {
		return nil, err
	}
	account.PublicKey = pubKey
	return &SigningAccount{
		Account: account,
		Signer:  signer,
	}, nil
}

func (accs *Accounts) SequentialSigningAccount(address crypto.Address) (*SequentialSigningAccount, error) {
	return &SequentialSigningAccount{
		Address:       address,
		accountLocker: accs.Mutex(address.Bytes()),
		getter: func() (*SigningAccount, error) {
			return accs.SigningAccount(address)
		},
	}, nil
}

type UnlockFunc func()

func (ssa *SequentialSigningAccount) Lock() (*SigningAccount, UnlockFunc, error) {
	ssa.accountLocker.Lock()
	account, err := ssa.getter()
	if err != nil {
		ssa.accountLocker.Unlock()
		return nil, nil, err
	}
	return account, ssa.accountLocker.Unlock, err
}
