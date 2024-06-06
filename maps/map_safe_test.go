package maps

import (
	"github.com/injoyai/conv"
	"net/http"
	_ "net/http/pprof"
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
	start := time.Now()
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
	t.Log(time.Since(start))
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

// [5.6s,6.6s,6.18s,5.6s,5.72,5.65s,5.61s,5.61s,5.91s]/千万次 1622MB
func TestNewMapSet(t *testing.T) {
	m := NewSafe(WithBase())
	for i := 0; i < 10000000; i++ {
		m.Set(i, i)
	}
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	t.Logf("使用内存: %d MB", ms.TotalAlloc/1024/1024)
}

// [4.37s,4.47s,4.37s,4.6s,4.57s]/千万次 内存使用量1393MB
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

// [9.95s,11.78s,11.72s,10.58s]/千万次写读
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

// [8.63s,8.8s,8.72s,8.55s,8.76s]/千万次写读 内存1479MB
func TestNewMap2(t *testing.T) {
	m := sync.Map{}
	for i := 0; i < 10000000; i++ {
		m.Store(i, i)
	}
	for i := 0; i < 10000000; i++ {
		_, _ = m.Load(i)
	}
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	t.Logf("使用内存: %d MB", ms.TotalAlloc/1024/1024)
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

func TestSet(t *testing.T) {
	num := 10000000

	last := uint64(0)
	printMem := func(start time.Time) {
		t.Log("耗时: ", time.Since(start))
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		t.Logf("使用内存: %d MB\n", (ms.TotalAlloc-last)/1024/1024)
		last = ms.TotalAlloc
	}

	{ //原生 内存[1393MB] 耗时[4.4s,4.2s,4.1s,4.5s,4.1s,4.5s,4.2s,4.2s,4.2s]
		start := time.Now()
		m := make(map[interface{}]interface{})
		for i := 0; i < num; i++ {
			m[i] = i
		}
		printMem(start)
	}

	{ //Safe-原生 内存[1621MB] 耗时[5.4s,5.8s,5.8s,5.7s,5.6s,5.7s,5.4s,5.9s,5.4s]
		start := time.Now()
		m := NewSafe(WithBase())
		for i := 0; i < num; i++ {
			m.Set(i, i)
		}
		printMem(start)
	}

	{ //Sync 内存[1478MB] 耗时[7.1s,6.9s,7.6s,6.7,7.2s,6.8s]
		start := time.Now()
		m := sync.Map{}
		for i := 0; i < num; i++ {
			m.Store(i, i)
		}
		printMem(start)
	}

	{ //Safe-Sync 内存[1707MB] 耗时[8.4s,8.8s,9.1s,8.5s]
		start := time.Now()
		m := NewSafe()
		for i := 0; i < num; i++ {
			m.Set(i, i)
		}
		printMem(start)
	}

}

func TestSetAndGet(t *testing.T) {
	num := 100000
	multi := 1000

	go http.ListenAndServe(":6060", nil)

	t.Logf("测试读读写少场景: 写(%d)次, 随后读(%d)次\n", num, num*multi)

	last := uint64(0)
	printMem := func(start time.Time) {
		t.Log("耗时: ", time.Since(start))
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		t.Logf("使用内存: %d MB\n", ms.HeapAlloc/1024/1024)
		t.Logf("累计使用内存: %d MB\n", (ms.TotalAlloc-last)/1024/1024)
		last = ms.TotalAlloc
	}

	{ //原生 内存[11MB] 耗时[3.7s,4.2s,3.8s,3.9s]
		start := time.Now()
		m := make(map[interface{}]interface{})
		for i := 0; i < num; i++ {
			m[i] = i
		}
		for i := 0; i < num*multi; i++ {
			_ = m[i]
		}
		printMem(start)
	}

	{ //Safe-原生 内存[774MB] 耗时[6.9s,6.8s,6.5s,7.0s]
		start := time.Now()
		m := NewSafe(WithBase())
		for i := 0; i < num; i++ {
			m.Set(i, i)
		}
		for i := 0; i < num*multi; i++ {
			m.Get(i)
		}
		printMem(start)
	}

	{ //Sync 内存[13MB] 耗时[4.3s,4.1s,4.3s]
		start := time.Now()
		m := sync.Map{}
		for i := 0; i < num; i++ {
			m.Store(i, i)
		}
		for i := 0; i < num*multi; i++ {
			m.Load(i)
		}
		printMem(start)
	}

	{ //test
		start := time.Now()
		//m := struct{ Map }{&_sync{}} //内存总分配776MB
		m := struct{ *_sync }{&_sync{}} //内存总分配13MB
		for i := 0; i < num; i++ {
			m.Set(i, i)
		}
		for i := 0; i < num*multi; i++ {
			m.Get(i)
		}
		printMem(start)
	}

	{ //Safe-Sync 内存[778MB] 耗时[6.4s,6.5s,7.2s]
		start := time.Now()
		m := NewSafe()
		for i := 0; i < num; i++ {
			m.Set(i, i)
		}
		for i := 0; i < num*multi; i++ {
			m.Get(i)
		}
		printMem(start)
	}

	select {}

}
