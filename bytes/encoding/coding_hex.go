package coding

import "encoding/hex"

type HEX struct{}

func (this HEX) Decode(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func (this HEX) Encode(bs []byte) string {
	return hex.EncodeToString(bs)
}

// DecodeHEX hex.DecodeString
func DecodeHEX(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

// EncodeHEX hex.EncodeToString
func EncodeHEX(bs []byte) string {
	return hex.EncodeToString(bs)
}
