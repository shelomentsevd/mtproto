package mtproto

import "errors"

func (m *MTProto) ContactsGetTopPeers(correspondents, botsPM, botsInline, groups, channels bool, offset, limit, hash int32) (error, *TL) {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		msg: TL_contacts_getTopPeers{
			Correspondents: correspondents,
			Bots_pm:        botsPM,
			Bots_inline:    botsInline,
			Groups:         groups,
			Channels:       channels,
			Offset:         offset,
			Limit:          limit,
			Hash:           hash,
		},
		resp: resp,
	}
	x := <-resp

	switch x.(type) {
	case TL_contacts_topPeersNotModified:
	case TL_contacts_topPeers:
	default:
		return errors.New("MTProto::ContactsGetTopPeers error: Unknown type"), nil
	}

	return nil, &x
}
