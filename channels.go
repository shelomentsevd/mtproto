package mtproto

func (m *MTProto) ChannelsGetParticipants(peer TL, filter TL, offset int32, limit int32) (*TL, error)  {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		msg: TL_channels_getParticipants{
			Channel: peer,
			Filter:  filter,
			Offset:  offset,
			Limit:   limit,
		},
		resp: resp,
	}
	x := <-resp

	return &x, nil
}
