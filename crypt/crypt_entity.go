package crypt

import (
	"crypto/hmac"
	"github.com/injoyai/base/types"
	"hash"
)

type Entity struct {
	hash.Hash
}

// Hmac 哈希加密 same New(h).Hmac(key)
func Hmac(h func() hash.Hash, key []byte) *Entity {
	return New(hmac.New(h, key))
}

// New 加密实例
func New(h hash.Hash) *Entity {
	return &Entity{Hash: h}
}

// Hmac 哈希加密
func (this *Entity) Hmac(key []byte) *Entity {
	this.Hash = hmac.New(func() hash.Hash {
		return this.Hash
	}, key)
	return this
}

// Encrypt 加密,返回bytes.Entity 等价于 []byte
func (this *Entity) Encrypt(bs []byte) types.Bytes {
	this.Write(bs)
	return this.Sum(nil)
}

// EncryptBytes 加密,返回字节 []byte
func (this *Entity) EncryptBytes(bs []byte) []byte {
	return this.Encrypt(bs)
}

// EncryptString 加密,返回字符
func (this *Entity) EncryptString(bs []byte) string {
	return this.Encrypt(bs).String()
}

// EncryptHEX 加密,返回hex,例如md5加密
func (this *Entity) EncryptHEX(bs []byte) string {
	return this.Encrypt(bs).HEX()
}

// EncryptBase64 加密,返回base64,大部分
func (this *Entity) EncryptBase64(bs []byte) string {
	return this.Encrypt(bs).Base64()
}

// EncryptHEXBase64 加密,先用hex编码 >>> ascii解码 >>> 再base64编码
func (this *Entity) EncryptHEXBase64(bs []byte) string {
	return this.Encrypt(bs).HEXBase64()
}
