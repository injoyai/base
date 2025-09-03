package aes

import (
	"crypto/aes"
	"errors"
	"github.com/injoyai/base/crypt"
	"github.com/injoyai/base/types"
)

//========================================EncryptECB========================================

// EncryptECB NoPadding
func EncryptECB(bs, key []byte) (types.Bytes, error) {
	return encryptECB(bs, key)
}

func EncryptECBBytes(bs, key []byte) ([]byte, error) {
	x, err := encryptECB(bs, key)
	return x.Bytes(), err
}

func EncryptECBString(bs, key []byte) (string, error) {
	x, err := encryptECB(bs, key)
	return x.String(), err
}

func EncryptECBHEX(bs, key []byte) (string, error) {
	x, err := encryptECB(bs, key)
	return x.HEX(), err
}

func EncryptECBBase64(bs, key []byte) (string, error) {
	x, err := encryptECB(bs, key)
	return x.Base64(), err
}

//========================================DecryptECB========================================

// DecryptECB NoPadding
func DecryptECB(bs, key []byte) (types.Bytes, error) {
	return decryptECB(bs, key)
}

func DecryptECBBytes(bs, key []byte) ([]byte, error) {
	x, err := decryptECB(bs, key)
	return x.Bytes(), err
}

func DecryptECBString(bs, key []byte) (string, error) {
	x, err := decryptECB(bs, key)
	return x.String(), err
}

func DecryptECBHEX(bs, key []byte) (string, error) {
	x, err := decryptECB(bs, key)
	return x.HEX(), err
}

func DecryptECBBase64(bs, key []byte) (string, error) {
	x, err := decryptECB(bs, key)
	return x.Base64(), err
}

//========================================inside========================================

// encryptECB NoPadding
func encryptECB(bs, key []byte) (types.Bytes, error) {
	block, err := aes.NewCipher(crypt.DealLength(key, 16))
	if err != nil {
		return nil, err
	}
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
	return dst, nil
}

// decryptECB NoPadding
func decryptECB(bs, key []byte) (types.Bytes, error) {
	block, err := aes.NewCipher(crypt.DealLength(key, 16))
	if err != nil {
		return nil, err
	}
	if block == nil || len(bs)%block.BlockSize() > 0 {
		return nil, errors.New("意外的错误")
	}
	var dst []byte
	tmpData := make([]byte, block.BlockSize())
	for index := 0; index < len(bs); index += block.BlockSize() {
		block.Decrypt(tmpData, bs[index:index+block.BlockSize()])
		dst = append(dst, tmpData...)
	}
	return dst, nil
}
