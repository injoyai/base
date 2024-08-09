package chans

import (
	"sync"
	"time"
)

func NewSubscribe() *Subscribe {
	return &Subscribe{}
}

/*
Subscribe
订阅
*/
type Subscribe struct {
	list []*Safe
	mu   sync.Mutex
	last interface{}
}

func (this *Subscribe) Len() int {
	return len(this.list)
}

func (this *Subscribe) Cap() int {
	return cap(this.list)
}

func (this *Subscribe) Last() interface{} {
	return this.last
}

func (this *Subscribe) Write(p []byte) (int, error) {
	this.Publish(p, 0)
	return len(p), nil
}

func (this *Subscribe) Publish(i interface{}, timeout ...time.Duration) {
	this.last = i
	for _, v := range this.list {
		v.Add(i, timeout...)
	}
}

func (this *Subscribe) Subscribe(cap ...uint) *Safe {
	s := NewSafe(cap...)
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

func (this *Subscribe) Close() error {
	for _, v := range this.list {
		v.C.Close()
	}
	this.list = nil
	return nil
}
