package maps

import (
	"strconv"
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
		t.Log(m.c.Load(key))
		c2.Close()
		t.Log(m.c.Load(key))
	}()
	for i := 0; ; i++ {
		time.Sleep(time.Second)
		m.Set(key, i)
	}
}

// 0s
func TestNewSafe2(t *testing.T) {
	c := &Chan{}
	var a interface{}
	for i := 0; i < 10000000; i++ {
		a = &c
	}
	_ = a
}

// 0.59s
func TestNewSafe3(t *testing.T) {
	c := &Chan{}
	_ = c

	for i := 0; i < 10000000; i++ {
		strconv.Itoa(int(time.Now().UnixNano()))
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
