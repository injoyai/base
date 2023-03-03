package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/injoyai/base/bytes"
)

//========================================EncryptCBC========================================

// EncryptCBC aes.cbc加密 ,16位长度
func EncryptCBC(bs, key []byte, ivs ...[]byte) (bytes.Entity, error) {
	return encryptCBC(bs, key, ivs...)
}

func EncryptCBCBytes(bs, key []byte, iv ...[]byte) ([]byte, error) {
	x, err := EncryptCBC(bs, key, iv...)
	return x.Bytes(), err
}

func EncryptCBCASCII(bs, key []byte, iv ...[]byte) (string, error) {
	x, err := EncryptCBC(bs, key, iv...)
	return x.ASCII(), err
}

func EncryptCBCHEX(bs, key []byte, iv ...[]byte) (string, error) {
	x, err := EncryptCBC(bs, key, iv...)
	return x.HEX(), err
}

func EncryptCBCBase64(bs, key []byte, iv ...[]byte) (string, error) {
	x, err := EncryptCBC(bs, key, iv...)
	return x.Base64(), err
}

//========================================DecryptCBC========================================

// DecryptCBC aes.cbc解密 ,16位长度
func DecryptCBC(bs, key []byte, ivs ...[]byte) (bytes.Entity, error) {
	return decryptCBC(bs, key, ivs...)
}

func DecryptCBCBytes(bs, key []byte, ivs ...[]byte) ([]byte, error) {
	x, err := decryptCBC(bs, key, ivs...)
	return x.Bytes(), err
}

func DecryptCBCASCII(bs, key []byte, ivs ...[]byte) (string, error) {
	x, err := decryptCBC(bs, key, ivs...)
	return x.ASCII(), err
}

func DecryptCBCHEX(bs, key []byte, ivs ...[]byte) (string, error) {
	x, err := decryptCBC(bs, key, ivs...)
	return x.HEX(), err
}

func DecryptCBCBase64(bs, key []byte, ivs ...[]byte) (string, error) {
	x, err := decryptCBC(bs, key, ivs...)
	return x.Base64(), err
}

//========================================inside========================================

// encryptCBC 字符串加密 CBC ,16位长度
func encryptCBC(bs, key []byte, ivs ...[]byte) (bytes.Entity, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	var iv []byte
	if len(ivs) == 0 {
		iv = key
	} else {
		iv = ivs[0]
	}
	bs = PKCS7Padding(bs, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv[:blockSize])
	result := make([]byte, len(bs))
	blockMode.CryptBlocks(result, bs)
	return result, nil
}

// decryptCBC 字符串解密 ,16位长度
func decryptCBC(bs, key []byte, ivs ...[]byte) (bytes.Entity, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	var iv []byte
	if len(ivs) == 0 {
		iv = key
	} else {
		iv = ivs[0]
	}
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(bs))
	blockMode.CryptBlocks(origData, bs)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	text := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, text...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	padding := int(origData[length-1])
	return origData[:(length - padding)]
}
