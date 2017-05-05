package mtproto

import ("fmt"
	"errors")

func (m *MTProto) AuthSendCode(phonenumber string) (error, *TL_auth_sentCode) {
	var authSentCode TL_auth_sentCode
	flag := true
	for flag {
		resp := make(chan TL, 1)
		m.queueSend <- packetToSend{
			msg: TL_auth_sendCode{
				Allow_flashcall: false,
				Phone_number:    phonenumber,
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
				return fmt.Errorf("RPC Error_code: %d", x.Error_code), nil
			}
			var newDc int32
			n, _ := fmt.Sscanf(x.Error_message, "PHONE_MIGRATE_%d", &newDc)
			if n != 1 {
				n, _ := fmt.Sscanf(x.Error_message, "NETWORK_MIGRATE_%d", &newDc)
				if n != 1 {
					return fmt.Errorf("RPC error_string: %s", x.Error_message), nil
				}
			}

			newDcAddr, ok := m.dclist[newDc]
			if !ok {
				return fmt.Errorf("Wrong DC index: %d", newDc), nil
			}
			err := m.reconnect(newDcAddr)
			if err != nil {
				return err, nil
			}
		default:
			return fmt.Errorf("Got: %T", x), nil
		}
	}

	return nil, &authSentCode
}

func (m *MTProto) AuthSignIn(phoneNumber, phoneCode, phoneCodeHash string) (error, *TL_auth_authorization) {
	if phoneNumber == "" || phoneCode == "" || phoneCodeHash == "" {
		return errors.New("MRProto::AuthSignIn one of function parameters is empty"), nil
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
		return fmt.Errorf("RPC: %#v", x), nil
	}

	return nil, &auth
}

