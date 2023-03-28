package pkg

import (
	"encoding/hex"
	"testing"
)

func TestPkg_Bytes(t *testing.T) {
	t.Log(New(20, []byte{0, 1, 2, 3, 5}).Bytes().HEX()) //88880011801400010203051348a34e8989

	{
		s := "88880011801400010203051348a34e8989"
		bs, err := hex.DecodeString(s)
		if err != nil {
			t.Error(err)
			return
		}
		p, err := Decode(bs)
		if err != nil {
			t.Error(err)
			return
		}
		t.Logf("%#v", p)
	}
}
