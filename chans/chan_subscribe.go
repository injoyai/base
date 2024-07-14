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
}

func (this *Subscribe) Len() int {
	return len(this.list)
}

func (this *Subscribe) Cap() int {
	return cap(this.list)
}

func (this *Subscribe) Publish(i interface{}, timeout ...time.Duration) {
	for _, v := range this.list {
		v.Add(i, timeout...)
	}
}

func (this *Subscribe) Subscribe(cap ...uint) *Safe {
	s := NewSafe(cap...)
	s.SetCloseFunc(func(error) error {
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
