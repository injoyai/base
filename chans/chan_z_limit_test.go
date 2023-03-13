package chans

import (
	"testing"
	"time"
)

func TestNewLimit(t *testing.T) {
	l := NewLimit(1)

	go func() {
		for {
			time.Sleep(time.Second * 3)

			l.Done()

		}
	}()

	for range Traverse(-1) {
		time.Sleep(time.Second)
		t.Log(l.Try())

	}

}
