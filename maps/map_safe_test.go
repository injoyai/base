package maps

import (
	"github.com/injoyai/conv"
	"runtime"
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
			t.Log(<-c.Chan())
		}
	}()
	c2 := m.Chan(key)
	go func() {
		for i := 0; i < 10; i++ {
			t.Log(<-c2.Chan())
		}

		c2.Close()
		t.Log(c2.Closed())
	}()
	for i := 0; ; i++ {
		time.Sleep(time.Second)
		m.Set(key, i)
	}
}

// 协程 2.9s,3.12s,3.11s,2.48s,3.5s,2.56s,3.68s,2.67
func TestNewMap3(t *testing.T) {
	m := NewSafe()
	c := make(chan struct{})
	go func() {
		for i := 0; i < 10000000; i++ {
			m.Set(i, i)
		}
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		t.Logf("使用内存: %d MB", ms.TotalAlloc/1024/1024)
		c <- struct{}{}
	}()
	for i := 0; i < 10000000; i++ {
		m.Get(i)
	}
	<-c
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	t.Logf("使用内存: %d MB", ms.TotalAlloc/1024/1024)
}

func TestNewMapGet(t *testing.T) {
	m := NewSafe()
	for i := 0; i < 10000000; i++ {
		m.Set(i, i)
	}
	start := time.Now()
	for i := 0; i < 10000000; i++ {
		m.Get(conv.String(i))
	}
	t.Log(time.Since(start))
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	t.Logf("使用内存: %d MB", ms.TotalAlloc/1024/1024)
}

func TestNewMapSet(t *testing.T) {
	m := NewSafe(WithBase())
	start := time.Now()
	for i := 0; i < 10000000; i++ {
		m.Set(i, i)
	}
	t.Log("耗时: ", time.Since(start))
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	t.Logf("使用内存: %d MB", ms.TotalAlloc/1024/1024)
}

func TestNewMapSet2(t *testing.T) {
	{
		m := make(map[interface{}]interface{})
		mu := sync.RWMutex{}
		for i := 0; i < 10000000; i++ {
			mu.Lock()
			m[i] = i
			mu.Unlock()
		}
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		t.Logf("使用内存: %d MB", ms.TotalAlloc/1024/1024)

		m = nil
		runtime.GC()

		runtime.ReadMemStats(&ms)
		t.Logf("使用内存: %d MB", ms.TotalAlloc/1024/1024)
	}
}

func TestNewMap8(t *testing.T) {
	m := NewSafe()
	for i := 0; i < 10000000; i++ {
		m.Set(i, i)
	}
	for i := 0; i < 10000000; i++ {
		m.Get(i)
	}
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	t.Logf("使用内存: %d MB", ms.TotalAlloc/1024/1024)
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
			t.Log(<-cc.Chan())
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
