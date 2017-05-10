package mtproto

import "errors"

func (m *MTProto) ContactsGetContacts(hash string) (*TL, error)  {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		msg: TL_contacts_getContacts{
			Hash: hash,
		},
		resp: resp,
	}
	x := <-resp

	return &x, nil
}

func (m *MTProto) ContactsGetTopPeers(correspondents, botsPM, botsInline, groups, channels bool, offset, limit, hash int32) (*TL, error) {
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
		return nil, errors.New("MTProto::ContactsGetTopPeers error: Unknown type")
	}

	return &x, nil
}
