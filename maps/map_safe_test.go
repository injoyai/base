package maps

import (
	"github.com/injoyai/conv"
	"sync"
	"testing"
	"time"
)

func TestNewSafe(t *testing.T) {
	key := "xxx"
	m := NewSafe()
	m.Set(key, -1)
	c := m.Chan(key, 10)
	m.Set(key, -2)
	go func() {
		for {
			t.Log(<-c.C)
		}
	}()
	c2 := m.Chan(key)
	go func() {
		for i := 0; i < 10; i++ {
			t.Log(<-c2.C)
		}

		c2.Close()
		t.Log(c2.Closed())
	}()
	for i := 0; ; i++ {
		time.Sleep(time.Second)
		m.Set(key, i)
	}
}

// 4.4s,6.41s,4.77s
func TestNewMap3(t *testing.T) {
	m := NewSafe()
	go func() {
		for i := 0; i < 10000000; i++ {
			m.Set(i, i)
		}
	}()
	for i := 0; i < 10000000; i++ {
		m.Get(i)
	}
}

// 3.85s,4.67s,4.07
func TestNewMap2(t *testing.T) {
	m := sync.Map{}
	go func() {
		for i := 0; i < 10000000; i++ {
			m.Store(i, i)
		}
	}()
	for i := 0; i < 10000000; i++ {
		val, _ := m.Load(i)
		if val != nil {
			_ = val.(int)
		}
	}
}

func TestNewMap4(t *testing.T) {
	m := NewSafe()
	go func() {
		cc := m.Chan(1)
		for {
			t.Log(<-cc.C)
		}
	}()
	m.Set(1, 2, time.Second)
	t.Log(m.Get(1))
	t.Log(m.Has(1))
	time.Sleep(time.Second)
	t.Log(m.Get(1))
	t.Log(m.Has(1))
}

func TestSafe_GetOrSetByHandler(t *testing.T) {
	m := NewSafe()
	for i := 0; i < 1000; i++ {
		go func(i int) {
			m.GetOrSetByHandler("", func() (interface{}, error) {
				m.Set(conv.String(i), "")
				return i, nil
			})
		}(i)
	}
	<-time.After(time.Second)
	m.Range(func(key, value interface{}) bool {
		t.Log(key)
		return true
	})
}
