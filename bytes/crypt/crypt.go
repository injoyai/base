package crypt

import (
	"crypto/hmac"
	"github.com/injoyai/base/bytes"
	"hash"
)

const (
	Padding = byte(58)
)

func DealLength(bs []byte, length int) []byte {
	beforeLength := len(bs)
	if beforeLength > length {
		bs = bs[:length]
	} else if beforeLength < length {
		for i := 0; i < length-beforeLength; i++ {
			bs = append(bs, Padding)
		}
	}
	return bs
}

type Crypt struct {
	hash.Hash
}

func Hmac(h func() hash.Hash, key []byte) *Crypt {
	return New(hmac.New(h, key))
}

func New(h hash.Hash) *Crypt {
	return &Crypt{Hash: h}
}

func (this *Crypt) Encrypt(bs []byte) bytes.Entity {
	this.Write(bs)
	return this.Sum(nil)
}

func (this *Crypt) EncryptBytes(bs []byte) []byte {
	return this.Encrypt(bs)
}

func (this *Crypt) EncryptASCII(bs []byte) string {
	return this.Encrypt(bs).ASCII()
}

func (this *Crypt) EncryptHEX(bs []byte) string {
	return this.Encrypt(bs).HEX()
}

func (this *Crypt) EncryptBase64(bs []byte) string {
	return this.Encrypt(bs).Base64()
}

func (this *Crypt) EncryptHEXBase64(bs []byte) string {
	return this.Encrypt(bs).HEXBase64()
}
