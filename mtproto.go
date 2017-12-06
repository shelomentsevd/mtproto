package mtproto

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

type MTProto struct {
	queueSend    chan packetToSend
	stopRoutines chan struct{}
	allDone      sync.WaitGroup

	network INetwork

	IPv6        bool
	authkeyfile string
	id          int32
	hash        string
	version     string
	device      string
	system      string
	language    string

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

type Option func(*options)

type options struct {
	Version       string
	DeviceModel   string
	SystemVersion string
	Language      string
	IPv6          bool
	AuthkeyFile   string
	ServerAddress string
	NewSession    bool
}

func WithVersion(version string) Option {
	return func(opts *options) {
		opts.Version = version
	}
}

func WithDevice(device string) Option {
	return func(opts *options) {
		opts.DeviceModel = device
	}
}

func WithSystem(system string) Option {
	return func(opts *options) {
		opts.SystemVersion = system
	}
}

func WithLanguage(language string) Option {
	return func(opts *options) {
		opts.Language = language
	}
}

func WithServer(server string, ipv6 bool) Option {
	return func(opts *options) {
		opts.ServerAddress = server
		opts.IPv6 = ipv6
	}
}

func WithAuthFile(authfile string, newSession bool) Option {
	return func(opts *options) {
		opts.AuthkeyFile = authfile
		opts.NewSession = newSession
	}
}

var defaultOptions = options{
	DeviceModel:   "Unknown",
	SystemVersion: runtime.GOOS + "/" + runtime.GOARCH,
	Language:      "en",
	IPv6:          false,
	AuthkeyFile:   os.Getenv("HOME") + "/mtproto.auth",
	ServerAddress: "149.154.167.50:443",
	Version:       "0.0.1",
	NewSession:    false,
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

// Current API Layer Version
const layer = 65

func NewMTProto(id int32, hash string, opts ...Option) (*MTProto, error) {
	var err error

	if id == 0 {
		return nil, fmt.Errorf("can't initialize mtproto: wrong application id")
	}

	if len(hash) == 0 {
		return nil, fmt.Errorf("can't initialize mtpoto: wrong application hash")
	}

	configuration := defaultOptions
	for _, option := range opts {
		option(&configuration)
	}

	m := new(MTProto)

	m.queueSend = make(chan packetToSend, 64)
	m.stopRoutines = make(chan struct{})
	m.allDone = sync.WaitGroup{}

	m.id = id
	m.hash = hash
	m.version = configuration.Version
	m.device = configuration.DeviceModel
	m.system = configuration.SystemVersion
	m.language = configuration.Language
	m.authkeyfile = configuration.AuthkeyFile
	m.IPv6 = configuration.IPv6

	if m.network, err = NewNetwork(configuration.NewSession, m.authkeyfile, m.queueSend, configuration.ServerAddress, m.IPv6); err != nil {
		return nil, err
	}

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
			Api_id:         m.id,
			Device_model:   m.device,
			System_version: m.system,
			App_version:    m.version,
			Lang_code:      m.language,
			Query:          TL_help_getConfig{},
		},
	}); err != nil {
		return
	}

	switch (*data).(type) {
	case TL_config:
		m.dclist = make(map[int32]string, 5)
		for _, v := range (*data).(TL_config).Dc_options {
			v := v.(TL_dcOption)
			if m.IPv6 && v.Ipv6 {
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
	if newaddr != m.network.Address() {
		m.network, err = NewNetwork(true, m.authkeyfile, m.queueSend, newaddr, m.IPv6)
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
			// TODO: m.InvokeSync()?
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
				// TODO: Last message to the server was lost. Fix it.
				// Connection closed by server, trying to reconnect
				err = m.reconnect(m.network.Address())
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

func (m *MTProto) InvokeSync(msg TL) (*TL, error) {
	x := <-m.InvokeAsync(msg)

	if x.err != nil {
		if err, ok := x.err.(TL_rpc_error); ok {
			switch err.Error_code {
			case errorSeeOther:
				var newDc int32
				n, _ := fmt.Sscanf(err.Error_message, "PHONE_MIGRATE_%d", &newDc)
				if n != 1 {
					n, _ := fmt.Sscanf(err.Error_message, "NETWORK_MIGRATE_%d", &newDc)
					if n != 1 {
						return nil, fmt.Errorf("RPC error_string: %s", err.Error_message)
					}
				}
				newDcAddr, ok := m.dclist[newDc]
				if !ok {
					return nil, fmt.Errorf("wrong DC index: %d", newDc)
				}
				err := m.reconnect(newDcAddr)
				if err != nil {
					return nil, err
				}
			default:
				return nil, x.err
			}
		}

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
