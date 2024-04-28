package bytes

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"github.com/injoyai/conv"
	"io"
	"math"
	"strconv"
)

func Copy(bs []byte) []byte {
	cp := make([]byte, len(bs))
	copy(cp, bs)
	return cp
}

func Equal(bs1, bs2 []byte) bool {
	if (bs1 == nil) != (bs2 == nil) {
		return false
	}
	return bytes.Equal(bs1, bs2)
}

// Sum 字节校验和
func Sum(bs []byte) byte {
	b := byte(0)
	for _, v := range bs {
		b += v
	}
	return b
}

func WriteTo(w io.Writer, bs []byte) (int, error) {
	return w.Write(bs)
}

// Reverse 字节倒序
func Reverse(bs []byte) []byte {
	x := make([]byte, len(bs))
	for i, v := range bs {
		x[len(bs)-i-1] = v
	}
	return x
}

// Upper 字节转大写
func Upper(bs []byte) []byte {
	return bytes.ToUpper(bs)
}

// Lower 字节转小写
func Lower(bs []byte) []byte {
	return bytes.ToLower(bs)
}

// UTF8 []{0x31,0x32} >>> "12"
func UTF8(bs []byte) string {
	return string(bs)
}

// ASCII 等效UTF8
func ASCII(bs []byte) string {
	return string(bs)
}

// HEX []{0x01,0x02} >>> "0102"
func HEX(bs []byte) string {
	return hex.EncodeToString(bs)
}

// Base64 same base64.StdEncoding.EncodeToString
func Base64(bs []byte) string {
	return base64.StdEncoding.EncodeToString(bs)
}

// BIN []{0x01,0x02} >>> "0000000100000002"
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

func UTF8ToInt(bs []byte) (int, error) {
	return strconv.Atoi(UTF8(bs))
}

// ASCIIToInt 历史问题,先保留该命名
func ASCIIToInt(bs []byte) (int, error) {
	return UTF8ToInt(bs)
}

func UTF8ToFloat(bs []byte, decimals int) (float64, error) {
	i, err := strconv.Atoi(ASCII(bs))
	return float64(i) / math.Pow10(decimals), err
}

// ASCIIToFloat 历史问题,先保留该命名
func ASCIIToFloat(bs []byte, decimals int) (float64, error) {
	return UTF8ToFloat(bs, decimals)
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
	for i := range bs {
		result[i] = bs[i] + b
	}
	return result
}

func SubByte(bs []byte, b byte) []byte {
	result := make([]byte, len(bs))
	for i := range bs {
		result[i] = bs[i] - b
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
