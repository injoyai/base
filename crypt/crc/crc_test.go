package crc

import (
	"testing"
)

func TestEncrypt16(t *testing.T) {
	t.Log(Encrypt16HEX(testData))
}
