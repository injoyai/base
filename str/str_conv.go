package str

import (
	"encoding/json"
	"github.com/injoyai/conv"
)

var (
	Int     = conv.Int
	Uint8   = conv.Uint8
	Uint16  = conv.Uint16
	Uint32  = conv.Uint32
	Uint64  = conv.Uint64
	Int8    = conv.Int8
	Int16   = conv.Int16
	Int32   = conv.Int32
	Int64   = conv.Int64
	Float32 = conv.Float32
	Float64 = conv.Float64
	Bool    = conv.Bool
	Select  = conv.Select[string]
)

func Interface(s string) any {
	s = `{"data":` + s + `}`
	m := make(map[string]any)
	_ = json.Unmarshal([]byte(s), &m)
	return m["data"]
}
