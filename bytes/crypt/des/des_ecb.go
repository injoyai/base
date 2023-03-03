package des

import (
	"bytes"
	"crypto/des"
	"encoding/base64"

	"github.com/injoyai/base/bytes/crypt"
)

// EncryptECB DES加密,ECB模式
func EncryptECB(str, key string) string {
	data := []byte(str)
	keyBs := crypt.DealLength([]byte(key), 8)
	block, err := des.NewCipher(keyBs)
	if err != nil {
		return ""
	}
	bs := block.BlockSize()
	data = func(ciphertext []byte, blockSize int) []byte {
		padding := blockSize - len(ciphertext)%blockSize
		padtext := bytes.Repeat([]byte{byte(padding)}, padding)
		return append(ciphertext, padtext...)
	}(data, bs)
	if len(data)%bs != 0 {
		return ""
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return base64.StdEncoding.EncodeToString(out)
}

// DecryptECB DES解密,ECB模式
func DecryptECB(str, key string) string {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	keyBs := crypt.DealLength([]byte(key), 8)
	block, err := des.NewCipher(keyBs)
	if err != nil {
		return ""
	}
	bs := block.BlockSize()
	if len(data)%bs != 0 {
		return ""
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Decrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	out = func(origData []byte) []byte {
		length := len(origData)
		unpadding := int(origData[length-1])
		return origData[:(length - unpadding)]
	}(out)
	return string(out)
}
