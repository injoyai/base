package chans

import (
	"sync"
	"time"
)

func NewSubscribe[T any]() *Subscribe[T] {
	return &Subscribe[T]{}
}

/*
Subscribe
订阅
*/
type Subscribe[T any] struct {
	list []*Safe[T]
	mu   sync.Mutex
	last T
}

func (this *Subscribe[T]) Len() int {
	return len(this.list)
}

func (this *Subscribe[T]) Cap() int {
	return cap(this.list)
}

func (this *Subscribe[T]) Last() T {
	return this.last
}

func (this *Subscribe[T]) Publish(i T, timeout ...time.Duration) {
	this.last = i
	for _, v := range this.list {
		v.Add(i, timeout...)
	}
}

func (this *Subscribe[T]) Subscribe(cap ...uint) *Safe[T] {
	s := NewSafe[T](cap...)
	s.SetCloseFunc(func(err error) error {
		for i, v := range this.list {
			if v == s {
				this.mu.Lock()
				this.list = append(this.list[:i], this.list[i+1:]...)
				this.mu.Unlock()
				break
			}
		}
		return nil
	})
	this.mu.Lock()
	this.list = append(this.list, s)
	this.mu.Unlock()
	return s
}

func (this *Subscribe[T]) Close() error {
	for _, v := range this.list {
		v.Close()
	}
	this.list = nil
	return nil
}
