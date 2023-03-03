package md5

import (
	"crypto/md5"
	"github.com/injoyai/base/bytes/crypt"
)

//Encrypt
//md5加密,默认hex
//000000,加密后为670b14728ad9902aecba32e22fa4f6bd
func Encrypt(s string) string {
	return crypt.New(md5.New()).EncryptHEX([]byte(s))
}

func EncryptHEX(s string) string {
	return crypt.New(md5.New()).EncryptHEX([]byte(s))
}

func EncryptASCII(s string) string {
	return crypt.New(md5.New()).EncryptASCII([]byte(s))
}

func EncryptBase64(s string) string {
	return crypt.New(md5.New()).EncryptBase64([]byte(s))
}

func EncryptBytes(s string) []byte {
	return crypt.New(md5.New()).EncryptBytes([]byte(s))
}

/*

 */

// Hmac
//HmacMD5,默认hex
func Hmac(s string, key string) string {
	return crypt.Hmac(md5.New, []byte(key)).EncryptHEX([]byte(s))
}

func HmacBytes(s string, key string) []byte {
	return crypt.Hmac(md5.New, []byte(key)).EncryptBytes([]byte(s))
}

func HmacASCII(s string, key string) string {
	return crypt.Hmac(md5.New, []byte(key)).EncryptASCII([]byte(s))
}

func HmacHEX(s string, key string) string {
	return crypt.Hmac(md5.New, []byte(key)).EncryptHEX([]byte(s))
}

func HmacBase64(s string, key string) string {
	return crypt.Hmac(md5.New, []byte(key)).EncryptBase64([]byte(s))
}
