package mtproto

import (
	"errors"
	"fmt"
)

func (m *MTProto) AuthSendCode(phonenumber string) (*TL_auth_sentCode, error) {
	var authSentCode TL_auth_sentCode
	flag := true
	for flag {
		resp := make(chan TL, 1)
		m.queueSend <- packetToSend{
			msg: TL_auth_sendCode{
				Allow_flashcall: false,
				Phone_number:    phonenumber,
				Current_number:  TL_boolTrue{},
				Api_id:          m.appConfig.Id,
				Api_hash:        m.appConfig.Hash,
			},
			resp: resp,
		}
		x := <-resp
		switch x.(type) {
		case TL_auth_sentCode:
			authSentCode = x.(TL_auth_sentCode)
			flag = false
		case TL_rpc_error:
			x := x.(TL_rpc_error)
			if x.Error_code != errorSeeOther {
				return nil, fmt.Errorf("RPC Error_code: %d", x.Error_code)
			}
			var newDc int32
			n, _ := fmt.Sscanf(x.Error_message, "PHONE_MIGRATE_%d", &newDc)
			if n != 1 {
				n, _ := fmt.Sscanf(x.Error_message, "NETWORK_MIGRATE_%d", &newDc)
				if n != 1 {
					return nil, fmt.Errorf("RPC error_string: %s", x.Error_message)
				}
			}

			newDcAddr, ok := m.dclist[newDc]
			if !ok {
				return nil, fmt.Errorf("Wrong DC index: %d", newDc)
			}
			err := m.reconnect(newDcAddr)
			if err != nil {
				return nil, err
			}
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

	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		msg: TL_auth_signIn{
			Phone_number:    phoneNumber,
			Phone_code_hash: phoneCodeHash,
			Phone_code:      phoneCode,
		},
		resp: resp,
	}
	x := <-resp
	auth, ok := x.(TL_auth_authorization)

	if !ok {
		return nil, fmt.Errorf("RPC: %#v", x)
	}

	return &auth, nil
}

func (m *MTProto) AuthLogOut() (bool, error) {
	var result bool
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		msg:  TL_auth_logOut{},
		resp: resp,
	}
	x := <-resp

	result, err := ToBool(x)
	if err != nil {
		return result, err
	}

	return result, err
}
