package aes

import (
	"crypto/aes"
	"encoding/base64"

	"github.com/injoyai/base/bytes/crypt"
)

// EncryptECB NoPadding
func EncryptECB(str, key string) string {
	block, err := aes.NewCipher(crypt.DealLength([]byte(key), 16))
	if err != nil {
		return ""
	}
	bs := []byte(str)
	if len(bs)%block.BlockSize() > 0 {
		for i := len(bs) % block.BlockSize(); i < block.BlockSize(); i++ {
			bs = append(bs, crypt.Padding)
		}
	}
	var dst []byte
	tmpData := make([]byte, block.BlockSize())
	for index := 0; index < len(bs); index += block.BlockSize() {
		block.Encrypt(tmpData, bs[index:index+block.BlockSize()])
		dst = append(dst, tmpData...)
	}
	return base64.StdEncoding.EncodeToString(dst)
}

// DecryptECB NoPadding
func DecryptECB(str, key string) string {
	bs, _ := base64.StdEncoding.DecodeString(str)
	block, _ := aes.NewCipher(crypt.DealLength([]byte(key), 16))
	if block == nil || len(bs)%block.BlockSize() > 0 {
		return ""
	}
	var dst []byte
	tmpData := make([]byte, block.BlockSize())
	for index := 0; index < len(bs); index += block.BlockSize() {
		block.Decrypt(tmpData, bs[index:index+block.BlockSize()])
		dst = append(dst, tmpData...)
	}
	return string(dst)
}
