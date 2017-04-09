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
		if flags & (1 << 1) != 0 {
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
		if flags & (1 << 0) != 0 {
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
		if flags & (1 << 0) != 0 {
			ipv6 = true
		}
		if flags & (1 << 1) != 0 {
			media_only = true
		}
		if flags & (1 << 2) != 0 {
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
	case crc_auth_codeTypeSms:
		// TODO
	case crc_auth_codeTypeCall:
		// TODO
	case crc_auth_codeTypeFlashCall:
		// TODO
	case crc_auth_sentCodeTypeApp:
		// TODO
	case crc_auth_sentCodeTypeSms:
		// TODO
	case crc_auth_sentCodeTypeCall:
		// TODO
	case crc_auth_sentCodeTypeFlashCall:
		// TODO
	case crc_auth_sentCode:
		// TODO
	case crc_auth_sendCode:
		// TODO
	case crc_auth_signIn:
		// TODO
	case crc_auth_authorization:
		// TODO
	case crc_userEmpty:
		// TODO
	case crc_user:
		// TODO
	case crc_userProfilePhotoEmpty:
		// TODO
	case crc_userProfilePhoto:
		// TODO
	case crc_fileLocationUnavailable:
		// TODO
	case crc_fileLocation:
		// TODO
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

//auth.codeTypeSms#72a3158c = auth.CodeType;
const crc_auth_codeTypeSms = 0x72a3158c
type TL_auth_codeTypeSms struct {}

func (e TL_auth_codeTypeSms) encode() []byte {
	x := NewEncodeBuf(4)
	x.UInt(crc_auth_codeTypeSms)
	return x.buf
}
//auth.codeTypeCall#741cd3e3 = auth.CodeType;
const crc_auth_codeTypeCall = 0x741cd3e3
type TL_auth_codeTypeCall struct {}

func (e TL_auth_codeTypeCall) encode() []byte  {
	x := NewEncodeBuf(4)
	x.UInt(crc_auth_codeTypeCall)
	return x.buf
}
//auth.codeTypeFlashCall#226ccefb = auth.CodeType;
const crc_auth_codeTypeFlashCall = 0x226ccefb
type TL_auth_codeTypeFlashCall struct {}

func (e TL_auth_codeTypeFlashCall) encode() []byte {
	x := NewEncodeBuf(4)
	x.UInt(crc_auth_codeTypeFlashCall)
	return x.buf
}

//auth.sentCodeTypeApp#3dbb5986 length:int = auth.SentCodeType;
const crc_auth_sentCodeTypeApp = 0x3dbb5986
type TL_auth_sentCodeTypeApp struct {}

func (e TL_auth_sentCodeTypeApp) encode() []byte  {
	x := NewEncodeBuf(4)
	x.UInt(crc_auth_sentCodeTypeApp)
	return x.buf
}

//auth.sentCodeTypeSms#c000bba2 length:int = auth.SentCodeType;
const crc_auth_sentCodeTypeSms = 0xc000bba2
type TL_auth_sentCodeTypeSms struct {}

func (e TL_auth_sentCodeTypeSms) encode() []byte  {
	x := NewEncodeBuf(4)
	x.UInt(crc_auth_sentCodeTypeSms)
	return x.buf
}

//auth.sentCodeTypeCall#5353e5a7 length:int = auth.SentCodeType;
const crc_auth_sentCodeTypeCall = 0x5353e5a7
type TL_auth_sentCodeTypeCall struct {}

func (e TL_auth_sentCodeTypeCall) encode() []byte {
	x := NewEncodeBuf(4)
	x.UInt(crc_auth_sentCodeTypeCall)
	return x.buf
}

//auth.sentCodeTypeFlashCall#ab03c6d9 pattern:string = auth.SentCodeType;
const crc_auth_sentCodeTypeFlashCall = 0xab03c6d9
type TL_auth_sentCodeTypeFlashCall struct {}

func (e TL_auth_sentCodeTypeFlashCall) encode() []byte  {
	x := NewEncodeBuf(4)
	x.UInt(crc_auth_sentCodeTypeFlashCall)
	return x.buf
}

// auth.sentCode#5e002502 flags:# phone_registered:flags.0?true type:auth.SentCodeType phone_code_hash:string next_type:flags.1?auth.CodeType timeout:flags.2?int = auth.SentCode;
const crc_auth_sentCode = 0x5e002502
type TL_auth_sentCode struct {
	flags int32
	phone_registered bool
	code_type TL // type:auth.SentCodeType
	phone_code_hash string
	next_type TL
	timeout int32
}

func (e TL_auth_sentCode) encode() []byte {
	var flags int32
	x := NewEncodeBuf(512)
	x.UInt(crc_auth_sentCode)
	// fill bits in flags
	if e.phone_registered {
		flags |= (1 << 0)
	}
	if _, ok := (e.next_type).(TL_null); !ok {
		flags |= (1 << 1)
	}
	if e.timeout > 0 {
		flags |= (1 << 2)
	}
	x.Int(flags)
	x.Bytes(e.code_type.encode())
	x.String(e.phone_code_hash)
	if _, ok := (e.next_type).(TL_null); !ok {
		x.Bytes(e.next_type.encode())
	}
	if e.timeout > 0 {
		x.Int(e.timeout)
	}
	return x.buf
}

// auth.sendCode#86aef0ec flags:# allow_flashcall:flags.0?true phone_number:string current_number:flags.0?Bool api_id:int api_hash:string = auth.SentCode;
const crc_auth_sendCode = 0x86aef0ec

type TL_auth_sendCode struct {
	flags int32
	allow_flashcall bool // allow_flashcall:flags.0?true
	phone_number string
	current_number TL // current_number:flags.0?Bool
	api_id int32
	api_hash string
}

func (e TL_auth_sendCode) encode() []byte {
	var flags int32
	x := NewEncodeBuf(512)
	x.UInt(crc_auth_sendCode)
	if e.allow_flashcall {
		flags |= (1 << 0)
	}
	x.Int(flags)
	x.String(e.phone_number)
	if e.allow_flashcall {
		x.Bytes(e.current_number.encode())
	}
	x.Int(e.api_id)
	x.String(e.api_hash)
	return x.buf
}

// auth.signIn#bcd51581 phone_number:string phone_code_hash:string phone_code:string = auth.Authorization;
const crc_auth_signIn = 0xbcd51581
type TL_auth_signIn struct {
	phone_number string
	phone_code_hash string
	phone_code string
}

func (e TL_auth_signIn) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_auth_signIn)
	x.String(e.phone_number)
	x.String(e.phone_code_hash)
	x.String(e.phone_code)
	return x.buf
}

// auth.authorization#cd050916 flags:# tmp_sessions:flags.0?int user:User = auth.Authorization;
const crc_auth_authorization = 0xcd050916
type TL_auth_authorization struct {
	flags int32
	tmp_sessions int32
	user TL
}

func (e TL_auth_authorization) encode() []byte {
	var flags int32
	x := NewEncodeBuf(512)
	x.UInt(crc_auth_authorization)
	// TODO: I am not sure about this condition. Check how serialization works in other libraries
	if e.tmp_sessions > 0 {
		flags |= (1 << 0)
	}
	x.Int(flags)
	if e.tmp_sessions > 0 {
		x.Int(e.tmp_sessions)
	}
	x.Bytes(e.user.encode())
	return x.buf
}

//fileLocationUnavailable#7c596b46 volume_id:long local_id:int secret:long = FileLocation;
const crc_fileLocationUnavailable = 0x7c596b46

type TL_fileLocationUnavailable struct {
	volume_id int64
	local_id int32
	secret int64
}

func (e TL_fileLocationUnavailable) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_fileLocationUnavailable)
	x.Long(e.volume_id)
	x.Int(e.local_id)
	x.Long(e.secret)
	return x.buf
}
//fileLocation#53d69076 dc_id:int volume_id:long local_id:int secret:long = FileLocation;
const crc_fileLocation = 0x53d69076
type TL_fileLocation struct {
	dc_id int32
	volume_id int64
	local_id int32
	secret int64
}

func (e TL_fileLocation) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_fileLocation)
	x.Int(e.dc_id)
	x.Long(e.volume_id)
	x.Int(e.local_id)
	x.Long(e.secret)
	return x.buf
}
//userProfilePhotoEmpty#4f11bae1 = UserProfilePhoto;
const crc_userProfilePhotoEmpty = 0x4f11bae1
type TL_userProfilePhotoEmpty struct {}

func (e TL_userProfilePhotoEmpty) encode() []byte  {
	x := NewEncodeBuf(4)
	x.UInt(crc_userProfilePhotoEmpty)
	return x.buf
}

//userProfilePhoto#d559d8c8 photo_id:long photo_small:FileLocation photo_big:FileLocation = UserProfilePhoto;
const crc_userProfilePhoto = 0xd559d8c8
type TL_userProfilePhoto struct {
	photo_id int64
	photo_small TL // FileLocation
	photo_big TL // FileLocation
}

func (e TL_userProfilePhoto) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_userProfilePhoto)
	x.Long(e.photo_id)
	x.Bytes(e.photo_small.encode())
	x.Bytes(e.photo_big.encode())
	return x.buf
}

//userEmpty#200250ba id:int = User;
const crc_userEmpty = 0x200250ba
type TL_userEmpty struct {
	id int32
}

func (e TL_userEmpty) encode() []byte {
	x := NewEncodeBuf(8)
	x.UInt(crc_userEmpty)
	x.Int(e.id)

	return x.buf
}

//user#d10d979a flags:# self:flags.10?true contact:flags.11?true mutual_contact:flags.12?true deleted:flags.13?true bot:flags.14?true bot_chat_history:flags.15?true bot_nochats:flags.16?true verified:flags.17?true restricted:flags.18?true min:flags.20?true bot_inline_geo:flags.21?true id:int access_hash:flags.0?long first_name:flags.1?string last_name:flags.2?string username:flags.3?string phone:flags.4?string photo:flags.5?UserProfilePhoto status:flags.6?UserStatus bot_info_version:flags.14?int restriction_reason:flags.18?string bot_inline_placeholder:flags.19?string = User;
const crc_user = 0xd10d979a
type TL_user struct {
	flags int32
	self bool// self:flags.10?true
	contact bool // contact:flags.11?true
	mutual_contact bool // mutual_contact:flags.12?true
	deleted bool // deleted:flags.13?true
	bot bool // bot:flags.14?true
	bot_chat_history bool // bot_chat_history:flags.15?true
	bot_nochats bool // bot_nochats:flags.16?true
	verified bool // verified:flags.17?true
	restricted bool // restricted:flags.18?true
	min bool // min:flags.20?true
	bot_inline_geo bool // bot_inline_geo:flags.21?true
	id int32 // id:int
	access_hash int64 // access_hash:flags.0?long
	first_name string // first_name:flags.1?string
	last_name string // last_name:flags.2?string
	username string // username:flags.3?string
	phone string // phone:flags.4?string
	photo TL // photo:flags.5?UserProfilePhoto
	status TL // status:flags.6?UserStatus
	bot_info_version int32 // bot_info_version:flags.14?int
	restriction_reason string // restriction_reason:flags.18?string
	bot_inline_placeholder string // bot_inline_placeholder:flags.19?string
}

func (e TL_user) encode() []byte  {
	var flags int32
	// fill bits in flags
	if e.self {
		flags |= (1 << 10)
	}
	if e.contact {
		flags |= (1 << 11)
	}
	if e.mutual_contact {
		flags |= (1 << 12)
	}
	if e.deleted {
		flags |= (1 << 13)
	}
	if e.bot {
		flags |= (1 << 14)
	}
	if e.bot_chat_history {
		flags |= (1 << 15)
	}
	if e.bot_nochats {
		flags |= (1 << 16)
	}
	if e.verified {
		flags |= (1 << 17)
	}
	if e.restricted {
		flags |= (1 << 18)
	}
	if e.min {
		flags |= (1 << 20)
	}
	if e.bot_inline_geo {
		flags |= (1 << 21)
	}
	if e.access_hash > 0 {
		flags |= (1 << 0)
	}
	if e.first_name != "" {
		flags |= (1 << 1)
	}
	if e.last_name != "" {
		flags |= (1 << 2)
	}
	if e.username != "" {
		flags |= (1 << 3)
	}
	if e.phone != "" {
		flags |= (1 << 4)
	}
	if _, ok := e.photo.(TL_userProfilePhoto); ok {
		flags |= (1 << 5)
	}
	if _, ok := e.status.(TL_null); !ok {
		flags |= (1 << 6)
	}
	if e.bot_info_version > 0 {
		flags |= (1 << 14)
	}
	if e.restriction_reason != "" {
		flags |= (1 << 18)
	}
	if e.bot_inline_placeholder != "" {
		flags |= (1 << 19)
	}
	x := NewEncodeBuf(512)
	x.UInt(crc_user)
	x.Int(flags)
	x.Int(e.id)
	if flags & (1 << 0) != 0 {
		x.Long(e.access_hash)
	}
	if flags & (1 << 1) != 0 {
		x.String(e.first_name)
	}
	if flags & (1 << 2) != 0 {
		x.String(e.last_name)
	}
	if flags & (1 << 3) != 0 {
		x.String(e.username)
	}
	if flags & (1 << 4) != 0 {
		x.String(e.phone)
	}
	if flags & (1 << 5) != 0 {
		x.Bytes(e.photo.encode())
	}
	if flags & (1 << 6) != 0 {
		x.Bytes(e.status.encode())
	}
	if flags & (1 << 14) != 0 {
		x.Int(e.bot_info_version)
	}
	if flags & (1 << 18) != 0 {
		x.String(e.restriction_reason)
	}
	if flags & (1 << 19) != 0 {
		x.String(e.bot_inline_placeholder)
	}

	return x.buf
}