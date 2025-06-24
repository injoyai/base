package types

import (
	"testing"
)

func TestBytes_Endian(t *testing.T) {
	bs := Bytes{11, 12, 13, 14, 15, 16, 17, 18}
	t.Log(bs.Endian("21").Bytes())       //{12,11,14,13,16,15,18,17}
	t.Log(bs.Endian("4321").Bytes())     //{14,13,12,11,18,17,16,15}
	t.Log(bs.Endian("43_1").Bytes())     //{14,13,11,18,17,15}
	t.Log(bs.Endian("87654321").Bytes()) //{18,17,16,15,14,13,12,11}
}

func TestBytes_Endian2(t *testing.T) {
	bs := Bytes{}
	t.Log(bs.Endian("21").Bytes())       //{}
	t.Log(bs.Endian("4321").Bytes())     //{}
	t.Log(bs.Endian("43_1").Bytes())     //{}
	t.Log(bs.Endian("87654321").Bytes()) //{}
}

func TestBytes_Endian3(t *testing.T) {
	bs := Bytes{11, 12, 13}
	t.Log(bs.Endian("21").Bytes())       //[12 11 0]
	t.Log(bs.Endian("4321").Bytes())     //[0 13 12]
	t.Log(bs.Endian("_3_1").Bytes())     //[0 13]
	t.Log(bs.Endian("87654321").Bytes()) //[0 0 0]
}

func TestBytes_Endian4(t *testing.T) {
	bs := Bytes(nil)
	t.Log(bs.Endian("21").Bytes())       //{}
	t.Log(bs.Endian("4321").Bytes())     //{}
	t.Log(bs.Endian("43_1").Bytes())     //{}
	t.Log(bs.Endian("87654321").Bytes()) //{}
}
