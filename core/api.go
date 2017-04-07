package core

import "fmt"

const crc_gzip_packed = 0x3072cfa1 // Processed manually

const crc_boolFalse = 0xbc799737

type TL_boolFalse struct {
}

const crc_boolTrue = 0x997275b5

type TL_boolTrue struct {
}

const crc_error = 0xc4b9f9bb

type TL_error struct {
	code int32
	text string
}

func (e TL_error) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_error)
	x.Int(e.code)
	x.String(e.text)
	return x.buf
}

const crc_null = 0x56730bcc

type TL_null struct {
}

func (m *DecodeBuf) ObjectGenerated(constructor uint32) (r TL) {
	switch constructor {
	case crc_boolFalse:
		r = TL_boolFalse{}

	case crc_boolTrue:
		r = TL_boolTrue{}

	case crc_error:
		r = TL_error{
			m.Int(),
			m.String(),
		}

	case crc_null:
		r = TL_null{}
	case crc_config:
		flags := m.Int()
		phonecalls_enabled := false
		if flags & (1 << 1) {
			phonecalls_enabled = true
		}
		date := m.Int()
		expires := m.Int()
		test_mode := m.Object()
		this_dc := m.Int()
		dc_options := m.Vector()
		chat_size_max := m.Int()
		megagroup_size_max := m.Int()
		forwarded_count_max := m.Int()
		online_update_period_ms := m.Int()
		offline_blur_timeout_ms := m.Int()
		offline_idle_timeout_ms := m.Int()
		online_cloud_timeout_ms := m.Int()
		notify_cloud_delay_ms := m.Int()
		notify_default_delay_ms := m.Int()
		chat_big_size := m.Int()
		push_chat_period_ms := m.Int()
		push_chat_limit := m.Int()
		saved_gifs_limit := m.Int()
		edit_time_limit := m.Int()
		rating_e_decay := m.Int()
		stickers_recent_limit := m.Int()
		var tmp_sessions int32
		if flags & (1 << 0) {
			tmp_sessions = m.Int()
		}
		pinned_dialogs_count_max := m.Int()
		call_receive_timeout_ms := m.Int()
		call_ring_timeout_ms := m.Int()
		call_connect_timeout_ms := m.Int()
		call_packet_timeout_ms := m.Int()
		me_url_prefix := m.String()
		disabled_features := m.Vector()
		r = TL_config{
			flags:                    flags,
			phonecalls_enabled:       phonecalls_enabled,
			date:                     date,
			expires:                  expires,
			test_mode:                test_mode,
			this_dc:                  this_dc,
			dc_options:               dc_options,
			chat_size_max:            chat_size_max,
			megagroup_size_max:       megagroup_size_max,
			forwarded_count_max:      forwarded_count_max,
			online_update_period_ms:  online_update_period_ms,
			offline_blur_timeout_ms:  offline_blur_timeout_ms,
			offline_idle_timeout_ms:  offline_idle_timeout_ms,
			online_cloud_timeout_ms:  online_cloud_timeout_ms,
			notify_cloud_delay_ms:    notify_cloud_delay_ms,
			notify_default_delay_ms:  notify_default_delay_ms,
			chat_big_size:            chat_big_size,
			push_chat_period_ms:      push_chat_period_ms,
			push_chat_limit:          push_chat_limit,
			saved_gifs_limit:         saved_gifs_limit,
			edit_time_limit:          edit_time_limit,
			rating_e_decay:           rating_e_decay,
			stickers_recent_limit:    stickers_recent_limit,
			tmp_sessions:             tmp_sessions,
			pinned_dialogs_count_max: pinned_dialogs_count_max,
			call_receive_timeout_ms:  call_receive_timeout_ms,
			call_ring_timeout_ms:     call_ring_timeout_ms,
			call_connect_timeout_ms:  call_connect_timeout_ms,
			call_packet_timeout_ms:   call_packet_timeout_ms,
			me_url_prefix:            me_url_prefix,
			disabled_features:        disabled_features,
		}
	case crc_dcOption:
		flags := m.Int()
		var ipv6, media_only, tcpo_only bool
		if flags & (1 << 0) {
			ipv6 = true
		}
		if flags & (1 << 1) {
			media_only = true
		}
		if flags & (1 << 2) {
			tcpo_only = true
		}
		id := m.Int()
		ip_address := m.String()
		port := m.Int()
		r = TL_dcOption{
			flags:      flags,
			ipv6:       ipv6,
			media_only: media_only,
			tcpo_only:  tcpo_only,
			id:         id,
			ip_address: ip_address,
			port:       port,
		}
	default:
		m.err = fmt.Errorf("Unknown constructor: \u002508x", constructor)
		return nil
	}
	return
}

// invokeWithLayer#da9b0d0d {X:Type} layer:int query:!X = X;
const crc_invokeWithLayer = 0xda9b0d0d

type TL_invokeWithLayer struct {
	layer int32
	query TL
}

func (e TL_invokeWithLayer) encode() []byte {
	// TODO: 512 is a magic number
	x := NewEncodeBuf(512)
	x.UInt(crc_invokeWithLayer)
	x.Int(e.layer)
	x.Bytes(e.query.encode())
	// TODO: Should I shrink a buffer to his actual size or not?
	return x.buf
}

// initConnection#69796de9 {X:Type} api_id:int device_model:string system_version:string app_version:string lang_code:string query:!X = X;
const crc_initConnection = 0x69796de9

type TL_initConnection struct {
	api_id         int32
	device_model   string
	system_version string
	app_version    string
	lang_code      string
	query          TL
}

func (e TL_initConnection) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_initConnection)
	x.Int(e.api_id)
	x.String(e.device_model)
	x.String(e.system_version)
	x.String(e.app_version)
	x.String(e.lang_code)
	x.Bytes(e.query.encode())
	return x.buf
}

// help.getConfig#c4f9186b = Config;
const crc_help_getConfig = 0xc4f9186b

type TL_help_getConfig struct {
}

func (e TL_help_getConfig) encode() []byte {
	x := NewEncodeBuf(4)
	x.UInt(crc_help_getConfig)
	return x.buf
}

// config#cb601684 flags:# phonecalls_enabled:flags.1?true date:int expires:int test_mode:Bool this_dc:int dc_options:Vector<DcOption> chat_size_max:int megagroup_size_max:int forwarded_count_max:int online_update_period_ms:int offline_blur_timeout_ms:int offline_idle_timeout_ms:int online_cloud_timeout_ms:int notify_cloud_delay_ms:int notify_default_delay_ms:int chat_big_size:int push_chat_period_ms:int push_chat_limit:int saved_gifs_limit:int edit_time_limit:int rating_e_decay:int stickers_recent_limit:int tmp_sessions:flags.0?int pinned_dialogs_count_max:int call_receive_timeout_ms:int call_ring_timeout_ms:int call_connect_timeout_ms:int call_packet_timeout_ms:int me_url_prefix:string disabled_features:Vector<DisabledFeature> = Config;
const crc_config = 0xcb601684

type TL_config struct {
	flags                    int32
	phonecalls_enabled       bool // flags.1?true TODO: TL_true
	date                     int32
	expires                  int32
	test_mode                TL // TL_boolFalse or TL_boolTrue
	this_dc                  int32
	dc_options               []TL // DcOption
	chat_size_max            int32
	megagroup_size_max       int32
	forwarded_count_max      int32
	online_update_period_ms  int32
	offline_blur_timeout_ms  int32
	offline_idle_timeout_ms  int32
	online_cloud_timeout_ms  int32
	notify_cloud_delay_ms    int32
	notify_default_delay_ms  int32
	chat_big_size            int32
	push_chat_period_ms      int32
	push_chat_limit          int32
	saved_gifs_limit         int32
	edit_time_limit          int32
	rating_e_decay           int32
	stickers_recent_limit    int32
	tmp_sessions             int32 // flags.0?int
	pinned_dialogs_count_max int32
	call_receive_timeout_ms  int32
	call_ring_timeout_ms     int32
	call_connect_timeout_ms  int32
	call_packet_timeout_ms   int32
	me_url_prefix            string
	disabled_features        []TL // DisabledFeature
}

func (e TL_config) encode() []byte { return nil }

// dcOption#5d8c6cc flags:# ipv6:flags.0?true media_only:flags.1?true tcpo_only:flags.2?true id:int ip_address:string port:int = DcOption;
const crc_dcOption = 0x5d8c6cc

type TL_dcOption struct {
	flags      int32
	ipv6       bool // ipv6:flags.0?true TODO: TL_true
	media_only bool // media_only:flags.1?true TODO: TL_true
	tcpo_only  bool // tcpo_only:flags.2?true TODO: TL_true
	id         int32
	ip_address string
	port       int32
}

func (e TL_dcOption) encode() []byte { return nil }
