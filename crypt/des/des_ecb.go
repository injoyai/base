package des

import (
	"bytes"
	"crypto/des"
	"encoding/base64"
	"github.com/injoyai/base/crypt"
	"github.com/injoyai/base/types"
)

//========================================EncryptECB========================================

// EncryptECB DES加密,ECB模式,默认ascii>>>base64
func EncryptECB(str, key string) string       { return encryptECB(str, key).Base64() }
func EncryptECBBytes(str, key string) []byte  { return encryptECB(str, key).Bytes() }
func EncryptECBString(str, key string) string { return encryptECB(str, key).String() }
func EncryptECBHEX(str, key string) string    { return encryptECB(str, key).HEX() }
func EncryptECBBase64(str, key string) string { return encryptECB(str, key).Base64() }

//========================================EncryptECB========================================

// DecryptECB DES解密,ECB模式,默认base64>>>ascii
func DecryptECB(str, key string) string       { return decryptECB(str, key).String() }
func DecryptECBBytes(str, key string) []byte  { return decryptECB(str, key).Bytes() }
func DecryptECBString(str, key string) string { return decryptECB(str, key).String() }
func DecryptECBHEX(str, key string) string    { return decryptECB(str, key).HEX() }
func DecryptECBBase64(str, key string) string { return decryptECB(str, key).Base64() }

//========================================inside========================================

// encryptECB DES加密,ECB模式
func encryptECB(str, key string) types.Bytes {
	data := []byte(str)
	keyBs := crypt.DealLength([]byte(key), 8)
	block, err := des.NewCipher(keyBs)
	if err != nil {
		return nil
	}
	bs := block.BlockSize()
	data = func(ciphertext []byte, blockSize int) []byte {
		padding := blockSize - len(ciphertext)%blockSize
		padtext := bytes.Repeat([]byte{byte(padding)}, padding)
		return append(ciphertext, padtext...)
	}(data, bs)
	if len(data)%bs != 0 {
		return nil
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return out
}

// decryptECB DES解密,ECB模式
func decryptECB(str, key string) types.Bytes {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil
	}
	keyBs := crypt.DealLength([]byte(key), 8)
	block, err := des.NewCipher(keyBs)
	if err != nil {
		return nil
	}
	bs := block.BlockSize()
	if len(data)%bs != 0 {
		return nil
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
	return out
}
