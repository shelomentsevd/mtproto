package mtproto

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

type MTProto struct {
	f            *os.File
	queueSend    chan packetToSend
	stopRoutines chan struct{}
	allDone      sync.WaitGroup

	session ISession
	network INetwork

	appConfig Configuration

	dclist map[int32]string
}

type packetToSend struct {
	msg  TL
	resp chan response
}

type response struct {
	data TL
	err  error
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

func NewMTProto(newSession bool, serverAddr string, useIPv6 bool, authkeyfile string, appConfig Configuration) (*MTProto, error) {
	var err error

	err = appConfig.Check()
	if err != nil {
		return nil, err
	}

	m := new(MTProto)
	m.appConfig = appConfig

	m.queueSend = make(chan packetToSend, 64)
	m.stopRoutines = make(chan struct{})
	m.allDone = sync.WaitGroup{}

	m.f, err = os.OpenFile(authkeyfile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	m.session = NewSession(m.f)

	rand.Seed(time.Now().UnixNano())
	m.session.SetSessionID(rand.Int63())

	if newSession {
		m.session.SetAddress(serverAddr)
		m.session.UseIPv6(useIPv6)
		m.session.Encrypted(false)
		return m, nil
	}

	err = m.session.Load()
	if err == nil {
		m.session.Encrypted(true)
	} else {
		m.session.SetAddress(serverAddr)
		m.session.UseIPv6(useIPv6)
		m.session.Encrypted(false)
	}

	m.network = NewNetwork(m.session, m.queueSend)

	return m, nil
}

func (m *MTProto) Connect() (err error) {
	m.network.Connect()

	// start goroutines
	go m.sendRoutine()
	go m.readRoutine()

	var data *TL

	// (help_getConfig)
	if data, err = m.InvokeSync(TL_invokeWithLayer{
		Layer: layer,
		Query: TL_initConnection{
			Api_id:         m.appConfig.Id,
			Device_model:   m.appConfig.DeviceModel,
			System_version: m.appConfig.SystemVersion,
			App_version:    m.appConfig.Version,
			Lang_code:      m.appConfig.Language,
			Query:          TL_help_getConfig{},
		},
	}); err !=nil {
		return
	}

	switch (*data).(type) {
	case TL_config:
		m.dclist = make(map[int32]string, 5)
		for _, v := range (*data).(TL_config).Dc_options {
			v := v.(TL_dcOption)
			if m.session.IsIPv6() && v.Ipv6 {
				m.dclist[v.Id] = fmt.Sprintf("[%s]:%d", v.Ip_address, v.Port)
			} else if !v.Ipv6 {
				m.dclist[v.Id] = fmt.Sprintf("%s:%d", v.Ip_address, v.Port)
			}
		}
	default:
		err = fmt.Errorf("Connection error: got: %T", data)
	}

	// start keep alive ping
	go m.pingRoutine()

	return
}

func (m *MTProto) Disconnect() error {
	// stop ping, send and read routine by closing channel stopRoutines
	close(m.stopRoutines)

	// Wait until all goroutines stopped
	m.allDone.Wait()

	// close send queue
	close(m.queueSend)

	return m.network.Disconnect()
}

func (m *MTProto) reconnect(newaddr string) error {
	err := m.Disconnect()
	if err != nil {
		return err
	}

	// renew connection
	m.session.Encrypted(true)
	if newaddr != m.session.GetAddress() {
		m.session = NewSession(m.f)
		m.session.SetAddress(newaddr)
		m.session.Encrypted(false)
	}

	err = m.Connect()
	return err
}

func (m *MTProto) pingRoutine() {
	m.allDone.Add(1)
	defer func() { m.allDone.Done() }()
	for {
		select {
		case <-m.stopRoutines:
			return
		case <-time.After(60 * time.Second):
			m.InvokeAsync(TL_ping{0xCADACAD})
		}
	}
}

func (m *MTProto) sendRoutine() {
	m.allDone.Add(1)
	defer func() { m.allDone.Done() }()
	for {
		select {
		case <-m.stopRoutines:
			return
		case x := <-m.queueSend:
			err := m.network.Send(x.msg, x.resp)
			if err != nil {
				log.Fatalln("SendRoutine:", err)
			}
		}
	}
}

func (m *MTProto) readRoutine() {
	m.allDone.Add(1)
	defer func() { m.allDone.Done() }()
	for {
		// Run async wait for data from server
		ch := make(chan interface{}, 1)
		go func(ch chan<- interface{}) {
			data, err := m.network.Read()
			if err == io.EOF {
				// Connection closed by server, trying to reconnect
				err = m.reconnect(m.session.GetAddress())
				if err != nil {
					log.Fatalln("ReadRoutine: ", err)
				}
			}
			if err != nil {
				log.Fatalln("ReadRoutine: ", err)
			}
			ch <- data
		}(ch)

		select {
		case <-m.stopRoutines:
			return
		case data := <-ch:
			if data == nil {
				return
			}
			m.network.Process(data)
		}
	}
}

func (m *MTProto) handleRPCError(rpcError TL_rpc_error) error {
	switch rpcError.Error_code {
	case errorSeeOther:
		var newDc int32
		n, _ := fmt.Sscanf(rpcError.Error_message, "PHONE_MIGRATE_%d", &newDc)
		if n != 1 {
			n, _ := fmt.Sscanf(rpcError.Error_message, "NETWORK_MIGRATE_%d", &newDc)
			if n != 1 {
				return fmt.Errorf("RPC error_string: %s", rpcError.Error_message)
			}
		}
		newDcAddr, ok := m.dclist[newDc]
		if !ok {
			return fmt.Errorf("Wrong DC index: %d", newDc)
		}
		err := m.reconnect(newDcAddr)
		if err != nil {
			return err
		}
		return fmt.Errorf("mtproto error: %d %s", rpcError.Error_code, rpcError.Error_message)
	case errorBadRequest, errorUnauthorized, errorFlood, errorInternal:
		return fmt.Errorf("mtproto error: %d %s", rpcError.Error_code, rpcError.Error_message)
	default:
		return fmt.Errorf("mtproto unknow error: %d %s", rpcError.Error_code, rpcError.Error_message)
	}
}

func (m *MTProto) InvokeSync(msg TL) (*TL, error) {
	x := <-m.InvokeAsync(msg)

	if x.err != nil {
		return nil, x.err
	}

	return &x.data, nil
}

func (m *MTProto) InvokeAsync(msg TL) chan response {
	resp := make(chan response, 1)
	m.queueSend <- packetToSend{
		msg:  msg,
		resp: resp,
	}
	return resp
}
