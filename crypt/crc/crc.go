package crc

import (
	"github.com/injoyai/base/types"
)

//========================================crc16========================================

func Encrypt16(bs []byte, params ...Param16) types.Bytes {
	param := CRC16_MODBUS
	if len(params) > 0 {
		param = params[0]
	}
	num := Checksum16(bs, MakeTable16(param))
	return []byte{byte(num >> 8), byte(num)}
}

func Encrypt16Bytes(bs []byte, params ...Param16) []byte { return Encrypt16(bs, params...).Bytes() }

func Encrypt16ASCII(bs []byte, params ...Param16) string { return Encrypt16(bs, params...).ASCII() }

func Encrypt16HEX(bs []byte, params ...Param16) string { return Encrypt16(bs, params...).HEX() }

func Encrypt16Base64(bs []byte, params ...Param16) string { return Encrypt16(bs, params...).Base64() }

//========================================crc8========================================

func Encrypt8(bs []byte, params ...Param8) types.Bytes {
	param := CRC8
	if len(params) > 0 {
		param = params[0]
	}
	return []byte{Checksum8(bs, MakeTable8(param))}
}

func Encrypt8Byte(bs []byte, params ...Param8) byte { return Encrypt8(bs, params...).Bytes()[0] }

func Encrypt8ASCII(bs []byte, params ...Param8) string { return Encrypt8(bs, params...).ASCII() }

func Encrypt8HEX(bs []byte, params ...Param8) string { return Encrypt8(bs, params...).HEX() }

func Encrypt8Base64(bs []byte, params ...Param8) string { return Encrypt8(bs, params...).Base64() }
