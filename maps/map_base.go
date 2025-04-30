package maps

import "sync"

func WithMutex[K comparable, V any]() *_base[K, V] { return &_base[K, V]{m: make(map[K]V)} }

/*



 */

type _base[K comparable, V any] struct {
	m  map[K]V
	mu sync.RWMutex
}

func (this *_base[K, V]) Set(key K, value V) {
	this.mu.Lock()
	this.m[key] = value
	this.mu.Unlock()
}

func (this *_base[K, V]) Get(key K) (V, bool) {
	this.mu.RLock()
	value, ok := this.m[key]
	this.mu.RUnlock()
	return value, ok
}

func (this *_base[K, V]) Del(key K) {
	this.mu.Lock()
	delete(this.m, key)
	this.mu.Unlock()
}

func (this *_base[K, V]) Range(fn func(key K, value V) bool) {
	this.mu.RLock()
	for k, v := range this.m {
		if !fn(k, v) {
			break
		}
	}
	this.mu.RUnlock()
}
