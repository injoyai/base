package coding

import (
	"encoding/base64"
)

type Base64 struct{}

func (this Base64) Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func (this Base64) Encode(bs []byte) string {
	return base64.StdEncoding.EncodeToString(bs)
}

// DecodeBase64 base64.StdEncoding.DecodeString
func DecodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// EncodeBase64 base64.StdEncoding.EncodeToString
func EncodeBase64(bs []byte) string {
	return base64.StdEncoding.EncodeToString(bs)
}
