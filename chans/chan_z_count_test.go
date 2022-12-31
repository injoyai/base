package chans

import (
	"testing"
	"time"
)

func TestCount(t *testing.T) {
	for a := range Count(10, time.Second) {
		t.Log(a)
	}
}
