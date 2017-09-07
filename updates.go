package mtproto

func (m *MTProto) UpdatesGetState() (*TL, error) {
	return m.InvokeSync(TL_updates_getState{})
}

func (m *MTProto) UpdatesGetDifference(pts, ptsTotalLimit, date, qts int32) (*TL, error) {
	return m.InvokeSync(TL_updates_getDifference{
		Pts:             pts,
		Pts_total_limit: ptsTotalLimit,
		Date:            date,
		Qts:             qts,
	})
}

func (m *MTProto) UpdatesGetChannelDifference(force bool, channel, filter TL, pts, limit int32) (*TL, error) {
	return m.InvokeSync(TL_updates_getChannelDifference{
		Force:   force,
		Channel: channel,
		Filter:  filter,
		Pts:     pts,
		Limit:   limit,
	})
}
