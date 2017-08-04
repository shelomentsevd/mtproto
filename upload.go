package mtproto

func (m *MTProto) UploadGetFile (location TL_inputFileLocation, offset int32, limit int32) (*TL, error) {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		msg: TL_upload_getFile{
			Location: location,
			Offset: offset,
			Limit: limit,
		},
		resp: resp,
	}
	x := <-resp

	return &x, nil
}
