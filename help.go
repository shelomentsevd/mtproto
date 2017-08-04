package mtproto

func (m *MTProto) HelpGetConfig() (*TL, error) {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		msg: TL_help_getConfig{},
		resp: resp,
	}
	x := <-resp

	return &x, nil
}

func (m *MTProto) HelpGetNearestDc() (*TL, error) {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		msg: TL_help_getNearestDc{},
		resp: resp,
	}
	x := <-resp

	return &x, nil
}
