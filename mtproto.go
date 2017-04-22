package mtproto

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
)

type MTProto struct {
	addr      string
	conn      *net.TCPConn
	f         *os.File
	queueSend chan packetToSend
	stopSend  chan struct{}
	stopRead  chan struct{}
	stopPing  chan struct{}
	allDone   chan struct{}

	authKey     []byte
	authKeyHash []byte
	serverSalt  []byte
	encrypted   bool
	sessionId   int64

	mutex        *sync.Mutex
	lastSeqNo    int32
	msgsIdToAck  map[int64]packetToSend
	msgsIdToResp map[int64]chan TL
	seqNo        int32
	msgId        int64

	appConfig *Configuration

	dclist map[int32]string
}

type packetToSend struct {
	msg  TL
	resp chan TL
}

type Configuration struct {
	Id            int32
	Hash          string
	Version       string
	DeviceModel   string
	SystemVersion string
	Language      string
}

// API Errors
const (
	errorSeeOther     = 303
	errorBadRequest   = 400
	errorUnauthorized = 401
	errorForbidden    = 403
	errorNotFound     = 404
	errorFlood        = 420
	errorInternal     = 500
)

const appConfigError = "App configuration error: %s"
const telegramAddr = "149.154.167.50:443"

// Current API Layer Version
const layer = 65

func NewConfiguration(id int32, hash, version, deviceModel, systemVersion, language string) (*Configuration, error) {
	appConfig := new(Configuration)

	if id == 0 || hash == "" || version == "" {
		return nil, fmt.Errorf(appConfigError, "Fields Id, Hash or Version are empty")
	}
	appConfig.Id = id
	appConfig.Hash = hash
	appConfig.Version = version

	appConfig.DeviceModel = deviceModel
	if deviceModel == "" {
		appConfig.DeviceModel = "Unknown"
	}

	appConfig.SystemVersion = systemVersion
	if systemVersion == "" {
		appConfig.SystemVersion = runtime.GOOS + "/" + runtime.GOARCH
	}

	appConfig.Language = language
	if language == "" {
		appConfig.Language = "en"
	}

	return appConfig, nil
}

func (appConfig Configuration) Check() error {
	if appConfig.Id == 0 || appConfig.Hash == "" || appConfig.Version == "" {
		return fmt.Errorf(appConfigError, "Configuration.Id, Configuration.Hash or Configuration.Version are empty")
	}

	if appConfig.DeviceModel == "" {
		return fmt.Errorf(appConfigError, "Configuration.DeviceModel is empty")
	}

	if appConfig.SystemVersion == "" {
		return fmt.Errorf(appConfigError, "Configuration.SystemVersion is empty")
	}

	if appConfig.Language == "" {
		return fmt.Errorf(appConfigError, "Configuration.Language is empty")
	}

	return nil
}

func NewMTProto(authkeyfile string, appConfig *Configuration) (*MTProto, error) {
	var err error

	err = appConfig.Check()
	if err != nil {
		return nil, err
	}

	m := new(MTProto)
	m.appConfig = appConfig

	m.f, err = os.OpenFile(authkeyfile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	err = m.readData()
	if err == nil {
		m.encrypted = true
	} else {
		m.addr = telegramAddr
		m.encrypted = false
	}
	rand.Seed(time.Now().UnixNano())
	m.sessionId = rand.Int63()

	return m, nil
}

func (m *MTProto) Connect() error {
	var err error
	var tcpAddr *net.TCPAddr

	// connect
	tcpAddr, err = net.ResolveTCPAddr("tcp", m.addr)
	if err != nil {
		return err
	}
	m.conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	// Packet Length is encoded by a single byte (see: https://core.telegram.org/mtproto)
	_, err = m.conn.Write([]byte{0xef})
	if err != nil {
		return err
	}
	// get new authKey if need
	if !m.encrypted {
		err = m.makeAuthKey()
		if err != nil {
			return err
		}
	}

	// start goroutines
	m.queueSend = make(chan packetToSend, 64)
	m.stopSend = make(chan struct{}, 1)
	m.stopRead = make(chan struct{}, 1)
	m.stopPing = make(chan struct{}, 1)
	m.allDone = make(chan struct{}, 3)
	m.msgsIdToAck = make(map[int64]packetToSend)
	m.msgsIdToResp = make(map[int64]chan TL)
	m.mutex = &sync.Mutex{}
	go m.sendRoutine()
	go m.readRoutine()

	// (help_getConfig)
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		msg: TL_invokeWithLayer{
			Layer: layer,
			Query: TL_initConnection{
				Api_id:         m.appConfig.Id,
				Device_model:   m.appConfig.DeviceModel,
				System_version: m.appConfig.SystemVersion,
				App_version:    m.appConfig.Version,
				Lang_code:      m.appConfig.Language,
				Query:          TL_help_getConfig{},
			},
		},
		resp: resp,
	}
	x := <-resp
	switch x.(type) {
	case TL_config:
		m.dclist = make(map[int32]string, 5)
		for _, v := range x.(TL_config).Dc_options {
			v := v.(TL_dcOption)
			m.dclist[v.Id] = fmt.Sprintf("%s:%d", v.Ip_address, v.Port)
		}
	default:
		return fmt.Errorf("Connection error: got: %T", x)
	}

	// start keep alive ping
	go m.pingRoutine()

	return nil
}

func (m *MTProto) reconnect(newaddr string) error {
	var err error

	// stop ping routine
	m.stopPing <- struct{}{}
	close(m.stopPing)

	// stop send routine
	m.stopSend <- struct{}{}
	close(m.stopSend)

	// stop read routine
	m.stopRead <- struct{}{}
	close(m.stopRead)

	<-m.allDone
	<-m.allDone
	<-m.allDone

	// close send queue
	close(m.queueSend)

	// close connection
	err = m.conn.Close()
	if err != nil {
		return err
	}

	// renew connection
	m.encrypted = false
	m.addr = newaddr
	err = m.Connect()
	return err
}

func (m *MTProto) AuthSendCode(phonenumber string) (error, *TL_auth_sentCode) {
	var authSentCode TL_auth_sentCode
	flag := true
	for flag {
		resp := make(chan TL, 1)
		m.queueSend <- packetToSend{
			msg: TL_auth_sendCode{
				Allow_flashcall: false,
				Phone_number:    phonenumber,
				Api_id:          m.appConfig.Id,
				Api_hash:        m.appConfig.Hash,
			},
			resp: resp,
		}
		x := <-resp
		switch x.(type) {
		case TL_auth_sentCode:
			authSentCode = x.(TL_auth_sentCode)
			flag = false
		case TL_rpc_error:
			x := x.(TL_rpc_error)
			if x.Error_code != errorSeeOther {
				return fmt.Errorf("RPC Error_code: %d", x.Error_code), nil
			}
			var newDc int32
			n, _ := fmt.Sscanf(x.Error_message, "PHONE_MIGRATE_%d", &newDc)
			if n != 1 {
				n, _ := fmt.Sscanf(x.Error_message, "NETWORK_MIGRATE_%d", &newDc)
				if n != 1 {
					return fmt.Errorf("RPC error_string: %s", x.Error_message), nil
				}
			}

			newDcAddr, ok := m.dclist[newDc]
			if !ok {
				return fmt.Errorf("Wrong DC index: %d", newDc), nil
			}
			err := m.reconnect(newDcAddr)
			if err != nil {
				return err, nil
			}
		default:
			return fmt.Errorf("Got: %T", x), nil
		}
	}

	return nil, &authSentCode
}

func (m *MTProto) AuthSignIn(phoneNumber, phoneCode, phoneCodeHash string) (error, *TL_auth_authorization) {
	if phoneNumber == "" || phoneCode == "" || phoneCodeHash == "" {
		return errors.New("MRProto::AuthSignIn one of function parameters is empty"), nil
	}

	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		msg: TL_auth_signIn{
			Phone_number:    phoneNumber,
			Phone_code_hash: phoneCodeHash,
			Phone_code:      phoneCode,
		},
		resp: resp,
	}
	x := <-resp
	auth, ok := x.(TL_auth_authorization)

	if !ok {
		return fmt.Errorf("RPC: %#v", x), nil
	}

	return nil, &auth
}

func (m *MTProto) GetTopPeers(correspondents, botsPM, botsInline, groups, channels bool, offset, limit, hash int32) (error, *TL) {
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
		return errors.New("MTProto::GetTopPeers error: Unknown type"), nil
	}

	return nil, &x
}

func (m *MTProto) GetHistory(peer TL, offsetId, offsetDate, addOffset, limit, maxId, minId int32) (error, *TL) {
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

func (m *MTProto) GetDialogs(excludePinned bool, offsetDate, offsetId int32, offsetPeer TL, limit int32) (error, *TL) {
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

func (m *MTProto) UpdatesGetState() (error, *TL) {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		msg: TL_updates_getState{},
		resp: resp,
	}
	x := <-resp

	return nil, &x
}

func (m *MTProto) UpdatesGetDifference(pts, ptsTotalLimit, date, qts int32) (error, *TL) {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		msg: TL_updates_getDifference{
			Pts: pts,
			Pts_total_limit: ptsTotalLimit,
			Date: date,
			Qts: qts,
		},
		resp: resp,
	}

	x := <-resp

	return nil, &x
}

func (m *MTProto) UpdatesGetChannelDifference(force bool, channel, filter TL, pts ,limit int32) (error, *TL) {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		msg: TL_updates_getChannelDifference{
			Force: force,
			Channel: channel,
			Filter: filter,
			Pts: pts,
			Limit: limit,
		},
		resp: resp,
	}

	x := <-resp

	return  nil, &x
}

func (m *MTProto) pingRoutine() {
	for {
		select {
		case <-m.stopPing:
			m.allDone <- struct{}{}
			return
		case <-time.After(60 * time.Second):
			m.queueSend <- packetToSend{TL_ping{0xCADACADA}, nil}
		}
	}
}

func (m *MTProto) sendRoutine() {
	for x := range m.queueSend {
		err := m.sendPacket(x.msg, x.resp)
		if err != nil {
			fmt.Println("SendRoutine:", err)
			os.Exit(2)
		}
	}

	m.allDone <- struct{}{}
}

func (m *MTProto) readRoutine() {
	for {
		data, err := m.read(m.stopRead)
		if err != nil {
			fmt.Println("ReadRoutine:", err)
			os.Exit(2)
		}
		if data == nil {
			m.allDone <- struct{}{}
			return
		}

		m.process(m.msgId, m.seqNo, data)
	}
}

func (m *MTProto) process(msgId int64, seqNo int32, data interface{}) interface{} {
	switch data.(type) {
	case TL_msg_container:
		data := data.(TL_msg_container).Items
		for _, v := range data {
			m.process(v.Msg_id, v.Seq_no, v.Data)
		}

	case TL_bad_server_salt:
		data := data.(TL_bad_server_salt)
		m.serverSalt = data.New_server_salt
		_ = m.saveData()
		m.mutex.Lock()
		for k, v := range m.msgsIdToAck {
			delete(m.msgsIdToAck, k)
			m.queueSend <- v
		}
		m.mutex.Unlock()

	case TL_new_session_created:
		data := data.(TL_new_session_created)
		m.serverSalt = data.Server_salt
		_ = m.saveData()

	case TL_ping:
		data := data.(TL_ping)
		m.queueSend <- packetToSend{TL_pong{msgId, data.Ping_id}, nil}

	case TL_pong:
		// ignore

	case TL_msgs_ack:
		data := data.(TL_msgs_ack)
		m.mutex.Lock()
		for _, v := range data.MsgIds {
			delete(m.msgsIdToAck, v)
		}
		m.mutex.Unlock()

	case TL_rpc_result:
		data := data.(TL_rpc_result)
		x := m.process(msgId, seqNo, data.Obj)
		m.mutex.Lock()
		v, ok := m.msgsIdToResp[data.Req_msg_id]
		if ok {
			v <- x.(TL)
			close(v)
			delete(m.msgsIdToResp, data.Req_msg_id)
		}
		delete(m.msgsIdToAck, data.Req_msg_id)
		m.mutex.Unlock()

	default:
		return data
	}

	// TODO: Check why I should do this
	if (seqNo & 1) == 1 {
		m.queueSend <- packetToSend{TL_msgs_ack{[]int64{msgId}}, nil}
	}

	return nil
}

// Save session
func (m *MTProto) saveData() (err error) {
	m.encrypted = true

	b := NewEncodeBuf(1024)
	b.StringBytes(m.authKey)
	b.StringBytes(m.authKeyHash)
	b.StringBytes(m.serverSalt)
	b.String(m.addr)

	err = m.f.Truncate(0)
	if err != nil {
		return err
	}

	_, err = m.f.WriteAt(b.buf, 0)
	if err != nil {
		return err
	}

	return nil
}

// Load session
func (m *MTProto) readData() (err error) {
	b := make([]byte, 1024*4)
	n, err := m.f.ReadAt(b, 0)
	if n <= 0 {
		return errors.New("New session")
	}

	d := NewDecodeBuf(b)
	m.authKey = d.StringBytes()
	m.authKeyHash = d.StringBytes()
	m.serverSalt = d.StringBytes()
	m.addr = d.String()

	if d.err != nil {
		return d.err
	}

	return nil
}
