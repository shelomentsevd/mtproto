package mtproto

func (m *MTProto) UpdatesGetState() (error, *TL) {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		msg:  TL_updates_getState{},
		resp: resp,
	}
	x := <-resp

	return nil, &x
}

func (m *MTProto) UpdatesGetDifference(pts, ptsTotalLimit, date, qts int32) (error, *TL) {
	resp := make(chan TL, 1)
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

	return nil, &x
}

func (m *MTProto) UpdatesGetChannelDifference(force bool, channel, filter TL, pts, limit int32) (error, *TL) {
	resp := make(chan TL, 1)
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

	return nil, &x
}