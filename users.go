package mtproto

import "fmt"

func (m * MTProto) UsersGetFullUsers(id TL) (error, *TL_userFull) {
	var user TL_userFull
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		msg: TL_users_getFullUser{
			Id:id,
		},
		resp: resp,
	}

	x := <-resp

	switch x.(type) {
	case TL_userFull:
		user = x.(TL_userFull)
	case TL_rpc_error:
		x := x.(TL_rpc_error)
		return fmt.Errorf("RPC code: %d message: %s", x.Error_code, x.Error_message), nil
	default:
		return fmt.Errorf("Got: %T", x), nil
	}

	return nil, &user
}
