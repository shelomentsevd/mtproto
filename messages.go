package mtproto

func (m *MTProto) MessagesGetHistory(peer TL, offsetId, offsetDate, addOffset, limit, maxId, minId int32) (error, *TL) {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		msg: TL_messages_getHistory{
			Peer:        peer,
			Offset_id:   offsetId,
			Offset_date: offsetDate,
			Add_offset:  addOffset,
			Limit:       limit,
			Max_id:      maxId,
			Min_id:      minId,
		},
		resp: resp,
	}
	x := <-resp

	return nil, &x
}

func (m *MTProto) MessagesGetDialogs(excludePinned bool, offsetDate, offsetId int32, offsetPeer TL, limit int32) (error, *TL) {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		msg: TL_messages_getDialogs{
			Exclude_pinned: excludePinned,
			Offset_date:    offsetDate,
			Offset_id:      offsetId,
			Offset_peer:    offsetPeer,
			Limit:          limit,
		},
		resp: resp,
	}
	x := <-resp

	return nil, &x
}