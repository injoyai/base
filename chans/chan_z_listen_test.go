package chans

import (
	"testing"
	"time"
)

func TestNewListen(t *testing.T) {
	l := NewListen()
	for i := 0; i < 3; i++ {
		go func() {
			c := l.Subscribe()
			for {
				t.Log(<-c.C)
			}
		}()
	}
	go func() {
		c := l.Subscribe()
		defer c.Close()
		for i := 0; i < 3; i++ {
			t.Log(<-c.C)
		}
	}()
	for i := 0; ; i++ {
		<-time.After(time.Second)
		l.Publish(i)
	}
}
