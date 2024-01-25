package maps

/*



以下测试结果是通过遍历1千万次设置数据,大致代码如下

func main() {
	//cache := cache.New(5*time.Minute, 10*time.Minute)
	//cache := maps.NewSafe()
	//cache := sync.Map{}
	//cache := make(map[interface{}]interface{})
	cache := &Map{m: make(map[string]interface{})}
	start := time.Now()
	for i := 0; i < 10000000; i++ {
		cache.Set(strconv.Itoa(i), i)
	}
	logs.Debug("耗时: ", time.Since(start))
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	logs.Debugf("使用内存: %d MB", ms.TotalAlloc/1024/1024)
}

type Map struct {
	m  map[string]interface{}
	mu sync.RWMutex
}

func (this *Map) Set(key string, value interface{}) {
	this.mu.Lock()
	this.m[key] = value
	this.mu.Unlock()
}

测试结果(千万次写入):
+-----------------------------------------------------------------------------------------------------------------------------------+
|方式1							|是否安全		|用时		|内存			|记录													|
|map[string]interface{}  		|否  			|4.3s   	|1392MB			|														|
|map[interface{}]interface{}	|否				|4.0s		|1392MB			|														|
|map[string]interface{} 加锁		|是				|4.4s		|1392MB			|														|
|github.com/patrickmn/go-cache  |是				|5.0s		|1760MB			|														|
|sync.Map 						|是 			|6.7s    	|1479MB			|6.8s,7.2s,6.7s,6.7s,6.88s,6.68s,6.95s,6.66s			|
|maps.NewSafe()					|是				|8.5s		|1937MB			|8.66s,8.42s,8.34s,8.38s,8.28s,8.45s,8.31s				|
|maps.NewSafe(WithAny)			|是				|5.8s		|1937MB			|5.88s,5.98s,5.75s,5.58s,5.63s,5.98s,5.69s,5.98			|
+-----------------------------------------------------------------------------------------------------------------------------------+

结论(就写入而言):
读写情况的瓶颈在写入,测试结果是写入的情况
能安全使用的情况,使用原生map和sync.RWMutex的性能是最好最稳定的
使用sync.Map,耗时提升50%
使用maps.NewSafe(),耗时提升100%













*/
