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
	Code int32
	Text string
}

func (e TL_error) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_error)
	x.Int(e.Code)
	x.String(e.Text)
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
		if flags&(1<<1) != 0 {
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
		if flags&(1<<0) != 0 {
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
			Flags:                    flags,
			Phonecalls_enabled:       phonecalls_enabled,
			Date:                     date,
			Expires:                  expires,
			Test_mode:                test_mode,
			This_dc:                  this_dc,
			Dc_options:               dc_options,
			Chat_size_max:            chat_size_max,
			Megagroup_size_max:       megagroup_size_max,
			Forwarded_count_max:      forwarded_count_max,
			Online_update_period_ms:  online_update_period_ms,
			Offline_blur_timeout_ms:  offline_blur_timeout_ms,
			Offline_idle_timeout_ms:  offline_idle_timeout_ms,
			Online_cloud_timeout_ms:  online_cloud_timeout_ms,
			Notify_cloud_delay_ms:    notify_cloud_delay_ms,
			Notify_default_delay_ms:  notify_default_delay_ms,
			Chat_big_size:            chat_big_size,
			Push_chat_period_ms:      push_chat_period_ms,
			Push_chat_limit:          push_chat_limit,
			Saved_gifs_limit:         saved_gifs_limit,
			Edit_time_limit:          edit_time_limit,
			Rating_e_decay:           rating_e_decay,
			Stickers_recent_limit:    stickers_recent_limit,
			Tmp_sessions:             tmp_sessions,
			Pinned_dialogs_count_max: pinned_dialogs_count_max,
			Call_receive_timeout_ms:  call_receive_timeout_ms,
			Call_ring_timeout_ms:     call_ring_timeout_ms,
			Call_connect_timeout_ms:  call_connect_timeout_ms,
			Call_packet_timeout_ms:   call_packet_timeout_ms,
			Me_url_prefix:            me_url_prefix,
			Disabled_features:        disabled_features,
		}
	case crc_dcOption:
		flags := m.Int()
		var ipv6, media_only, tcpo_only bool
		if flags&(1<<0) != 0 {
			ipv6 = true
		}
		if flags&(1<<1) != 0 {
			media_only = true
		}
		if flags&(1<<2) != 0 {
			tcpo_only = true
		}
		id := m.Int()
		ip_address := m.String()
		port := m.Int()
		r = TL_dcOption{
			Flags:      flags,
			Ipv6:       ipv6,
			Media_only: media_only,
			Tcpo_only:  tcpo_only,
			Id:         id,
			Ip_address: ip_address,
			Port:       port,
		}
	case crc_auth_codeTypeSms:
		r = TL_auth_codeTypeSms{}
	case crc_auth_codeTypeCall:
		r = TL_auth_codeTypeCall{}
	case crc_auth_codeTypeFlashCall:
		r = TL_auth_codeTypeFlashCall{}
	case crc_auth_sentCodeTypeApp:
		r = TL_auth_sentCodeTypeApp{
			Length: m.Int(),
		}
	case crc_auth_sentCodeTypeSms:
		r = TL_auth_sentCodeTypeSms{
			Length: m.Int(),
		}
	case crc_auth_sentCodeTypeCall:
		r = TL_auth_sentCodeTypeCall{
			Length: m.Int(),
		}
	case crc_auth_sentCodeTypeFlashCall:
		r = TL_auth_sentCodeTypeFlashCall{
			Pattern: m.String(),
		}
	case crc_auth_sentCode:
		flags := m.Int()
		phone_registered := false
		if flags&(1<<0) != 0 {
			phone_registered = true
		}
		code_type := m.Object()
		phone_code_hash := m.String()
		var next_type TL
		next_type = TL_null{}
		if flags&(1<<1) != 0 {
			next_type = m.Object()
		}
		var timeout int32
		if flags&(1<<2) != 0 {
			timeout = m.Int()
		}
		r = TL_auth_sentCode{
			Flags:            flags,
			Phone_registered: phone_registered,
			Code_type:        code_type,
			Phone_code_hash:  phone_code_hash,
			Next_type:        next_type,
			Timeout:          timeout,
		}
	case crc_auth_sendCode:
		flags := m.Int()
		allow_flashcall := false
		if flags&(1<<0) != 0 {
			allow_flashcall = true
		}
		phone_number := m.String()
		var current_number TL
		if flags&(1<<0) != 0 {
			current_number = m.Object()
		}
		api_id := m.Int()
		api_hash := m.String()
		r = TL_auth_sendCode{
			Flags:           flags,
			Allow_flashcall: allow_flashcall,
			Phone_number:    phone_number,
			Current_number:  current_number,
			Api_id:          api_id,
			Api_hash:        api_hash,
		}
	case crc_auth_signIn:
		r = TL_auth_signIn{
			Phone_number:    m.String(),
			Phone_code_hash: m.String(),
			Phone_code:      m.String(),
		}
	case crc_auth_authorization:
		flags := m.Int()
		var tmp_sessions int32
		if flags&(1<<0) != 0 {
			tmp_sessions = m.Int()
		}
		user := m.Object()
		r = TL_auth_authorization{
			Flags:        flags,
			Tmp_sessions: tmp_sessions,
			User:         user,
		}
	case crc_userEmpty:
		r = TL_userEmpty{
			Id: m.Int(),
		}
	case crc_user:
		flags := m.Int()
		self := false
		if flags&(1<<10) != 0 {
			self = true
		}
		contact := false
		if flags&(1<<11) != 0 {
			contact = true
		}
		mutual_contact := false
		if flags&(1<<12) != 0 {
			mutual_contact = true
		}
		deleted := false
		if flags&(1<<13) != 0 {
			deleted = true
		}
		bot := false
		if flags&(1<<14) != 0 {
			bot = true
		}
		bot_chat_history := false
		if flags&(1<<15) != 0 {
			bot_chat_history = true
		}
		bot_nochats := false
		if flags&(1<<16) != 0 {
			bot_nochats = true
		}
		verified := false
		if flags&(1<<17) != 0 {
			verified = true
		}
		restricted := false
		if flags&(1<<18) != 0 {
			restricted = true
		}
		min := false
		if flags&(1<<20) != 0 {
			min = true
		}
		bot_inline_geo := false
		if flags&(1<<21) != 0 {
			bot_inline_geo = true
		}
		id := m.Int()
		var access_hash int64
		if flags&(1<<0) != 0 {
			access_hash = m.Long()
		}
		var first_name, last_name, username, phone string
		if flags&(1<<1) != 0 {
			first_name = m.String()
		}
		if flags&(1<<2) != 0 {
			last_name = m.String()
		}
		if flags&(1<<3) != 0 {
			username = m.String()
		}
		if flags&(1<<4) != 0 {
			phone = m.String()
		}
		var photo, status TL
		if flags&(1<<5) != 0 {
			photo = m.Object()
		}
		if flags&(1<<6) != 0 {
			status = m.Object()
		}
		var bot_info_version int32
		if flags&(1<<14) != 0 {
			bot_info_version = m.Int()
		}
		var restriction_reason, bot_inline_placeholder string
		if flags&(1<<18) != 0 {
			restriction_reason = m.String()
		}
		if flags&(1<<19) != 0 {
			bot_inline_placeholder = m.String()
		}
		r = TL_user{
			Flags:                  flags,
			Self:                   self,
			Contact:                contact,
			Mutual_contact:         mutual_contact,
			Deleted:                deleted,
			Bot:                    bot,
			Bot_chat_history:       bot_chat_history,
			Bot_nochats:            bot_nochats,
			Verified:               verified,
			Restricted:             restricted,
			Min:                    min,
			Bot_inline_geo:         bot_inline_geo,
			Id:                     id,
			Access_hash:            access_hash,
			First_name:             first_name,
			Last_name:              last_name,
			Username:               username,
			Phone:                  phone,
			Photo:                  photo,
			Status:                 status,
			Bot_info_version:       bot_info_version,
			Restriction_reason:     restriction_reason,
			Bot_inline_placeholder: bot_inline_placeholder,
		}
	case crc_userProfilePhotoEmpty:
		r = TL_userProfilePhotoEmpty{}
	case crc_userProfilePhoto:
		r = TL_userProfilePhoto{
			Photo_id:    m.Long(),
			Photo_small: m.Object(),
			Photo_big:   m.Object(),
		}
	case crc_fileLocationUnavailable:
		r = TL_fileLocationUnavailable{
			Volume_id: m.Long(),
			Local_id:  m.Int(),
			Secret:    m.Long(),
		}
	case crc_fileLocation:
		r = TL_fileLocation{
			Dc_id:     m.Int(),
			Volume_id: m.Long(),
			Local_id:  m.Int(),
			Secret:    m.Long(),
		}
	case crc_userStatusEmpty:
		r = TL_userStatusEmpty{}
	case crc_userStatusOnline:
		r = TL_userStatusOnline{
			Expires: m.Int(),
		}
	case crc_userStatusOffline:
		r = TL_userStatusOffline{
			Was_online: m.Int(),
		}
	case crc_userStatusRecently:
		r = TL_userStatusRecently{}
	case crc_userStatusLastWeek:
		r = TL_userStatusLastWeek{}
	case crc_userStatusLastMonth:
		r = TL_userStatusLastMonth{}
	case crc_contacts_topPeersNotModified:
		r = TL_contacts_topPeersNotModified{}
	case crc_contacts_topPeers:
		r = TL_contacts_topPeers{
			Categories: m.Vector(),
			Chats:      m.Vector(),
			Users:      m.Vector(),
		}
	case crc_topPeerCategoryPeers:
		r = TL_topPeerCategoryPeers{
			Category: m.Object(),
			Count:    m.Int(),
			Peers:    m.Vector(),
		}
	case crc_topPeer:
		r = TL_topPeer{
			Peer:   m.Object(),
			Rating: m.Double(),
		}
	case crc_peerUser:
		r = TL_peerUser{
			User_id: m.Int(),
		}
	case crc_peerChat:
		r = TL_peerChat{
			Chat_id: m.Int(),
		}
	case crc_peerChannel:
		r = TL_peerChannel{
			Channel_id: m.Int(),
		}
	case crc_topPeerCategoryBotsPM:
		r = TL_topPeerCategoryBotsPM{}
	case crc_topPeerCategoryBotsInline:
		r = TL_topPeerCategoryBotsInline{}
	case crc_topPeerCategoryCorrespondents:
		r = TL_topPeerCategoryCorrespondents{}
	case crc_topPeerCategoryGroups:
		r = TL_topPeerCategoryGroups{}
	case crc_topPeerCategoryChannels:
		r = TL_topPeerCategoryChannels{}
	case crc_chatEmpty:
		r = TL_chatEmpty{
			Id: m.Int(),
		}
	case crc_chat:
		flags := m.Int()
		creator := flags&(1<<0) != 0
		kicked := flags&(1<<1) != 0
		left := flags&(1<<2) != 0
		admins_enabled := flags&(1<<3) != 0
		admin := flags&(1<<4) != 0
		deactivated := flags&(1<<5) != 0
		id := m.Int()
		title := m.String()
		photo := m.Object()
		participants_count := m.Int()
		date := m.Int()
		version := m.Int()
		var migrated_to TL
		migrated_to = TL_null{}
		if flags&(1<<6) != 0 {
			migrated_to = m.Object()
		}
		r = TL_chat{
			Flags:              flags,
			Creator:            creator,
			Kicked:             kicked,
			Left:               left,
			Admins_enabled:     admins_enabled,
			Admin:              admin,
			Deactivated:        deactivated,
			Id:                 id,
			Title:              title,
			Photo:              photo,
			Participants_count: participants_count,
			Date:               date,
			Version:            version,
			Migrated_to:        migrated_to,
		}
	case crc_chatForbidden:
		r = TL_chatForbidden{
			Id:    m.Int(),
			Title: m.String(),
		}
	case crc_chatPhotoEmpty:
		r = TL_chatPhotoEmpty{}
	case crc_chatPhoto:
		r = TL_chatPhoto{
			Photo_small: m.Object(),
			Photo_big:   m.Object(),
		}
	case crc_channel:
		flags := m.Int()
		creator := flags&(1<<0) != 0
		kicked := flags&(1<<1) != 0
		left := flags&(1<<2) != 0
		editor := flags&(1<<3) != 0
		moderator := flags&(1<<4) != 0
		broadcast := flags&(1<<5) != 0
		verified := flags&(1<<7) != 0
		megagroup := flags&(1<<8) != 0
		restricted := flags&(1<<9) != 0
		democracy := flags&(1<<10) != 0
		signatures := flags&(1<<11) != 0
		min := flags&(1<<12) != 0
		id := m.Int()
		var access_hash int64
		if flags&(1<<13) != 0 {
			access_hash = m.Long()
		}
		title := m.String()
		var username string
		if flags&(1<<6) != 0 {
			username = m.String()
		}
		photo := m.Object()
		date := m.Int()
		version := m.Int()
		var restriction_reason string
		if flags&(1<<9) != 0 {
			restriction_reason = m.String()
		}
		r = TL_channel{
			Flags:              flags,
			Creator:            creator,
			Kicked:             kicked,
			Left:               left,
			Editor:             editor,
			Moderator:          moderator,
			Broadcast:          broadcast,
			Verified:           verified,
			Megagroup:          megagroup,
			Restricted:         restricted,
			Democracy:          democracy,
			Signatures:         signatures,
			Min:                min,
			Id:                 id,
			Access_hash:        access_hash,
			Title:              title,
			Username:           username,
			Photo:              photo,
			Date:               date,
			Version:            version,
			Restriction_reason: restriction_reason,
		}
	case crc_channelForbidden:
		flags := m.Int()
		broadcast := flags&(1<<5) != 0
		megagroup := flags&(1<<8) != 0
		id := m.Int()
		access_hash := m.Long()
		title := m.String()
		r = TL_channelForbidden{
			Flags:       flags,
			Broadcast:   broadcast,
			Megagroup:   megagroup,
			Id:          id,
			Access_hash: access_hash,
			Title:       title,
		}
	case crc_inputChannelEmpty:
		r = TL_inputChannelEmpty{}
	case crc_inputChannel:
		r = TL_inputChannel{
			Channel_id:  m.Int(),
			Access_hash: m.Long(),
		}
	case crc_messageEmpty:
		// TODO
	case crc_messages_messages:
		// TODO
	case crc_messages_messagesSlice:
		// TODO
	case crc_messages_channelMessages:
		// TODO
	case crc_messages_getHistory:
		// TODO
	default:
		m.err = fmt.Errorf("Unknown constructor: %x", constructor)
		return nil
	}
	return
}

// invokeWithLayer#da9b0d0d {X:Type} Layer:int Query:!X = X;
const crc_invokeWithLayer = 0xda9b0d0d

type TL_invokeWithLayer struct {
	Layer int32
	Query TL
}

func (e TL_invokeWithLayer) encode() []byte {
	// TODO: 512 is a magic number
	x := NewEncodeBuf(512)
	x.UInt(crc_invokeWithLayer)
	x.Int(e.Layer)
	x.Bytes(e.Query.encode())
	// TODO: Should I shrink a buffer to his actual Size or not?
	return x.buf
}

// initConnection#69796de9 {X:Type} Api_id:int Device_model:string System_version:string App_version:string Lang_code:string Query:!X = X;
const crc_initConnection = 0x69796de9

type TL_initConnection struct {
	Api_id         int32
	Device_model   string
	System_version string
	App_version    string
	Lang_code      string
	Query          TL
}

func (e TL_initConnection) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_initConnection)
	x.Int(e.Api_id)
	x.String(e.Device_model)
	x.String(e.System_version)
	x.String(e.App_version)
	x.String(e.Lang_code)
	x.Bytes(e.Query.encode())
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

// config#cb601684 Flags:# Phonecalls_enabled:Flags.1?true Date:int Expires:int Test_mode:Bool This_dc:int Dc_options:Vector<DcOption> Chat_size_max:int Megagroup_size_max:int Forwarded_count_max:int Online_update_period_ms:int Offline_blur_timeout_ms:int Offline_idle_timeout_ms:int Online_cloud_timeout_ms:int Notify_cloud_delay_ms:int Notify_default_delay_ms:int Chat_big_size:int Push_chat_period_ms:int Push_chat_limit:int Saved_gifs_limit:int Edit_time_limit:int Rating_e_decay:int Stickers_recent_limit:int Tmp_sessions:Flags.0?int Pinned_dialogs_count_max:int Call_receive_timeout_ms:int Call_ring_timeout_ms:int Call_connect_timeout_ms:int Call_packet_timeout_ms:int Me_url_prefix:string Disabled_features:Vector<DisabledFeature> = Config;
const crc_config = 0xcb601684

type TL_config struct {
	Flags                    int32
	Phonecalls_enabled       bool // Flags.1?true TODO: TL_true
	Date                     int32
	Expires                  int32
	Test_mode                TL // TL_boolFalse or TL_boolTrue
	This_dc                  int32
	Dc_options               []TL // DcOption
	Chat_size_max            int32
	Megagroup_size_max       int32
	Forwarded_count_max      int32
	Online_update_period_ms  int32
	Offline_blur_timeout_ms  int32
	Offline_idle_timeout_ms  int32
	Online_cloud_timeout_ms  int32
	Notify_cloud_delay_ms    int32
	Notify_default_delay_ms  int32
	Chat_big_size            int32
	Push_chat_period_ms      int32
	Push_chat_limit          int32
	Saved_gifs_limit         int32
	Edit_time_limit          int32
	Rating_e_decay           int32
	Stickers_recent_limit    int32
	Tmp_sessions             int32 // Flags.0?int
	Pinned_dialogs_count_max int32
	Call_receive_timeout_ms  int32
	Call_ring_timeout_ms     int32
	Call_connect_timeout_ms  int32
	Call_packet_timeout_ms   int32
	Me_url_prefix            string
	Disabled_features        []TL // DisabledFeature
}

func (e TL_config) encode() []byte { return nil }

// dcOption#5d8c6cc Flags:# Ipv6:Flags.0?true Media_only:Flags.1?true Tcpo_only:Flags.2?true Id:int Ip_address:string Port:int = DcOption;
const crc_dcOption = 0x5d8c6cc

type TL_dcOption struct {
	Flags      int32
	Ipv6       bool // Ipv6:Flags.0?true TODO: TL_true
	Media_only bool // Media_only:Flags.1?true TODO: TL_true
	Tcpo_only  bool // Tcpo_only:Flags.2?true TODO: TL_true
	Id         int32
	Ip_address string
	Port       int32
}

func (e TL_dcOption) encode() []byte { return nil }

//auth.codeTypeSms#72a3158c = auth.CodeType;
const crc_auth_codeTypeSms = 0x72a3158c

type TL_auth_codeTypeSms struct{}

func (e TL_auth_codeTypeSms) encode() []byte {
	x := NewEncodeBuf(4)
	x.UInt(crc_auth_codeTypeSms)
	return x.buf
}

//auth.codeTypeCall#741cd3e3 = auth.CodeType;
const crc_auth_codeTypeCall = 0x741cd3e3

type TL_auth_codeTypeCall struct{}

func (e TL_auth_codeTypeCall) encode() []byte {
	x := NewEncodeBuf(4)
	x.UInt(crc_auth_codeTypeCall)
	return x.buf
}

//auth.codeTypeFlashCall#226ccefb = auth.CodeType;
const crc_auth_codeTypeFlashCall = 0x226ccefb

type TL_auth_codeTypeFlashCall struct{}

func (e TL_auth_codeTypeFlashCall) encode() []byte {
	x := NewEncodeBuf(4)
	x.UInt(crc_auth_codeTypeFlashCall)
	return x.buf
}

//auth.sentCodeTypeApp#3dbb5986 Length:int = auth.SentCodeType;
const crc_auth_sentCodeTypeApp = 0x3dbb5986

type TL_auth_sentCodeTypeApp struct {
	Length int32
}

func (e TL_auth_sentCodeTypeApp) encode() []byte {
	x := NewEncodeBuf(8)
	x.UInt(crc_auth_sentCodeTypeApp)
	x.Int(e.Length)
	return x.buf
}

//auth.sentCodeTypeSms#c000bba2 Length:int = auth.SentCodeType;
const crc_auth_sentCodeTypeSms = 0xc000bba2

type TL_auth_sentCodeTypeSms struct {
	Length int32
}

func (e TL_auth_sentCodeTypeSms) encode() []byte {
	x := NewEncodeBuf(8)
	x.UInt(crc_auth_sentCodeTypeSms)
	x.Int(e.Length)
	return x.buf
}

//auth.sentCodeTypeCall#5353e5a7 Length:int = auth.SentCodeType;
const crc_auth_sentCodeTypeCall = 0x5353e5a7

type TL_auth_sentCodeTypeCall struct {
	Length int32
}

func (e TL_auth_sentCodeTypeCall) encode() []byte {
	x := NewEncodeBuf(8)
	x.UInt(crc_auth_sentCodeTypeCall)
	x.Int(e.Length)
	return x.buf
}

//auth.sentCodeTypeFlashCall#ab03c6d9 Pattern:string = auth.SentCodeType;
const crc_auth_sentCodeTypeFlashCall = 0xab03c6d9

type TL_auth_sentCodeTypeFlashCall struct {
	Pattern string
}

func (e TL_auth_sentCodeTypeFlashCall) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_auth_sentCodeTypeFlashCall)
	x.String(e.Pattern)
	return x.buf
}

// auth.sentCode#5e002502 Flags:# Phone_registered:Flags.0?true type:auth.SentCodeType Phone_code_hash:string Next_type:Flags.1?auth.CodeType timeout:Flags.2?int = auth.SentCode;
const crc_auth_sentCode = 0x5e002502

type TL_auth_sentCode struct {
	Flags            int32
	Phone_registered bool
	Code_type        TL // type:auth.SentCodeType
	Phone_code_hash  string
	Next_type        TL
	Timeout          int32
}

func (e TL_auth_sentCode) encode() []byte {
	var flags int32
	x := NewEncodeBuf(512)
	x.UInt(crc_auth_sentCode)
	// fill bits in Flags
	if e.Phone_registered {
		flags |= (1 << 0)
	}
	if _, ok := (e.Next_type).(TL_null); !ok {
		flags |= (1 << 1)
	}
	if e.Timeout > 0 {
		flags |= (1 << 2)
	}
	x.Int(flags)
	x.Bytes(e.Code_type.encode())
	x.String(e.Phone_code_hash)
	if _, ok := (e.Next_type).(TL_null); !ok {
		x.Bytes(e.Next_type.encode())
	}
	if e.Timeout > 0 {
		x.Int(e.Timeout)
	}
	return x.buf
}

// auth.sendCode#86aef0ec Flags:# Allow_flashcall:Flags.0?true Phone_number:string Current_number:Flags.0?Bool Api_id:int Api_hash:string = auth.SentCode;
const crc_auth_sendCode = 0x86aef0ec

type TL_auth_sendCode struct {
	Flags           int32
	Allow_flashcall bool // Allow_flashcall:Flags.0?true
	Phone_number    string
	Current_number  TL // Current_number:Flags.0?Bool
	Api_id          int32
	Api_hash        string
}

func (e TL_auth_sendCode) encode() []byte {
	var flags int32
	x := NewEncodeBuf(512)
	x.UInt(crc_auth_sendCode)
	if e.Allow_flashcall {
		flags |= (1 << 0)
	}
	x.Int(flags)
	x.String(e.Phone_number)
	if e.Allow_flashcall {
		x.Bytes(e.Current_number.encode())
	}
	x.Int(e.Api_id)
	x.String(e.Api_hash)
	return x.buf
}

// auth.signIn#bcd51581 Phone_number:string Phone_code_hash:string Phone_code:string = auth.Authorization;
const crc_auth_signIn = 0xbcd51581

type TL_auth_signIn struct {
	Phone_number    string
	Phone_code_hash string
	Phone_code      string
}

func (e TL_auth_signIn) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_auth_signIn)
	x.String(e.Phone_number)
	x.String(e.Phone_code_hash)
	x.String(e.Phone_code)
	return x.buf
}

// auth.authorization#cd050916 Flags:# Tmp_sessions:Flags.0?int User:User = auth.Authorization;
const crc_auth_authorization = 0xcd050916

type TL_auth_authorization struct {
	Flags        int32
	Tmp_sessions int32
	User         TL
}

func (e TL_auth_authorization) encode() []byte {
	var flags int32
	x := NewEncodeBuf(512)
	x.UInt(crc_auth_authorization)
	// TODO: I am not sure about this condition. Check how serialization works in other libraries
	if e.Tmp_sessions > 0 {
		flags |= (1 << 0)
	}
	x.Int(flags)
	if e.Tmp_sessions > 0 {
		x.Int(e.Tmp_sessions)
	}
	x.Bytes(e.User.encode())
	return x.buf
}

//fileLocationUnavailable#7c596b46 Volume_id:long Local_id:int Secret:long = FileLocation;
const crc_fileLocationUnavailable = 0x7c596b46

type TL_fileLocationUnavailable struct {
	Volume_id int64
	Local_id  int32
	Secret    int64
}

func (e TL_fileLocationUnavailable) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_fileLocationUnavailable)
	x.Long(e.Volume_id)
	x.Int(e.Local_id)
	x.Long(e.Secret)
	return x.buf
}

//fileLocation#53d69076 Dc_id:int Volume_id:long Local_id:int Secret:long = FileLocation;
const crc_fileLocation = 0x53d69076

type TL_fileLocation struct {
	Dc_id     int32
	Volume_id int64
	Local_id  int32
	Secret    int64
}

func (e TL_fileLocation) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_fileLocation)
	x.Int(e.Dc_id)
	x.Long(e.Volume_id)
	x.Int(e.Local_id)
	x.Long(e.Secret)
	return x.buf
}

//userProfilePhotoEmpty#4f11bae1 = UserProfilePhoto;
const crc_userProfilePhotoEmpty = 0x4f11bae1

type TL_userProfilePhotoEmpty struct{}

func (e TL_userProfilePhotoEmpty) encode() []byte {
	x := NewEncodeBuf(4)
	x.UInt(crc_userProfilePhotoEmpty)
	return x.buf
}

//userProfilePhoto#d559d8c8 Photo_id:long Photo_small:FileLocation Photo_big:FileLocation = UserProfilePhoto;
const crc_userProfilePhoto = 0xd559d8c8

type TL_userProfilePhoto struct {
	Photo_id    int64
	Photo_small TL // FileLocation
	Photo_big   TL // FileLocation
}

func (e TL_userProfilePhoto) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_userProfilePhoto)
	x.Long(e.Photo_id)
	x.Bytes(e.Photo_small.encode())
	x.Bytes(e.Photo_big.encode())
	return x.buf
}

// userStatusEmpty#9d05049 = UserStatus;
const crc_userStatusEmpty = 0x9d05049

type TL_userStatusEmpty struct{}

func (e TL_userStatusEmpty) encode() []byte {
	x := NewEncodeBuf(4)
	x.UInt(crc_userStatusEmpty)
	return x.buf
}

// userStatusOnline#edb93949 Expires:int = UserStatus;
const crc_userStatusOnline = 0xedb93949

type TL_userStatusOnline struct {
	Expires int32
}

func (e TL_userStatusOnline) encode() []byte {
	x := NewEncodeBuf(8)
	x.UInt(crc_userStatusOnline)
	x.Int(e.Expires)
	return x.buf
}

// userStatusOffline#8c703f Was_online:int = UserStatus;
const crc_userStatusOffline = 0x8c703f

type TL_userStatusOffline struct {
	Was_online int32
}

func (e TL_userStatusOffline) encode() []byte {
	x := NewEncodeBuf(8)
	x.UInt(crc_userStatusOffline)
	x.Int(e.Was_online)
	return x.buf
}

// userStatusRecently#e26f42f1 = UserStatus;
const crc_userStatusRecently = 0xe26f42f1

type TL_userStatusRecently struct{}

func (e TL_userStatusRecently) encode() []byte {
	x := NewEncodeBuf(4)
	x.UInt(crc_userStatusRecently)
	return x.buf
}

// userStatusLastWeek#7bf09fc = UserStatus;
const crc_userStatusLastWeek = 0x7bf09fc

type TL_userStatusLastWeek struct{}

func (e TL_userStatusLastWeek) encode() []byte {
	x := NewEncodeBuf(4)
	x.UInt(crc_userStatusLastWeek)
	return x.buf
}

// userStatusLastMonth#77ebc742 = UserStatus;
const crc_userStatusLastMonth = 0x77ebc742

type TL_userStatusLastMonth struct{}

func (e TL_userStatusLastMonth) encode() []byte {
	x := NewEncodeBuf(4)
	x.UInt(crc_userStatusLastMonth)
	return x.buf
}

//userEmpty#200250ba Id:int = User;
const crc_userEmpty = 0x200250ba

type TL_userEmpty struct {
	Id int32
}

func (e TL_userEmpty) encode() []byte {
	x := NewEncodeBuf(8)
	x.UInt(crc_userEmpty)
	x.Int(e.Id)

	return x.buf
}

//User#d10d979a Flags:# Self:Flags.10?true Contact:Flags.11?true Mutual_contact:Flags.12?true Deleted:Flags.13?true Bot:Flags.14?true Bot_chat_history:Flags.15?true Bot_nochats:Flags.16?true Verified:Flags.17?true Restricted:Flags.18?true Min:Flags.20?true Bot_inline_geo:Flags.21?true Id:int Access_hash:Flags.0?long First_name:Flags.1?string Last_name:Flags.2?string Username:Flags.3?string Phone:Flags.4?string Photo:Flags.5?UserProfilePhoto Status:Flags.6?UserStatus Bot_info_version:Flags.14?int Restriction_reason:Flags.18?string Bot_inline_placeholder:Flags.19?string = User;
const crc_user = 0xd10d979a

type TL_user struct {
	Flags                  int32
	Self                   bool   // Self:Flags.10?true
	Contact                bool   // Contact:Flags.11?true
	Mutual_contact         bool   // Mutual_contact:Flags.12?true
	Deleted                bool   // Deleted:Flags.13?true
	Bot                    bool   // Bot:Flags.14?true
	Bot_chat_history       bool   // Bot_chat_history:Flags.15?true
	Bot_nochats            bool   // Bot_nochats:Flags.16?true
	Verified               bool   // Verified:Flags.17?true
	Restricted             bool   // Restricted:Flags.18?true
	Min                    bool   // Min:Flags.20?true
	Bot_inline_geo         bool   // Bot_inline_geo:Flags.21?true
	Id                     int32  // Id:int
	Access_hash            int64  // Access_hash:Flags.0?long
	First_name             string // First_name:Flags.1?string
	Last_name              string // Last_name:Flags.2?string
	Username               string // Username:Flags.3?string
	Phone                  string // Phone:Flags.4?string
	Photo                  TL     // Photo:Flags.5?UserProfilePhoto
	Status                 TL     // Status:Flags.6?UserStatus
	Bot_info_version       int32  // Bot_info_version:Flags.14?int
	Restriction_reason     string // Restriction_reason:Flags.18?string
	Bot_inline_placeholder string // Bot_inline_placeholder:Flags.19?string
}

func (e TL_user) encode() []byte {
	var flags int32
	// fill bits in Flags
	if e.Self {
		flags |= (1 << 10)
	}
	if e.Contact {
		flags |= (1 << 11)
	}
	if e.Mutual_contact {
		flags |= (1 << 12)
	}
	if e.Deleted {
		flags |= (1 << 13)
	}
	if e.Bot {
		flags |= (1 << 14)
	}
	if e.Bot_chat_history {
		flags |= (1 << 15)
	}
	if e.Bot_nochats {
		flags |= (1 << 16)
	}
	if e.Verified {
		flags |= (1 << 17)
	}
	if e.Restricted {
		flags |= (1 << 18)
	}
	if e.Min {
		flags |= (1 << 20)
	}
	if e.Bot_inline_geo {
		flags |= (1 << 21)
	}
	if e.Access_hash > 0 {
		flags |= (1 << 0)
	}
	if e.First_name != "" {
		flags |= (1 << 1)
	}
	if e.Last_name != "" {
		flags |= (1 << 2)
	}
	if e.Username != "" {
		flags |= (1 << 3)
	}
	if e.Phone != "" {
		flags |= (1 << 4)
	}
	if _, ok := e.Photo.(TL_userProfilePhoto); ok {
		flags |= (1 << 5)
	}
	if _, ok := e.Status.(TL_null); !ok {
		flags |= (1 << 6)
	}
	if e.Bot_info_version > 0 {
		flags |= (1 << 14)
	}
	if e.Restriction_reason != "" {
		flags |= (1 << 18)
	}
	if e.Bot_inline_placeholder != "" {
		flags |= (1 << 19)
	}
	x := NewEncodeBuf(512)
	x.UInt(crc_user)
	x.Int(flags)
	x.Int(e.Id)
	if flags&(1<<0) != 0 {
		x.Long(e.Access_hash)
	}
	if flags&(1<<1) != 0 {
		x.String(e.First_name)
	}
	if flags&(1<<2) != 0 {
		x.String(e.Last_name)
	}
	if flags&(1<<3) != 0 {
		x.String(e.Username)
	}
	if flags&(1<<4) != 0 {
		x.String(e.Phone)
	}
	if flags&(1<<5) != 0 {
		x.Bytes(e.Photo.encode())
	}
	if flags&(1<<6) != 0 {
		x.Bytes(e.Status.encode())
	}
	if flags&(1<<14) != 0 {
		x.Int(e.Bot_info_version)
	}
	if flags&(1<<18) != 0 {
		x.String(e.Restriction_reason)
	}
	if flags&(1<<19) != 0 {
		x.String(e.Bot_inline_placeholder)
	}

	return x.buf
}

//peerUser#9db1bc6d user_id:int = Peer;
const crc_peerUser = 0x9db1bc6d

type TL_peerUser struct {
	User_id int32
}

func (e TL_peerUser) encode() []byte {
	x := NewEncodeBuf(8)
	x.UInt(crc_peerUser)
	x.Int(e.User_id)
	return x.buf
}

//peerChat#bad0e5bb chat_id:int = Peer;
const crc_peerChat = 0xbad0e5bb

type TL_peerChat struct {
	Chat_id int32
}

func (e TL_peerChat) encode() []byte {
	x := NewEncodeBuf(8)
	x.UInt(crc_peerChat)
	x.Int(e.Chat_id)
	return x.buf
}

//peerChannel#bddde532 channel_id:int = Peer;
const crc_peerChannel = 0xbddde532

type TL_peerChannel struct {
	Channel_id int32
}

func (e TL_peerChannel) encode() []byte {
	x := NewEncodeBuf(8)
	x.UInt(crc_peerChannel)
	x.Int(e.Channel_id)
	return x.buf
}

//topPeerCategoryBotsPM#ab661b5b = TopPeerCategory;
const crc_topPeerCategoryBotsPM = 0xab661b5b

type TL_topPeerCategoryBotsPM struct{}

func (e TL_topPeerCategoryBotsPM) encode() []byte { return nil }

//topPeerCategoryBotsInline#148677e2 = TopPeerCategory;
const crc_topPeerCategoryBotsInline = 0x148677e2

type TL_topPeerCategoryBotsInline struct{}

func (e TL_topPeerCategoryBotsInline) encode() []byte { return nil }

//topPeerCategoryCorrespondents#637b7ed = TopPeerCategory;
const crc_topPeerCategoryCorrespondents = 0x637b7ed

type TL_topPeerCategoryCorrespondents struct{}

func (e TL_topPeerCategoryCorrespondents) encode() []byte { return nil }

//topPeerCategoryGroups#bd17a14a = TopPeerCategory;
const crc_topPeerCategoryGroups = 0xbd17a14a

type TL_topPeerCategoryGroups struct{}

func (e TL_topPeerCategoryGroups) encode() []byte { return nil }

//topPeerCategoryChannels#161d9628 = TopPeerCategory;
const crc_topPeerCategoryChannels = 0x161d9628

type TL_topPeerCategoryChannels struct{}

func (e TL_topPeerCategoryChannels) encode() []byte { return nil }

//topPeer#edcdc05b peer:Peer rating:double = TopPeer;
const crc_topPeer = 0xedcdc05b

type TL_topPeer struct {
	Peer   TL
	Rating float64
}

func (e TL_topPeer) encode() []byte { return nil }

//contacts.topPeersNotModified#de266ef5 = contacts.TopPeers;
const crc_contacts_topPeersNotModified = 0xde266ef5

type TL_contacts_topPeersNotModified struct{}

func (e TL_contacts_topPeersNotModified) encode() []byte { return nil }

//contacts.topPeers#70b772a8 categories:Vector<TopPeerCategoryPeers> chats:Vector<Chat> users:Vector<User> = contacts.TopPeers;
const crc_contacts_topPeers = 0x70b772a8

type TL_contacts_topPeers struct {
	Categories []TL
	Chats      []TL
	Users      []TL
}

func (e TL_contacts_topPeers) encode() []byte { return nil }

//topPeerCategoryPeers#fb834291 category:TopPeerCategory count:int peers:Vector<TopPeer> = TopPeerCategoryPeers;
const crc_topPeerCategoryPeers = 0xfb834291

type TL_topPeerCategoryPeers struct {
	Category TL
	Count    int32
	Peers    []TL
}

func (e TL_topPeerCategoryPeers) encode() []byte { return nil }

//contacts.getTopPeers#d4982db5 flags:# correspondents:flags.0?true bots_pm:flags.1?true bots_inline:flags.2?true groups:flags.10?true channels:flags.15?true offset:int limit:int hash:int = contacts.TopPeers;
const crc_contacts_getTopPeers = 0xd4982db5

type TL_contacts_getTopPeers struct {
	Correspondents bool // flags.0?true
	Bots_pm        bool // flags.1?true
	Bots_inline    bool // flags.2?true
	Groups         bool // flags.10?true
	Channels       bool // flags.15?true
	Offset         int32
	Limit          int32
	Hash           int32
}

func (e TL_contacts_getTopPeers) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_contacts_getTopPeers)
	var flags int32
	if e.Correspondents {
		flags |= (1 << 0)
	}
	if e.Bots_pm {
		flags |= (1 << 1)
	}
	if e.Bots_inline {
		flags |= (1 << 2)
	}
	if e.Groups {
		flags |= (1 << 10)
	}
	if e.Channels {
		flags |= (1 << 15)
	}
	x.Int(flags)
	x.Int(e.Offset)
	x.Int(e.Limit)
	x.Int(e.Hash)
	return x.buf
}

//chatPhotoEmpty#37c1011c = ChatPhoto;
const crc_chatPhotoEmpty = 0x37c1011c

type TL_chatPhotoEmpty struct{}

func (e TL_chatPhotoEmpty) encode() []byte {
	x := NewEncodeBuf(4)
	x.UInt(crc_chatPhotoEmpty)
	return x.buf
}

//chatPhoto#6153276a photo_small:FileLocation photo_big:FileLocation = ChatPhoto;
const crc_chatPhoto = 0x6153276a

type TL_chatPhoto struct {
	Photo_small TL
	Photo_big   TL
}

func (e TL_chatPhoto) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_chatPhoto)
	x.Bytes(e.Photo_small.encode())
	x.Bytes(e.Photo_big.encode())
	return x.buf
}

//chatEmpty#9ba2d800 id:int = Chat;
const crc_chatEmpty = 0x9ba2d800

type TL_chatEmpty struct {
	Id int32
}

func (e TL_chatEmpty) encode() []byte {
	x := NewEncodeBuf(8)
	x.UInt(crc_chatEmpty)
	x.Int(e.Id)
	return x.buf
}

//chat#d91cdd54 flags:# creator:flags.0?true kicked:flags.1?true left:flags.2?true admins_enabled:flags.3?true admin:flags.4?true deactivated:flags.5?true id:int title:string photo:ChatPhoto participants_count:int date:int version:int migrated_to:flags.6?InputChannel = Chat;
const crc_chat = 0xd91cdd54

type TL_chat struct {
	Flags              int32
	Creator            bool // flags.0?true
	Kicked             bool // flags.1?true
	Left               bool // flags.2?true
	Admins_enabled     bool // flags.3?true
	Admin              bool // flags.4?true
	Deactivated        bool // flags.5?true
	Id                 int32
	Title              string
	Photo              TL
	Participants_count int32
	Date               int32
	Version            int32
	Migrated_to        TL // flags6?InputChannel
}

func (e TL_chat) encode() []byte { return nil }

//chatForbidden#7328bdb id:int title:string = Chat;
const crc_chatForbidden = 0x7328bdb

type TL_chatForbidden struct {
	Id    int32
	Title string
}

func (e TL_chatForbidden) encode() []byte { return nil }

//channel#a14dca52 flags:# creator:flags.0?true kicked:flags.1?true left:flags.2?true editor:flags.3?true moderator:flags.4?true broadcast:flags.5?true verified:flags.7?true megagroup:flags.8?true restricted:flags.9?true democracy:flags.10?true signatures:flags.11?true min:flags.12?true id:int access_hash:flags.13?long title:string username:flags.6?string photo:ChatPhoto date:int version:int restriction_reason:flags.9?string = Chat;
const crc_channel = 0xa14dca52

type TL_channel struct {
	Flags              int32
	Creator            bool // flags.0?true
	Kicked             bool // flags.1?true
	Left               bool // flags.2?true
	Editor             bool // flags.3?true
	Moderator          bool // flags.4?true
	Broadcast          bool // flags.5?true
	Verified           bool // flags.7?true
	Megagroup          bool // flags.8?true
	Restricted         bool // flags.9?true
	Democracy          bool // flags.10?true
	Signatures         bool // flags.11?true
	Min                bool // flags.12?true
	Id                 int32
	Access_hash        int64 // flags.13?true
	Title              string
	Username           string // flags.6?true
	Photo              TL     // ChatPhoto
	Date               int32
	Version            int32
	Restriction_reason string // // flags.9?true
}

func (e TL_channel) encode() []byte { return nil }

//channelForbidden#8537784f flags:# broadcast:flags.5?true megagroup:flags.8?true id:int access_hash:long title:string = Chat;
const crc_channelForbidden = 0x8537784f

type TL_channelForbidden struct {
	Flags       int32
	Broadcast   bool // flags.5?true
	Megagroup   bool // flags.8?true
	Id          int32
	Access_hash int64
	Title       string
}

func (e TL_channelForbidden) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_chatForbidden)
	var flags int32
	if e.Broadcast {
		flags |= (1 << 5)
	}
	if e.Megagroup {
		flags |= (1 << 8)
	}
	x.Int(flags)
	x.Int(e.Id)
	x.Long(e.Access_hash)
	x.String(e.Title)
	return x.buf
}

//inputChannelEmpty#ee8c1e86 = InputChannel;
const crc_inputChannelEmpty = 0xee8c1e86

type TL_inputChannelEmpty struct{}

func (e TL_inputChannelEmpty) encode() []byte {
	x := NewEncodeBuf(4)
	x.UInt(crc_inputChannelEmpty)
	return x.buf
}

//inputChannel#afeb712e channel_id:int access_hash:long = InputChannel;

const crc_inputChannel = 0xafeb712e

type TL_inputChannel struct {
	Channel_id  int32
	Access_hash int64
}

func (e TL_inputChannel) encode() []byte {
	x := NewEncodeBuf(16)
	x.UInt(crc_inputChannel)
	x.Int(e.Channel_id)
	x.Long(e.Access_hash)
	return x.buf
}

//replyKeyboardHide#a03e5b85 flags:# selective:flags.2?true = ReplyMarkup;
//replyKeyboardForceReply#f4108aa0 flags:# single_use:flags.1?true selective:flags.2?true = ReplyMarkup;
//replyKeyboardMarkup#3502758c flags:# resize:flags.0?true single_use:flags.1?true selective:flags.2?true rows:Vector<KeyboardButtonRow> = ReplyMarkup;
//replyInlineMarkup#48a30254 rows:Vector<KeyboardButtonRow> = ReplyMarkup;

//messageMediaEmpty#3ded6320 = MessageMedia;
//messageMediaPhoto#3d8ce53d photo:Photo caption:string = MessageMedia;
//messageMediaGeo#56e0d474 geo:GeoPoint = MessageMedia;
//messageMediaContact#5e7d2f39 phone_number:string first_name:string last_name:string user_id:int = MessageMedia;
//messageMediaUnsupported#9f84f49e = MessageMedia;
//messageMediaDocument#f3e02ea8 document:Document caption:string = MessageMedia;
//messageMediaWebPage#a32dd600 webpage:WebPage = MessageMedia;
//messageMediaVenue#7912b71f geo:GeoPoint title:string address:string provider:string venue_id:string = MessageMedia;
//messageMediaGame#fdb19008 game:Game = MessageMedia;
//messageMediaInvoice#84551347 flags:# shipping_address_requested:flags.1?true test:flags.3?true title:string description:string photo:flags.0?WebDocument receipt_msg_id:flags.2?int currency:string total_amount:long start_param:string = MessageMedia;

//messageFwdHeader#c786ddcb flags:# from_id:flags.0?int date:int channel_id:flags.1?int channel_post:flags.2?int = MessageFwdHeader;

//messageEmpty#83e5de54 id:int = Message;
const crc_messageEmpty = 0x83e5de54

type TL_messageEmpty struct {
	Id int32
}

func (e TL_messageEmpty) encode() []byte {
	x := NewEncodeBuf(8)
	x.UInt(crc_messageEmpty)
	x.Int(e.Id)

	return x.buf
}

//message#c09be45f flags:# out:flags.1?true mentioned:flags.4?true media_unread:flags.5?true silent:flags.13?true post:flags.14?true id:int from_id:flags.8?int to_id:Peer fwd_from:flags.2?MessageFwdHeader via_bot_id:flags.11?int reply_to_msg_id:flags.3?int date:int message:string media:flags.9?MessageMedia reply_markup:flags.6?ReplyMarkup entities:flags.7?Vector<MessageEntity> views:flags.10?int edit_date:flags.15?int = Message;
//messageService#9e19a1f6 flags:# out:flags.1?true mentioned:flags.4?true media_unread:flags.5?true silent:flags.13?true post:flags.14?true id:int from_id:flags.8?int to_id:Peer reply_to_msg_id:flags.3?int date:int action:MessageAction = Message;

//messages.messages#8c718e87 messages:Vector<Message> chats:Vector<Chat> users:Vector<User> = messages.Messages;
const crc_messages_messages = 0x8c718e87

type TL_messages_messages struct {
	Messages []TL
	Chats    []TL
	Users    []TL
}

func (e TL_messages_messages) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_messages_messages)
	x.Vector(e.Messages)
	x.Vector(e.Chats)
	x.Vector(e.Users)
	return x.buf
}

//messages.messagesSlice#b446ae3 count:int messages:Vector<Message> chats:Vector<Chat> users:Vector<User> = messages.Messages;

const crc_messages_messagesSlice = 0xb446ae3

type TL_messages_messagesSlice struct {
	Count    int32
	Messages []TL
	Chats    []TL
	Users    []TL
}

func (e TL_messages_messagesSlice) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_messages_messagesSlice)
	x.Int(e.Count)
	x.Vector(e.Messages)
	x.Vector(e.Chats)
	x.Vector(e.Users)

	return x.buf
}

//messages.channelMessages#99262e37 flags:# pts:int count:int messages:Vector<Message> chats:Vector<Chat> users:Vector<User> = messages.Messages;
const crc_messages_channelMessages = 0x99262e37

type TL_messages_channelMessages struct {
	Flags    int32
	Pts      int32
	Count    int32
	Messages []TL
	Chats    []TL
	Users    []TL
}

func (e TL_messages_channelMessages) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_messages_channelMessages)
	var flags int32
	x.Int(flags)
	x.Int(e.Pts)
	x.Int(e.Count)
	x.Vector(e.Messages)
	x.Vector(e.Chats)
	x.Vector(e.Users)
	return x.buf
}

//messages.getHistory#afa92846 peer:InputPeer offset_id:int offset_date:int add_offset:int limit:int max_id:int min_id:int = messages.Messages;

const crc_messages_getHistory = 0xafa92846

type TL_messages_getHistory struct {
	Peer        TL
	Offset_id   int32
	Offset_date int32
	Add_offset  int32
	Limit       int32
	Max_id      int32
	Min_id      int32
}

func (e TL_messages_getHistory) encode() []byte {
	x := NewEncodeBuf(512)
	x.UInt(crc_messages_getHistory)
	x.Bytes(e.Peer.encode())
	x.Int(e.Offset_id)
	x.Int(e.Offset_date)
	x.Int(e.Add_offset)
	x.Int(e.Limit)
	x.Int(e.Max_id)
	x.Int(e.Min_id)

	return x.buf
}
