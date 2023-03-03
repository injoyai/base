package sha

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"github.com/injoyai/base/bytes/crypt"
)

func Encrypt256(s string) string {
	return crypt.New(sha256.New()).EncryptBase64([]byte(s))
}

func Encrypt256Bytes(s string) []byte {
	return crypt.New(sha256.New()).EncryptBytes([]byte(s))
}

func Encrypt256ASCII(s string) string {
	return crypt.New(sha256.New()).EncryptASCII([]byte(s))
}

func Encrypt256HEX(s string) string {
	return crypt.New(sha256.New()).EncryptHEX([]byte(s))
}

func Encrypt256Base64(s string) string {
	return crypt.New(sha256.New()).EncryptBase64([]byte(s))
}

func Encrypt1(s string) string {
	return crypt.New(sha1.New()).EncryptBase64([]byte(s))
}

func Encrypt512(s string) string {
	return crypt.New(sha512.New()).EncryptBase64([]byte(s))
}
