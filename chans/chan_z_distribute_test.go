package chans

import (
	"testing"
	"time"
)

func TestNewDistribute(t *testing.T) {
	s := NewDistribute()
	go func() {
		c := s.ChanOne(10).C
		for {
			t.Log("1: ", <-c)
		}
	}()
	go func() {
		c := s.ChanOne(10).C
		for {
			t.Log("2: ", <-c)
		}
	}()
	go func() {
		c := s.ChanOne(10).C
		for {
			t.Log("3: ", <-c)
		}
	}()
	go func() {
		c := s.ChanOne(10).C
		for {
			t.Log("4: ", <-c)
		}
	}()
	for i := 0; ; i++ {
		<-time.After(time.Second)
		s.Add(i)
	}
}

func TestNewDistributeAll(t *testing.T) {
	s := NewDistribute()
	go func() {
		c := s.ChanAll(10).C
		for {
			t.Log("1: ", <-c)
		}
	}()
	go func() {
		c := s.ChanAll(10).C
		for {
			t.Log("2: ", <-c)
		}
	}()
	go func() {
		c := s.ChanAll(10).C
		for {
			t.Log("3: ", <-c)
		}
	}()
	go func() {
		c := s.ChanAll(10).C
		for {
			t.Log("4: ", <-c)
		}
	}()
	for i := 0; ; i++ {
		<-time.After(time.Second)
		s.Add(i)
	}
}
