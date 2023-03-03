package md5

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	if Encrypt("000000") != "670b14728ad9902aecba32e22fa4f6bd" {
		t.Log(Encrypt("000000"))
		t.Error("加密错误")
		return
	}
	if Encrypt("test") != "098f6bcd4621d373cade4e832627b4f6" {
		t.Log(Encrypt("test"))
		t.Error("加密错误")
		return
	}
}
