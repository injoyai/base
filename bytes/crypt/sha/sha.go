package sha

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"github.com/injoyai/base/bytes"
	"github.com/injoyai/base/bytes/crypt"
)

//========================================Sha1========================================

func Encrypt1(bs []byte) bytes.Entity { return crypt.New(sha1.New()).Encrypt(bs) }
func Encrypt1Bytes(bs []byte) []byte  { return Encrypt1(bs).Bytes() }
func Encrypt1ASCII(bs []byte) string  { return Encrypt1(bs).ASCII() }
func Encrypt1HEX(bs []byte) string    { return Encrypt1(bs).HEX() }
func Encrypt1Base64(bs []byte) string { return Encrypt1(bs).Base64() }

func Hmac1(data, secret []byte) bytes.Entity    { return crypt.Hmac(sha1.New, secret).Encrypt(data) }
func Hmac1Bytes(data, secret []byte) []byte     { return Hmac1(data, secret).Bytes() }
func Hmac1ASCII(data, secret []byte) string     { return Hmac1(data, secret).ASCII() }
func Hmac1HEX(data, secret []byte) string       { return Hmac1(data, secret).HEX() }
func Hmac1Base64(data, secret []byte) string    { return Hmac1(data, secret).Base64() }
func Hmac1HEXBase64(data, secret []byte) string { return Hmac1(data, secret).HEXBase64() }

//========================================Sha256========================================

func Encrypt256(bs []byte) bytes.Entity { return crypt.New(sha256.New()).Encrypt(bs) }
func Encrypt256Bytes(bs []byte) []byte  { return Encrypt256(bs).Bytes() }
func Encrypt256ASCII(bs []byte) string  { return Encrypt256(bs).ASCII() }
func Encrypt256HEX(bs []byte) string    { return Encrypt256(bs).HEX() }
func Encrypt256Base64(bs []byte) string { return Encrypt256(bs).Base64() }

func Hmac256(data, secret []byte) bytes.Entity    { return crypt.Hmac(sha256.New, secret).Encrypt(data) }
func Hmac256Bytes(data, secret []byte) []byte     { return Hmac256(data, secret).Bytes() }
func Hmac256ASCII(data, secret []byte) string     { return Hmac256(data, secret).ASCII() }
func Hmac256HEX(data, secret []byte) string       { return Hmac256(data, secret).HEX() }
func Hmac256Base64(data, secret []byte) string    { return Hmac256(data, secret).Base64() }
func Hmac256HEXBase64(data, secret []byte) string { return Hmac256(data, secret).HEXBase64() }

//========================================Sha512========================================

func Encrypt512(bs []byte) bytes.Entity { return crypt.New(sha512.New()).Encrypt(bs) }
func Encrypt512Bytes(bs []byte) []byte  { return Encrypt512(bs).Bytes() }
func Encrypt512ASCII(bs []byte) string  { return Encrypt512(bs).ASCII() }
func Encrypt512HEX(bs []byte) string    { return Encrypt512(bs).HEX() }
func Encrypt512Base64(bs []byte) string { return Encrypt512(bs).Base64() }

func Hmac512(data, secret []byte) bytes.Entity    { return crypt.Hmac(sha512.New, secret).Encrypt(data) }
func Hmac512Bytes(data, secret []byte) []byte     { return Hmac512(data, secret).Bytes() }
func Hmac512ASCII(data, secret []byte) string     { return Hmac512(data, secret).ASCII() }
func Hmac512HEX(data, secret []byte) string       { return Hmac512(data, secret).HEX() }
func Hmac512Base64(data, secret []byte) string    { return Hmac512(data, secret).Base64() }
func Hmac512HEXBase64(data, secret []byte) string { return Hmac512(data, secret).HEXBase64() }
