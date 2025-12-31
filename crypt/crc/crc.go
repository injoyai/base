package crc

import (
	"encoding/base64"
	"encoding/hex"
)

//========================================crc16========================================

func Encrypt16(bs []byte, params ...Param16) []byte {
	param := CRC16_MODBUS
	if len(params) > 0 {
		param = params[0]
	}
	num := Checksum16(bs, MakeTable16(param))
	return []byte{byte(num >> 8), byte(num)}
}

func Encrypt16String(bs []byte, params ...Param16) string {
	return string(Encrypt16(bs, params...))
}

func Encrypt16HEX(bs []byte, params ...Param16) string {
	return hex.EncodeToString(Encrypt16(bs, params...))
}

func Encrypt16Base64(bs []byte, params ...Param16) string {
	return base64.StdEncoding.EncodeToString(Encrypt16(bs, params...))
}

//========================================crc8========================================

func Encrypt8(bs []byte, params ...Param8) byte {
	param := CRC8
	if len(params) > 0 {
		param = params[0]
	}
	return Checksum8(bs, MakeTable8(param))
}
