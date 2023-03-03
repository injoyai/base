package md5

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	t.Log(Encrypt("test"))
	t.Log(EncryptHEX("test"))
}
