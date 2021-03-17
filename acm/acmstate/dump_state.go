package acmstate

import (
	"bytes"
	"encoding/hex"
	"encoding/json"

	"github.com/klye-dev/hivesmartchain/acm"
	"github.com/klye-dev/hivesmartchain/crypto"
)

type DumpState struct {
	bytes.Buffer
}

func (dw *DumpState) UpdateAccount(updatedAccount *acm.Account) error {
	dw.WriteString("UpdateAccount\n")
	bs, err := json.Marshal(updatedAccount)
	if err != nil {
		return err
	}
	dw.Write(bs)
	dw.WriteByte('\n')
	return nil
}

func (dw *DumpState) RemoveAccount(address crypto.Address) error {
	dw.WriteString("RemoveAccount\n")
	dw.WriteString(address.String())
	dw.WriteByte('\n')
	return nil
}

func (dw *DumpState) SetStorage(address crypto.Address, key, value []byte) error {
	dw.WriteString("SetStorage\n")
	dw.WriteString(address.String())
	dw.WriteByte('/')
	dw.WriteString(hex.EncodeToString(key[:]))
	dw.WriteByte('/')
	dw.WriteString(hex.EncodeToString(value[:]))
	dw.WriteByte('\n')
	return nil
}
