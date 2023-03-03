package hmac

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"

	"github.com/injoyai/base/bytes/crypt"
)

func Sha256(str string, secret string) string {
	return crypt.Hmac(sha256.New, []byte(secret)).EncryptHEXBase64([]byte(str))
}

func Sha1(str string, secret string) string {
	return crypt.Hmac(sha1.New, []byte(secret)).EncryptHEXBase64([]byte(str))
}

func Sha512(str string, secret string) string {
	return crypt.Hmac(sha512.New, []byte(secret)).EncryptHEXBase64([]byte(str))
}
