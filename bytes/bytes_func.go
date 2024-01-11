package bytes

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"github.com/injoyai/conv"
	"math"
	"strconv"
)

func Copy(bs []byte) []byte {
	cp := make([]byte, len(bs))
	copy(cp, bs)
	return cp
}

// Sum 字节校验和
func Sum(bs []byte) byte {
	b := byte(0)
	for _, v := range bs {
		b += v
	}
	return b
}

// Reverse 字节倒序
func Reverse(bs []byte) []byte {
	x := make([]byte, len(bs))
	for i, v := range bs {
		x[len(bs)-i-1] = v
	}
	return x
}

func Upper(bs []byte) []byte {
	return bytes.ToUpper(bs)
}

func Lower(bs []byte) []byte {
	return bytes.ToLower(bs)
}

func UTF8(bs []byte) string {
	return string(bs)
}

func ASCII(bs []byte) string {
	return string(bs)
}

func HEX(bs []byte) string {
	return hex.EncodeToString(bs)
}

func Base64(bs []byte) string {
	return base64.StdEncoding.EncodeToString(bs)
}

func BIN(bs []byte) string {
	return conv.BINStr(bs)
}

func HEXBase64(bs []byte) string {
	return Base64([]byte(HEX(bs)))
}

func HEXToInt(bs []byte) (int, error) {
	return strconv.Atoi(HEX(bs))
}

func HEXToFloat(bs []byte, decimals int) (float64, error) {
	i, err := HEXToInt(bs)
	return float64(i) / math.Pow10(decimals), err
}

func ASCIIToInt(bs []byte) (int, error) {
	return strconv.Atoi(ASCII(bs))
}

func ASCIIToFloat(bs []byte, decimals int) (float64, error) {
	i, err := strconv.Atoi(ASCII(bs))
	return float64(i) / math.Pow10(decimals), err
}

func Uint8(bs []byte) uint8 {
	return uint8(Uint64(bs))
}

func Uint16(bs []byte) uint16 {
	return uint16(Uint64(bs))
}

func Uint32(bs []byte) uint32 {
	return uint32(Uint64(bs))
}

func Uint64(bs []byte) uint64 {
	cp := Copy(bs)
	for len(cp) < 8 {
		cp = append([]byte{0}, cp...)
	}
	return binary.BigEndian.Uint64(cp)
}

func Int8(bs []byte) int8 {
	return int8(Uint64(bs))
}

func Int16(bs []byte) int16 {
	return int16(Uint64(bs))
}

func Int(bs []byte) int {
	return int(Uint64(bs))
}

func Int32(bs []byte) int32 {
	return int32(Uint64(bs))
}

func Int64(bs []byte) int64 {
	return int64(Uint64(bs))
}

func AddByte(bs []byte, b byte) []byte {
	result := make([]byte, len(bs))
	for _, v := range bs {
		result = append(result, v+b)
	}
	return result
}

func SubByte(bs []byte, b byte) []byte {
	result := make([]byte, len(bs))
	for _, v := range bs {
		result = append(result, v-b)
	}
	return result
}

// Sub0x33ReverseHEXToInt DLT645专用
func Sub0x33ReverseHEXToInt(bs []byte) (int, error) {
	bs = SubByte(bs, 0x33)
	bs = Reverse(bs)
	return HEXToInt(bs)
}

// Sub0x33ReverseHEXToFloat DLT645专用
func Sub0x33ReverseHEXToFloat(bs []byte, decimals int) (float64, error) {
	bs = SubByte(bs, 0x33)
	bs = Reverse(bs)
	return HEXToFloat(bs, decimals)
}
