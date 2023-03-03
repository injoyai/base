package crypt

import (
	"crypto/md5"
	"testing"
)

func TestNew(t *testing.T) {
	data := []byte("test")
	if Hmac(md5.New, data).EncryptHEX(data) !=
		New(md5.New()).Hmac(data).EncryptHEX(data) {
		t.Error("函数错误")
		return
	}
}
