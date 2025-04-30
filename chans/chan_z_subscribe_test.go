package chans

import (
	"testing"
	"time"
)

func TestNewDistribute(t *testing.T) {
	s := NewSubscribe[any]()
	go func() {
		c := s.Subscribe(10).Chan
		for {
			t.Log("1: ", <-c)
		}
	}()
	go func() {
		c := s.Subscribe(10).Chan
		for {
			t.Log("2: ", <-c)
		}
	}()
	go func() {
		c := s.Subscribe(10).Chan
		for {
			t.Log("3: ", <-c)
		}
	}()
	go func() {
		c := s.Subscribe(10).Chan
		for {
			t.Log("4: ", <-c)
		}
	}()
	for i := 0; ; i++ {
		<-time.After(time.Second)
		s.Publish(i)
	}
}
