package mtproto

import "fmt"

func (m *MTProto) UsersGetFullUsers(id TL) (*TL_userFull, error) {
	var user TL_userFull
	tl, err := m.InvokeSync(TL_users_getFullUser{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	switch (*tl).(type) {
	case TL_userFull:
		user = (*tl).(TL_userFull)
	default:
		return nil, fmt.Errorf("Got: %T", *tl)
	}

	return &user, nil
}
