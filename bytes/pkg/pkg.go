package pkg

import (
	"fmt"
	"github.com/injoyai/base/bytes"
	"github.com/injoyai/conv"
	"hash/crc32"
)

/*
包构成(大端):
.===================================.
|构成	|字节	|类型	|说明		|
|-----------------------------------|
|帧头 	|1字节 	|Byte	|固定0x8888	|
|-----------------------------------|
|帧长  	|2字节	|HEX	|总字节长度	|
|-----------------------------------|
|帧类型	|1字节	|Bin	|详见帧类型	|
|-----------------------------------|
|消息号	|1字节	|Byte	|消息id		|
|-----------------------------------|
|内容	|可变	|Byte	|数据内容	|
|-----------------------------------|
|校验和	|4字节	|Byte	|crc IEEE	|
|-----------------------------------|
|帧尾 	|1字节	|Byte	|固定0x8989	|
^===================================^

包类型:
.===================================================================================.
|bit0				|bit1			|bit2	|bit3	|bit4	|bit5	|bit6	|bit7	|
|-----------------------------------------------------------------------------------|
|数据方向0请求,1响应	|1测试通讯,无内容	|预留	|预留	|预留	|预留	|预留	|预留	|						|
^===================================================================================^
*/

var (
	start = []byte{0x88, 0x88} //帧头
	end   = []byte{0x89, 0x89} //帧尾
)

const (
	minLength      = 12
	bitBack   byte = 0x80
	bitPing   byte = 0x40
)

func Ping() []byte {
	//01000000
	return (&Pkg{Type: bitPing}).Bytes()
}

func Pong() []byte {
	//11000000
	return (&Pkg{Type: bitBack + bitPing}).Bytes()
}

func New(msgID uint8, data []byte) *Pkg {
	return &Pkg{
		Type:  bitBack,
		MsgID: msgID,
		Data:  data,
	}
}

type Pkg struct {
	Type  uint8
	MsgID uint8
	Data  []byte
}

func (this *Pkg) String() string {
	return this.Bytes().HEX()
}

func (this *Pkg) Bytes() bytes.Entity {
	data := []byte(nil)
	data = append(data, start...)
	length := len(this.Data) + minLength
	data = append(data, byte(length>>8), byte(length))
	data = append(data, this.Type)
	data = append(data, this.MsgID)
	data = append(data, this.Data...)
	data = append(data, conv.Bytes(crc32.ChecksumIEEE(data))...)
	data = append(data, end...)
	return data
}

func (this *Pkg) Resp(bs []byte) *Pkg {
	this.Type += bitBack
	this.Data = bs
	return this
}

func (this *Pkg) IsCall() bool {
	return this.Type>>7 == 0
}

func (this *Pkg) IsBack() bool {
	return this.Type>>7 == 1
}

func (this *Pkg) IsPing() bool {
	return this.Type>>6 == 1
}

func (this *Pkg) IsPong() bool {
	return this.Type>>6 == 3
}

func Decode(bs []byte) (*Pkg, error) {

	//校验基础数据长度
	if len(bs) <= 10 {
		return nil, fmt.Errorf("数据长度小于(%d)", minLength)
	}

	//校验帧头
	if bs[0] != start[0] && bs[1] != start[1] {
		return nil, fmt.Errorf("帧头错误,需要(%x),得到(%x)", start, bs[:2])
	}

	//获取总数据长度
	length := conv.Int(bs[2:4])

	//校验总长度
	if len(bs) != length {
		return nil, fmt.Errorf("数据总长度错误,需要(%d),得到(%d)", length, len(bs))
	}

	//校验crc32
	if crc1, crc2 := crc32.ChecksumIEEE(bs[:length-6]), conv.Uint32(bs[length-6:length-2]); crc1 != crc2 {
		return nil, fmt.Errorf("数据CRC校验错误,需要(%x),得到(%x)", crc1, crc2)
	}

	//校验帧尾
	if bs[length-2] != end[0] && bs[length-1] != bs[1] {
		return nil, fmt.Errorf("帧尾错误,需要(%x),得到(%x)", end, bs[length-2:])
	}

	return &Pkg{
		Type:  bs[4],
		MsgID: bs[5],
		Data:  bs[6 : length-6],
	}, nil

}
