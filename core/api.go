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
	default:
		m.err = fmt.Errorf("Unknown constructor: \u002508x", constructor)
		return nil
	}
	return
}
