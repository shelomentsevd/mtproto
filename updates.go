package mtproto

func (m *MTProto) UpdatesGetState() (*TL, error) {
	resp := make(chan response, 1)
	m.queueSend <- packetToSend{
		msg:  TL_updates_getState{},
		resp: resp,
	}
	x := <-resp
	if x.err != nil {
		return nil, x.err
	}

	return &x.data, nil
}

func (m *MTProto) UpdatesGetDifference(pts, ptsTotalLimit, date, qts int32) (*TL, error) {
	resp := make(chan response, 1)
	m.queueSend <- packetToSend{
		msg: TL_updates_getDifference{
			Pts:             pts,
			Pts_total_limit: ptsTotalLimit,
			Date:            date,
			Qts:             qts,
		},
		resp: resp,
	}

	x := <-resp
	if x.err != nil {
		return nil, x.err
	}

	return &x.data, nil
}

func (m *MTProto) UpdatesGetChannelDifference(force bool, channel, filter TL, pts, limit int32) (*TL, error) {
	resp := make(chan response, 1)
	m.queueSend <- packetToSend{
		msg: TL_updates_getChannelDifference{
			Force:   force,
			Channel: channel,
			Filter:  filter,
			Pts:     pts,
			Limit:   limit,
		},
		resp: resp,
	}

	x := <-resp
	if x.err != nil {
		return nil, x.err
	}

	return &x.data, nil
}
