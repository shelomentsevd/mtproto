package core

import "math/big"

type TL interface {
	encode() []byte
}

const crc_vector = 0x1cb5c415 // Processed manually

const crc_msg_container = 0x73f1f8dc

type TL_msg_container struct {
	items []TL_MT_message
}

type TL_MT_message struct {
	msg_id int64
	seq_no int32
	size   int32
	data   interface{}
}

const crc_req_pq = 0x60469778

type TL_req_pq struct {
	nonce []byte
}

const crc_p_q_inner_data = 0x83c95aec

type TL_p_q_inner_data struct {
	pq           *big.Int
	p            *big.Int
	q            *big.Int
	nonce        []byte
	server_nonce []byte
	new_nonce    []byte
}

const crc_req_DH_params = 0xd712e4be

type TL_req_DH_params struct {
	nonce        []byte
	server_nonce []byte
	p            *big.Int
	q            *big.Int
	fp           uint64
	encdata      []byte
}

const crc_client_DH_inner_data = 0x6643b654

type TL_client_DH_inner_data struct {
	nonce        []byte
	server_nonce []byte
	retry        int64
	g_b          *big.Int
}

const crc_set_client_DH_params = 0xf5045f1f

type TL_set_client_DH_params struct {
	nonce        []byte
	server_nonce []byte
	encdata      []byte
}

const crc_resPQ = 0x05162463

type TL_resPQ struct {
	nonce        []byte
	server_nonce []byte
	pq           *big.Int
	fingerprints []int64
}

const crc_server_DH_params_ok = 0xd0e8075c

type TL_server_DH_params_ok struct {
	nonce            []byte
	server_nonce     []byte
	encrypted_answer []byte
}

const crc_server_DH_params_fail = 0x79cb045d

type TL_server_DH_params_fail struct {
	nonce          []byte
	server_nonce   []byte
	new_nonce_hash []byte
}

const crc_server_DH_inner_data = 0xb5890dba

type TL_server_DH_inner_data struct {
	nonce        []byte
	server_nonce []byte
	g            int32
	dh_prime     *big.Int
	g_a          *big.Int
	server_time  int32
}

const crc_new_session_created = 0x9ec20908

type TL_new_session_created struct {
	first_msg_id int64
	unique_id    int64
	server_salt  []byte
}

const crc_bad_server_salt = 0xedab447b

type TL_bad_server_salt struct {
	bad_msg_id      int64
	bad_msg_seqno   int32
	error_code      int32
	new_server_salt []byte
}

const crc_bad_msg_notification = 0xa7eff811

type TL_bad_msg_notification struct {
	bad_msg_id    int64
	bad_msg_seqno int32
	error_code    int32
}

const crc_msgs_ack = 0x62d6b459

type TL_msgs_ack struct {
	msgIds []int64
}

const crc_rpc_result = 0xf35c6d01

type TL_rpc_result struct {
	req_msg_id int64
	obj        interface{}
}

const crc_rpc_error = 0x2144ca19

type TL_rpc_error struct {
	error_code    int32
	error_message string
}

const crc_dh_gen_ok = 0x3bcbf734

type TL_dh_gen_ok struct {
	nonce           []byte
	server_nonce    []byte
	new_nonce_hash1 []byte
}

const crc_ping = 0x7abe77ec

type TL_ping struct {
	ping_id int64
}

const crc_pong = 0x347773c5

type TL_pong struct {
	msg_id  int64
	ping_id int64
}
