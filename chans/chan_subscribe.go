package chans

import (
	"sync"
	"time"
)

func NewSubscribe[K comparable, V any]() *Subscribe[K, V] {
	return &Subscribe[K, V]{m: map[K][]*Safe[V]{}}
}

/*
Subscribe
订阅
*/
type Subscribe[K comparable, V any] struct {
	m    map[K][]*Safe[V]
	mu   sync.RWMutex
	last V
}

func (this *Subscribe[K, V]) Last() V {
	return this.last
}

func (this *Subscribe[K, V]) Publish(topic K, value V, timeout ...time.Duration) {
	this.last = value
	this.mu.RLock()
	ls, ok := this.m[topic]
	this.mu.RUnlock()
	if ok {
		for _, v := range ls {
			v.Add(value, timeout...)
		}
	}
}

func (this *Subscribe[K, V]) Subscribe(topic K, cap ...int) *Safe[V] {
	s := NewSafe[V](cap...)
	s.SetCloseFunc(func(err error) error {
		this.mu.RLock()
		ls, ok := this.m[topic]
		this.mu.RUnlock()
		if ok {
			for i, v := range ls {
				if v == s {
					this.mu.Lock()
					this.m[topic] = append(this.m[topic][:i], this.m[topic][i+1:]...)
					this.mu.Unlock()
					break
				}
			}
		}
		return nil
	})
	this.mu.Lock()
	this.m[topic] = append(this.m[topic], s)
	this.mu.Unlock()
	return s
}

func (this *Subscribe[K, V]) Close() error {
	for _, ls := range this.m {
		for _, v := range ls {
			v.Close()
		}
	}
	return nil
}
