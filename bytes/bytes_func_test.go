package bytes

import (
	"encoding/hex"
	"testing"
)

func TestUint64(t *testing.T) {
	t.Log(Uint64([]byte{01, 02}))
	t.Log(Uint64([]byte{0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 03}))
	t.Log(Int64([]byte{01, 02}))
	t.Log(Int64([]byte{0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 03}))
}

func TestAddByte(t *testing.T) {
	t.Log(hex.EncodeToString(AddByte([]byte{0x34, 0x32, 0x33}, 0x33)))
	t.Log(hex.EncodeToString(SubByte([]byte{0x34, 0x32, 0x33}, 0x33)))
}
