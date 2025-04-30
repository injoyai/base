package types

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

type Bs = Bytes

type Bytes []byte

// Len 数据长度
func (this Bytes) Len() int {
	return len(this)
}

func (this Bytes) Cap() int {
	return cap(this)
}

func (this Bytes) Error() string {
	return this.String()
}

func (this Bytes) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(this)
	return int64(n), err
}

func (this Bytes) Sum() byte {
	b := byte(0)
	for _, v := range this {
		b += v
	}
	return b
}

func (this Bytes) Copy() Bytes {
	cp := make([]byte, len(this))
	copy(cp, this)
	return cp
}

func (this Bytes) Equal(bs Bytes) bool {
	if (this == nil) != (bs == nil) {
		return false
	}
	return bytes.Equal(this, bs)
}

func (this Bytes) Upper() Bytes {
	return bytes.ToUpper(this)
}

func (this Bytes) Lower() Bytes {
	return bytes.ToLower(this)
}

// String []{0x31,0x32} >>> "12"
func (this Bytes) String() string {
	return string(this)
}

// UTF8 []{0x31,0x32} >>> "12"
func (this Bytes) UTF8() string {
	return string(this)
}

// ASCII []{0x31,0x32} >>> "12"
func (this Bytes) ASCII() string {
	return string(this)
}

// HEX []{0x01,0x02} >>> "0102"
func (this Bytes) HEX() string {
	return hex.EncodeToString(this)
}

// Base64 same base64.StdEncoding.EncodeToString
func (this Bytes) Base64() string {
	return base64.StdEncoding.EncodeToString(this)
}

// HEXBase64 HEX() then Base64()
func (this Bytes) HEXBase64() string {
	return Bytes(this.HEX()).Base64()
}

// Bytes 字节数组
func (this Bytes) Bytes() []byte {
	return this
}

// Reader io.Reader
func (this Bytes) Reader() io.Reader {
	return bytes.NewReader(this.Bytes())
}

// Buffer bytes.Buffer
func (this Bytes) Buffer() *bytes.Buffer {
	return bytes.NewBuffer(this.Bytes())
}

// GetFirst 获取第一个元素,不存在返回0
func (this Bytes) GetFirst() byte {
	if this.Len() > 0 {
		return this[0]
	}
	return 0
}

// GetLast 获取最后一个元素,不存在则返回0
func (this Bytes) GetLast() byte {
	if this.Len() > 0 {
		return this[this.Len()-1]
	}
	return 0
}

// Int64 字节数组转int64 大端模式
func (this Bytes) Int64() int64 {
	return int64(this.Uint64())
}

// Uint64 字节数组转uint64 大端模式
func (this Bytes) Uint64() uint64 {
	cp := this.Copy()
	for len(cp) < 8 {
		cp = append([]byte{0}, cp...)
	}
	return binary.BigEndian.Uint64(cp)
}

// BINStr 字节转2进制字符串
func (this Bytes) BINStr() string {
	return conv.BINStr(this)
}

// BIN 字节转2进制字符串
func (this Bytes) BIN() string {
	return conv.BINStr(this)
}

// Append just append
func (this Bytes) Append(b ...byte) Bytes {
	return append(this, b...)
}

// ASCIIToInt []{0x31,0x32} >>> 12
func (this Bytes) ASCIIToInt() (int, error) {
	return strconv.Atoi(this.ASCII())
}

// ASCIIToFloat64 字节ascii编码再转int,再转float64
func (this Bytes) ASCIIToFloat64(decimals int) (float64, error) {
	i, err := strconv.Atoi(this.ASCII())
	return float64(i) / math.Pow10(decimals), err
}

// HEXToInt []{0x01,0x02} >>> 102
func (this Bytes) HEXToInt() (int, error) {
	return strconv.Atoi(this.HEX())
}

// HEXToFloat64 字节hex编码再转int,再转float64
func (this Bytes) HEXToFloat64(decimals int) (float64, error) {
	i, err := this.HEXToInt()
	return float64(i) / math.Pow10(decimals), err
}

// Reverse 倒序
func (this Bytes) Reverse() Bytes {
	x := make([]byte, len(this))
	for i, v := range this {
		x[len(this)-i-1] = v
	}
	return x
}

// ReverseASCII 倒序再ASCII
func (this Bytes) ReverseASCII() string {
	return this.Reverse().ASCII()
}

// ReverseHEX 倒序再hex
func (this Bytes) ReverseHEX() string {
	return this.Reverse().HEX()
}

// ReverseBase64 倒序再base64
func (this Bytes) ReverseBase64() string {
	return this.Reverse().Base64()
}

// SubByte 每个字节减sub
func (this Bytes) SubByte(sub byte) Bytes {
	result := make([]byte, len(this))
	for i := range this {
		result[i] = this[i] - sub
	}
	return result
}

// AddByte 每个字节加add
func (this Bytes) AddByte(add byte) Bytes {
	result := make([]byte, len(this))
	for i := range this {
		result[i] = this[i] + add
	}
	return result
}

// Sub0x33 每个字节减0x33
func (this Bytes) Sub0x33() Bytes {
	return this.SubByte(0x33)
}

// Add0x33 每个字节加0x33
func (this Bytes) Add0x33() Bytes {
	return this.AddByte(0x33)
}

// Sub0x33ReverseHEXToInt DLT645协议流程,先减0x33,再倒序,再转hex,再转int
func (this Bytes) Sub0x33ReverseHEXToInt() (int, error) {
	return this.Sub0x33().Reverse().HEXToInt()
}

// Sub0x33ReverseHEXToFloat DLT645协议流程,先减0x33,再倒序,再转hex,再转float64
func (this Bytes) Sub0x33ReverseHEXToFloat(decimals int) (float64, error) {
	return this.Sub0x33().Reverse().HEXToFloat64(decimals)
}
