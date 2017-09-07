package mtproto

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (m *MTProto) AuthSendCode(phonenumber string) (*TL_auth_sentCode, error) {
	var authSentCode TL_auth_sentCode
	flag := true
	for flag {
		x := <-m.InvokeAsync(TL_auth_sendCode{
			Allow_flashcall: false,
			Phone_number:    phonenumber,
			Current_number:  TL_boolTrue{},
			Api_id:          m.appConfig.Id,
			Api_hash:        m.appConfig.Hash,
		})

		if x.err != nil {
			// TODO: Maybe there are different ways to do it
			// MTProto connected to new DC(see handleRPCError), trying to get data again
			if strings.Contains(x.err.Error(), strconv.Itoa(errorSeeOther)) {
				continue
			}
			return nil, x.err
		}

		switch x.data.(type) {
		case TL_auth_sentCode:
			authSentCode = x.data.(TL_auth_sentCode)
			flag = false
		default:
			return nil, fmt.Errorf("Got: %T", x)
		}
	}

	return &authSentCode, nil
}

func (m *MTProto) AuthSignIn(phoneNumber, phoneCode, phoneCodeHash string) (*TL_auth_authorization, error) {
	if phoneNumber == "" || phoneCode == "" || phoneCodeHash == "" {
		return nil, errors.New("MRProto::AuthSignIn one of function parameters is empty")
	}

	x := <-m.InvokeAsync(TL_auth_signIn{
		Phone_number:    phoneNumber,
		Phone_code_hash: phoneCodeHash,
		Phone_code:      phoneCode,
	})
	if x.err != nil {
		return nil, x.err
	}

	auth, ok := x.data.(TL_auth_authorization)

	if !ok {
		return nil, fmt.Errorf("RPC: %#v", x)
	}

	return &auth, nil
}

func (m *MTProto) AuthLogOut() (bool, error) {
	var result bool
	x := <-m.InvokeAsync(TL_auth_logOut{})
	if x.err != nil {
		return result, x.err
	}

	result, err := ToBool(x.data)
	if err != nil {
		return result, err
	}

	return result, err
}
