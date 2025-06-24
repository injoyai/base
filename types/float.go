package types

import (
	"bytes"
	"encoding/binary"
	"github.com/injoyai/conv"
	"math"
	"strconv"
)

type (
	F64   = Float64
	Float = Float64
	F32   = Float32
)

type Float64 float64

func (this Float64) Float() float64 { return float64(this) }

func (this Float64) Int() int { return int(this) }

func (this Float64) Int64() int64 { return int64(this) }

func (this Float64) String() string { return strconv.FormatFloat(this.Float(), 'f', -1, 64) }

func (this Float64) Bytes() Bytes {
	u := math.Float64bits(this.Float())
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, u)
	return bytesBuffer.Bytes()
}

func (this Float64) Decimals(d ...int) float64 {
	b := math.Pow10(conv.Default(2, d...))
	return float64(int64(this.Float()*b+0.5)) / b
}

type Float32 float32

func (this Float32) Float() float32 { return float32(this) }

func (this Float32) Int() int { return int(this) }

func (this Float32) Int64() int64 { return int64(this) }

func (this Float32) String() string { return strconv.FormatFloat(float64(this.Float()), 'f', -1, 64) }

func (this Float32) Bytes() Bytes {
	u := math.Float32bits(this.Float())
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, u)
	return bytesBuffer.Bytes()
}

func (this Float32) Decimals(d ...int) float32 {
	b := float32(math.Pow10(conv.Default(2, d...)))
	return float32(int64(this.Float()*b+0.5)) / b
}
