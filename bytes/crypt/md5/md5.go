package md5

import (
	"crypto/md5"
	"github.com/injoyai/base/bytes/crypt"
)

//========================================md5========================================

// Encrypt
// md5加密,默认hex
// 000000,加密后为670b14728ad9902aecba32e22fa4f6bd
func Encrypt(s string) string       { return encrypt().EncryptHEX([]byte(s)) }
func EncryptBytes(s string) []byte  { return encrypt().EncryptBytes([]byte(s)) }
func EncryptASCII(s string) string  { return encrypt().EncryptASCII([]byte(s)) }
func EncryptHEX(s string) string    { return encrypt().EncryptHEX([]byte(s)) }
func EncryptBase64(s string) string { return encrypt().EncryptBase64([]byte(s)) }
func encrypt() *crypt.Entity        { return crypt.New(md5.New()) }

//========================================hmac md5========================================

func Hmac(s string, key string) string       { return hmac(key).EncryptHEX([]byte(s)) }
func HmacBytes(s string, key string) []byte  { return hmac(key).EncryptBytes([]byte(s)) }
func HmacASCII(s string, key string) string  { return hmac(key).EncryptASCII([]byte(s)) }
func HmacHEX(s string, key string) string    { return hmac(key).EncryptHEX([]byte(s)) }
func HmacBase64(s string, key string) string { return hmac(key).EncryptBase64([]byte(s)) }
func hmac(key string) *crypt.Entity          { return crypt.Hmac(md5.New, []byte(key)) }
