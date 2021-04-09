package crypto

import (
	"encoding/json"
	"fmt"
	//hex "github.com/tmthrgd/go-hex"
)

type UserNameInterface interface {
	UserNameJSON()
}

type UserNameJSON struct {
	NativeName string
}

func (u *UserNameJSON) IsValid() bool {
	if u != nil {
		return true
	} else {
		return false
	}
}

func (u UserNameJSON) GetUserName() UserNameJSON {
	return u
}

func (u *UserNameJSON) String() string {
	return fmt.Sprintf("NativeName<UserName:%X>", u)
}

func (u *UserNameJSON) MarshalJSON() ([]byte, error) {
	jStruct := UserNameJSON{
		NativeName: fmt.Sprintf("%s", u),
	}
	txt, err := json.Marshal(jStruct)
	return txt, err
}

//func (u *UserName) String() string {
//	return hex.EncodeUpperToString(u)
//}
