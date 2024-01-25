package maps

import "sync"

type Map interface {
	Get(key interface{}) (interface{}, bool)
	Set(key interface{}, value interface{})
	Del(key interface{})
	Range(fn func(key, value interface{}) bool)
}

func WithBase() *_base { return &_base{m: make(map[interface{}]interface{})} }

func WithSync() *_sync { return &_sync{} }

/*



 */

type _base struct {
	m  map[interface{}]interface{}
	mu sync.RWMutex
}

func (this *_base) Set(key interface{}, value interface{}) {
	this.mu.Lock()
	this.m[key] = value
	this.mu.Unlock()
}

func (this *_base) Get(key interface{}) (interface{}, bool) {
	this.mu.RLock()
	value, ok := this.m[key]
	this.mu.RUnlock()
	return value, ok
}

func (this *_base) Del(key interface{}) {
	this.mu.Lock()
	delete(this.m, key)
	this.mu.Unlock()
}

func (this *_base) Range(fn func(key, value interface{}) bool) {
	this.mu.RLock()
	for k, v := range this.m {
		if !fn(k, v) {
			break
		}
	}
	this.mu.RUnlock()
}

type _sync struct {
	sync.Map
}

func (this *_sync) Set(key, value interface{}) {
	this.Map.Store(key, value)
}

func (this *_sync) Get(key interface{}) (interface{}, bool) {
	return this.Map.Load(key)
}

func (this *_sync) Del(key interface{}) {
	this.Map.Delete(key)
}

func (this *_sync) Range(fn func(key, value interface{}) bool) {
	this.Map.Range(func(key, value interface{}) bool {
		return fn(key, value)
	})
}
