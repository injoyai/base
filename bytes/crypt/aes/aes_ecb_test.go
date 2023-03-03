package aes

import (
	"testing"
)

func TestEncryptECB(t *testing.T) {
	t.Log(EncryptECB("test", ""))
	t.Log(EncryptECB("1234567890123456", "1234567890123456"))
	t.Log(EncryptECB("test", "jdgobvhfobhofhob"))
	t.Log(DecryptECB("dXzNDNxckOrb7uz2ON0AAA==", "1234567890123456"))
	t.Log(DecryptECB("ssss==", "1234567890123456"))
}
