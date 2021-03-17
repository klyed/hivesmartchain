package validator

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/klye-dev/hsc-main/acm"
	"github.com/klye-dev/hsc-main/crypto"
	"github.com/stretchr/testify/assert"
)

func TestValidators_AlterPower(t *testing.T) {
	vs := NewSet()
	pow1 := big.NewInt(2312312321)
	pubA := pubKey(1)
	vs.ChangePower(pubA, pow1)
	assert.Equal(t, pow1, vs.TotalPower())
	vs.ChangePower(pubA, big.NewInt(0))
	assertZero(t, vs.TotalPower())
}

func pubKey(secret interface{}) *crypto.PublicKey {
	return acm.NewAccountFromSecret(fmt.Sprintf("%v", secret)).PublicKey
}
