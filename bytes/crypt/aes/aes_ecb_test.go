package aes

import (
	"testing"
)

func TestEncryptECB(t *testing.T) {
	t.Log(EncryptECBBase64([]byte("test"), []byte("")))
	t.Log(EncryptECBBase64([]byte("1234567890123456"), []byte("1234567890123456")))
	t.Log(EncryptECBBase64([]byte("test"), []byte("jdgobvhfobhofhob")))
	//t.Log(DecryptECB("dXzNDNxckOrb7uz2ON0AAA==", "1234567890123456"))
	//t.Log(DecryptECB("ssss==", "1234567890123456"))
}
