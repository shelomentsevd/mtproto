package mtproto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// low-level communication with telegram server
type INetwork interface {
	Connect() error
	Disconnect() error
	Send(msg TL, resp chan response) error
	Read() (interface{}, error)
	Process(data interface{}) interface{}
}

type Network struct {
	session ISession

	conn *net.TCPConn

	mutex        *sync.Mutex
	msgsIdToAck  map[int64]packetToSend
	msgsIdToResp map[int64]chan response

	queueSend chan packetToSend
	lastSeqNo int32
	seqNo     int32
	msgId     int64
}

func NewNetwork(session ISession, queueSend chan packetToSend) INetwork {
	nw := new(Network)

	nw.session = session
	nw.queueSend = queueSend
	nw.msgsIdToAck = make(map[int64]packetToSend)
	nw.msgsIdToResp = make(map[int64]chan response)
	nw.mutex = &sync.Mutex{}

	return nw
}

func (nw *Network) Connect() error {
	var err error
	var tcpAddr *net.TCPAddr

	// connect
	tcpAddr, err = net.ResolveTCPAddr("tcp", nw.session.GetAddress())
	if err != nil {
		return err
	}
	nw.conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	// Packet Length is encoded by a single byte (see: https://core.telegram.org/mtproto)
	_, err = nw.conn.Write([]byte{0xef})
	if err != nil {
		return err
	}
	// get new authKey if need
	if !nw.session.IsEncrypted() {
		err = nw.makeAuthKey()
		if err != nil {
			return err
		}
	}

	return nil
}

func (nw *Network) Disconnect() error {

	return nw.conn.Close()
}

func (nw *Network) Send(msg TL, resp chan response) error {
	obj := msg.encode()

	x := NewEncodeBuf(256)

	// padding for tcpsize
	x.Int(0)

	if nw.session.IsEncrypted() {
		needAck := true
		switch msg.(type) {
		case TL_ping, TL_msgs_ack:
			needAck = false
		}
		z := NewEncodeBuf(256)
		newMsgId := GenerateMessageId()
		z.Bytes(nw.session.GetServerSalt())
		z.Long(nw.session.GetSessionID())
		z.Long(newMsgId)
		if needAck {
			z.Int(nw.lastSeqNo | 1)
		} else {
			z.Int(nw.lastSeqNo)
		}
		z.Int(int32(len(obj)))
		z.Bytes(obj)

		msgKey := sha1(z.buf)[4:20]
		aesKey, aesIV := generateAES(msgKey, nw.session.GetAuthKey(), false)

		y := make([]byte, len(z.buf)+((16-(len(obj)%16))&15))
		copy(y, z.buf)
		encryptedData, err := doAES256IGEencrypt(y, aesKey, aesIV)
		if err != nil {
			return err
		}

		nw.lastSeqNo += 2
		if needAck {
			nw.mutex.Lock()
			nw.msgsIdToAck[newMsgId] = packetToSend{msg, resp}
			nw.mutex.Unlock()
		}

		x.Bytes(nw.session.GetAuthKeyHash())
		x.Bytes(msgKey)
		x.Bytes(encryptedData)

		if resp != nil {
			nw.mutex.Lock()
			nw.msgsIdToResp[newMsgId] = resp
			nw.mutex.Unlock()
		}

	} else {
		x.Long(0)
		x.Long(GenerateMessageId())
		x.Int(int32(len(obj)))
		x.Bytes(obj)

	}

	// minus padding
	size := len(x.buf)/4 - 1

	if size < 127 {
		x.buf[3] = byte(size)
		x.buf = x.buf[3:]
	} else {
		binary.LittleEndian.PutUint32(x.buf, uint32(size<<8|127))
	}
	_, err := nw.conn.Write(x.buf)
	if err != nil {
		return err
	}

	return nil
}

func (nw *Network) Read() (interface{}, error) {
	var err error
	var n int
	var size int
	var data interface{}

	err = nw.conn.SetReadDeadline(time.Now().Add(300 * time.Second))
	if err != nil {
		return nil, err
	}
	b := make([]byte, 1)
	n, err = nw.conn.Read(b)
	if err != nil {
		return nil, err
	}

	if b[0] < 127 {
		size = int(b[0]) << 2
	} else {
		b := make([]byte, 3)
		n, err = nw.conn.Read(b)
		if err != nil {
			return nil, err
		}
		size = (int(b[0]) | int(b[1])<<8 | int(b[2])<<16) << 2
	}

	left := size
	buf := make([]byte, size)
	for left > 0 {
		n, err = nw.conn.Read(buf[size-left:])
		if err != nil {
			return nil, err
		}
		left -= n
	}

	if size == 4 {
		return nil, fmt.Errorf("Server response error: %d", int32(binary.LittleEndian.Uint32(buf)))
	}

	dbuf := NewDecodeBuf(buf)

	authKeyHash := dbuf.Bytes(8)
	if binary.LittleEndian.Uint64(authKeyHash) == 0 {
		nw.msgId = dbuf.Long()
		messageLen := dbuf.Int()
		if int(messageLen) != dbuf.size-20 {
			return nil, fmt.Errorf("Message len: %d (need %d)", messageLen, dbuf.size-20)
		}
		nw.seqNo = 0

		data = dbuf.Object()
		if dbuf.err != nil {
			return nil, dbuf.err
		}

	} else {
		msgKey := dbuf.Bytes(16)
		encryptedData := dbuf.Bytes(dbuf.size - 24)
		aesKey, aesIV := generateAES(msgKey, nw.session.GetAuthKey(), true)
		x, err := doAES256IGEdecrypt(encryptedData, aesKey, aesIV)
		if err != nil {
			return nil, err
		}
		dbuf = NewDecodeBuf(x)
		_ = dbuf.Long() // salt
		_ = dbuf.Long() // session_id
		nw.msgId = dbuf.Long()
		nw.seqNo = dbuf.Int()
		messageLen := dbuf.Int()
		if int(messageLen) > dbuf.size-32 {
			return nil, fmt.Errorf("Message len: %d (need less than %d)", messageLen, dbuf.size-32)
		}
		if !bytes.Equal(sha1(dbuf.buf[0 : 32+messageLen])[4:20], msgKey) {
			return nil, errors.New("Wrong msg_key")
		}

		data = dbuf.Object()
		if dbuf.err != nil {
			return nil, dbuf.err
		}

	}
	mod := nw.msgId & 3
	if mod != 1 && mod != 3 {
		return nil, fmt.Errorf("Wrong bits of message_id: %d", mod)
	}

	return data, nil
}

func (nw *Network) makeAuthKey() error {
	var x []byte
	var err error
	var data interface{}

	// (send) req_pq
	nonceFirst := GenerateNonce(16)
	err = nw.Send(TL_req_pq{nonceFirst}, nil)
	if err != nil {
		return err
	}

	// (parse) resPQ
	data, err = nw.Read()
	if err != nil {
		return err
	}
	res, ok := data.(TL_resPQ)
	if !ok {
		return errors.New("Handshake: Need resPQ")
	}
	if !bytes.Equal(nonceFirst, res.Nonce) {
		return errors.New("Handshake: Wrong Nonce")
	}
	found := false
	for _, b := range res.Fingerprints {
		if uint64(b) == telegramPublicKey_FP {
			found = true
			break
		}
	}
	if !found {
		return errors.New("Handshake: No fingerprint")
	}

	// (encoding) p_q_inner_data
	p, q := splitPQ(res.Pq)
	nonceSecond := GenerateNonce(32)
	nonceServer := res.Server_nonce
	innerData1 := (TL_p_q_inner_data{res.Pq, p, q, nonceFirst, nonceServer, nonceSecond}).encode()

	x = make([]byte, 255)
	copy(x[0:], sha1(innerData1))
	copy(x[20:], innerData1)
	encryptedData1 := doRSAencrypt(x)
	// (send) req_DH_params
	err = nw.Send(TL_req_DH_params{nonceFirst, nonceServer, p, q, telegramPublicKey_FP, encryptedData1}, nil)
	if err != nil {
		return err
	}

	// (parse) server_DH_params_{ok, fail}
	data, err = nw.Read()
	if err != nil {
		return err
	}
	dh, ok := data.(TL_server_DH_params_ok)
	if !ok {
		return errors.New("Handshake: Need server_DH_params_ok")
	}
	if !bytes.Equal(nonceFirst, dh.Nonce) {
		return errors.New("Handshake: Wrong Nonce")
	}
	if !bytes.Equal(nonceServer, dh.Server_nonce) {
		return errors.New("Handshake: Wrong Server_nonce")
	}
	t1 := make([]byte, 48)
	copy(t1[0:], nonceSecond)
	copy(t1[32:], nonceServer)
	hash1 := sha1(t1)

	t2 := make([]byte, 48)
	copy(t2[0:], nonceServer)
	copy(t2[16:], nonceSecond)
	hash2 := sha1(t2)

	t3 := make([]byte, 64)
	copy(t3[0:], nonceSecond)
	copy(t3[32:], nonceSecond)
	hash3 := sha1(t3)

	tmpAESKey := make([]byte, 32)
	tmpAESIV := make([]byte, 32)

	copy(tmpAESKey[0:], hash1)
	copy(tmpAESKey[20:], hash2[0:12])

	copy(tmpAESIV[0:], hash2[12:20])
	copy(tmpAESIV[8:], hash3)
	copy(tmpAESIV[28:], nonceSecond[0:4])

	// (parse-thru) server_DH_inner_data
	decodedData, err := doAES256IGEdecrypt(dh.Encrypted_answer, tmpAESKey, tmpAESIV)
	if err != nil {
		return err
	}
	innerbuf := NewDecodeBuf(decodedData[20:])
	data = innerbuf.Object()
	if innerbuf.err != nil {
		return innerbuf.err
	}
	dhi, ok := data.(TL_server_DH_inner_data)
	if !ok {
		return errors.New("Handshake: Need server_DH_inner_data")
	}
	if !bytes.Equal(nonceFirst, dhi.Nonce) {
		return errors.New("Handshake: Wrong Nonce")
	}
	if !bytes.Equal(nonceServer, dhi.Server_nonce) {
		return errors.New("Handshake: Wrong Server_nonce")
	}

	_, g_b, g_ab := makeGAB(dhi.G, dhi.G_a, dhi.Dh_prime)
	authKey := g_ab.Bytes()
	if authKey[0] == 0 {
		authKey = authKey[1:]
	}
	authKeyHash := sha1(authKey)[12:20]
	t4 := make([]byte, 32+1+8)
	copy(t4[0:], nonceSecond)
	t4[32] = 1
	copy(t4[33:], sha1(authKey)[0:8])
	nonceHash1 := sha1(t4)[4:20]
	serverSalt := make([]byte, 8)
	copy(serverSalt, nonceSecond[:8])
	xor(serverSalt, nonceServer[:8])

	nw.session.SetAuthKey(authKey)
	nw.session.SetAuthKeyHash(authKeyHash)
	nw.session.SetServerSalt(serverSalt)

	// (encoding) client_DH_inner_data
	innerData2 := (TL_client_DH_inner_data{nonceFirst, nonceServer, 0, g_b}).encode()
	x = make([]byte, 20+len(innerData2)+(16-((20+len(innerData2))%16))&15)
	copy(x[0:], sha1(innerData2))
	copy(x[20:], innerData2)
	encryptedData2, err := doAES256IGEencrypt(x, tmpAESKey, tmpAESIV)

	// (send) set_client_DH_params
	err = nw.Send(TL_set_client_DH_params{nonceFirst, nonceServer, encryptedData2}, nil)
	if err != nil {
		return err
	}

	// (parse) dh_gen_{ok, Retry, fail}
	data, err = nw.Read()
	if err != nil {
		return err
	}
	dhg, ok := data.(TL_dh_gen_ok)
	if !ok {
		return errors.New("Handshake: Need dh_gen_ok")
	}
	if !bytes.Equal(nonceFirst, dhg.Nonce) {
		return errors.New("Handshake: Wrong Nonce")
	}
	if !bytes.Equal(nonceServer, dhg.Server_nonce) {
		return errors.New("Handshake: Wrong Server_nonce")
	}
	if !bytes.Equal(nonceHash1, dhg.New_nonce_hash1) {
		return errors.New("Handshake: Wrong New_nonce_hash1")
	}

	// (all ok)
	err = nw.session.Save()
	if err != nil {
		return err
	}
	nw.session.Encrypted(true)

	return nil
}

func (nw *Network) Process(data interface{}) interface{} {
	return nw.process(nw.msgId, nw.seqNo, data)
}

func (nw *Network) process(msgId int64, seqNo int32, data interface{}) interface{} {
	switch data.(type) {
	case TL_msg_container:
		data := data.(TL_msg_container).Items
		for _, v := range data {
			nw.process(v.Msg_id, v.Seq_no, v.Data)
		}

	case TL_bad_server_salt:
		data := data.(TL_bad_server_salt)
		nw.session.SetServerSalt(data.New_server_salt)
		_ = nw.session.Save()
		nw.mutex.Lock()
		defer nw.mutex.Unlock()
		for k, v := range nw.msgsIdToAck {
			delete(nw.msgsIdToAck, k)
			nw.queueSend <- v
		}

	case TL_new_session_created:
		data := data.(TL_new_session_created)
		nw.session.SetServerSalt(data.Server_salt)
		_ = nw.session.Save()

	case TL_ping:
		data := data.(TL_ping)
		nw.queueSend <- packetToSend{TL_pong{msgId, data.Ping_id}, nil}

	case TL_pong:
		// ignore

	case TL_msgs_ack:
		data := data.(TL_msgs_ack)
		nw.mutex.Lock()
		defer nw.mutex.Unlock()
		for _, v := range data.MsgIds {
			delete(nw.msgsIdToAck, v)
		}

	case TL_rpc_result:
		data := data.(TL_rpc_result)
		x := nw.process(msgId, seqNo, data.Obj)
		nw.mutex.Lock()
		defer nw.mutex.Unlock()
		if v, ok := nw.msgsIdToResp[data.Req_msg_id]; ok {
			// TODO: response struct is useless
			v <- response{x.(TL), nil}
			close(v)
		}
		delete(nw.msgsIdToAck, data.Req_msg_id)
	default:
		return data
	}

	// TODO: Check why I should do this
	if (seqNo & 1) == 1 {
		nw.queueSend <- packetToSend{TL_msgs_ack{[]int64{msgId}}, nil}
	}

	return nil
}
