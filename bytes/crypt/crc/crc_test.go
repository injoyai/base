package crc

import (
	"log"
	"testing"
)

func TestEncrypt16(t *testing.T) {
	log.Println(Encrypt16HEX(testData))
}
