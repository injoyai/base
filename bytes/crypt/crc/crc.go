package crc

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
)

func Encrypt16ASCII(bs []byte) string {
	return string(Encrypt16Bytes(bs))
}

func Encrypt16HEX(bs []byte) string {
	return hex.EncodeToString(Encrypt16Bytes(bs))
}

func Encrypt16Base64(bs []byte) string {
	return base64.StdEncoding.EncodeToString(Encrypt16Bytes(bs))
}

func Encrypt16Bytes(bs []byte, params ...Param16) []byte {
	num := Encrypt16(bs, params...)
	int16buf := new(bytes.Buffer)
	_ = binary.Write(int16buf, binary.LittleEndian, num)
	for int16buf.Len() < 2 {
		int16buf.Write([]byte{0x00})
	}
	return int16buf.Bytes()
}

func Encrypt16(bs []byte, params ...Param16) uint16 {
	param := CRC16_MODBUS
	if len(params) > 0 {
		param = params[0]
	}
	return Checksum16(bs, MakeTable16(param))
}

func Encrypt8ASCII(bs []byte) string {
	return string([]byte{Encrypt8(bs)})
}

func Encrypt8HEX(bs []byte) string {
	return hex.EncodeToString([]byte{Encrypt8(bs)})
}

func Encrypt8Base64(bs []byte) string {
	return base64.StdEncoding.EncodeToString([]byte{Encrypt8(bs)})
}

func Encrypt8(bs []byte, params ...Param8) uint8 {
	param := CRC8
	if len(params) > 0 {
		param = params[0]
	}
	return Checksum8(bs, MakeTable8(param))
}
