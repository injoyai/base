package main

import (
	"github.com/injoyai/base/maps"
	"log"
	"runtime"
	"time"
)

func main() {
	cache := maps.NewSafe(maps.WithBase())
	start := time.Now()
	for i := 0; i < 10000000; i++ {
		cache.Set(i, i)
	}
	log.Println("耗时: ", time.Since(start))
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	log.Printf("使用内存: %d MB\n", ms.TotalAlloc/1024/1024)
}
