package bytes

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"github.com/injoyai/conv"
	"io"
	"strconv"
)

type Entity []byte

// Len 数据长度
func (this Entity) Len() int {
	return len(this)
}

func (this Entity) Cap() int {
	return cap(this)
}

func (this Entity) Error() string {
	return this.String()
}

// String []{0x31,0x32} >>> "12"
func (this Entity) String() string {
	return string(this)
}

// ASCII []{0x31,0x32} >>> "12"
func (this Entity) ASCII() string {
	return string(this)
}

// HEX []{0x01,0x02} >>> "0102"
func (this Entity) HEX() string {
	return hex.EncodeToString(this)
}

// Base64 same base64.StdEncoding.EncodeToString
func (this Entity) Base64() string {
	return base64.StdEncoding.EncodeToString(this)
}

// HEXBase64 HEX() then Base64()
func (this Entity) HEXBase64() string {
	return Entity(this.HEX()).Base64()
}

// Bytes 字节数组
func (this Entity) Bytes() []byte {
	return this
}

// Reader io.Reader
func (this Entity) Reader() io.Reader {
	return bytes.NewReader(this.Bytes())
}

// Buffer bytes.Buffer
func (this Entity) Buffer() *bytes.Buffer {
	return bytes.NewBuffer(this.Bytes())
}

// SumByte 累加转byte
func (this Entity) SumByte() byte {
	return Sum(this)
}

// GetFirst 获取第一个元素,不存在返回0
func (this Entity) GetFirst() byte {
	if this.Len() > 0 {
		return this[0]
	}
	return 0
}

// GetLast 获取最后一个元素,不存在则返回0
func (this Entity) GetLast() byte {
	if this.Len() > 0 {
		return this[this.Len()-1]
	}
	return 0
}

// Int64 字节数组转int64 大端模式
func (this Entity) Int64() int64 {
	return Int64(this.Bytes())
}

// Uint64 字节数组转uint64 大端模式
func (this Entity) Uint64() uint64 {
	return Uint64(this.Bytes())
}

// BINStr 字节转2进制字符串
func (this Entity) BINStr() string {
	return conv.BINStr(this)
}

// Append just append
func (this Entity) Append(b ...byte) Entity {
	return append(this, b...)
}

// ASCIIToInt []{0x31,0x32} >>> 12
func (this Entity) ASCIIToInt() (int, error) {
	return strconv.Atoi(this.ASCII())
}

// ASCIIToFloat64 字节ascii编码再转int,再转float64
func (this Entity) ASCIIToFloat64(decimals int) (float64, error) {
	return ASCIIToFloat(this, decimals)
}

// HEXToInt []{0x01,0x02} >>> 102
func (this Entity) HEXToInt() (int, error) {
	return strconv.Atoi(this.HEX())
}

// HEXToFloat64 字节hex编码再转int,再转float64
func (this Entity) HEXToFloat64(decimals int) (float64, error) {
	return HEXToFloat(this, decimals)
}

// Reverse 倒序
func (this Entity) Reverse() Entity {
	return Reverse(this)
}

// ReverseASCII 倒序再ASCII
func (this Entity) ReverseASCII() string {
	return this.Reverse().ASCII()
}

// ReverseHEX 倒序再hex
func (this Entity) ReverseHEX() string {
	return this.Reverse().HEX()
}

// ReverseBase64 倒序再base64
func (this Entity) ReverseBase64() string {
	return this.Reverse().Base64()
}

// SubByte 每个字节减sub
func (this Entity) SubByte(sub byte) (result Entity) {
	return SubByte(this, sub)
}

// AddByte 每个字节加add
func (this Entity) AddByte(add byte) Entity {
	return AddByte(this, add)
}

// Sub0x33 每个字节减0x33
func (this Entity) Sub0x33() Entity {
	return this.SubByte(0x33)
}

// Add0x33 每个字节加0x33
func (this Entity) Add0x33() Entity {
	return this.AddByte(0x33)
}

// Sub0x33ReverseHEXToInt DLT645协议流程,先减0x33,再倒序,再转hex,再转int
func (this Entity) Sub0x33ReverseHEXToInt() (int, error) {
	return Sub0x33ReverseHEXToInt(this)
}

// Sub0x33ReverseHEXToFloat DLT645协议流程,先减0x33,再倒序,再转hex,再转float64
func (this Entity) Sub0x33ReverseHEXToFloat(decimals int) (float64, error) {
	return Sub0x33ReverseHEXToFloat(this, decimals)
}

// Sub0x33ReverseHEXToFloat64 DLT645协议流程,先减0x33,再倒序,再转hex,再转float64
func (this Entity) Sub0x33ReverseHEXToFloat64(decimals int) (float64, error) {
	return Sub0x33ReverseHEXToFloat(this, decimals)
}
