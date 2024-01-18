package chans

import (
	"fmt"
	"time"
)

func NewListen() *Listen {
	return &Listen{}
}

type Listen struct {
	list []interface{}
}

func (this *Listen) Len() int {
	return len(this.list)
}

func (this *Listen) Cap() int {
	return cap(this.list)
}

func (this *Listen) Publish(value interface{}) {
	for _, v := range this.list {
		v.(*Subscribe).s.Try(value)
	}
}

func (this *Listen) TryPublish(value interface{}) {
	for _, v := range this.list {
		v.(*Subscribe).s.Try(value)
	}
}

func (this *Listen) MustPublish(value interface{}) {
	for _, v := range this.list {
		v.(*Subscribe).s.Must(value)
	}
}

func (this *Listen) TimeoutPublish(value interface{}, timeout time.Duration) {
	for _, v := range this.list {
		v.(*Subscribe).s.Timeout(value, timeout)
	}
}

func (this *Listen) Subscribe(cap ...uint) *Subscribe {
	c := NewSafe(cap...)
	k := fmt.Sprintf("%p", c)
	s := &Subscribe{
		k: k,
		s: c,
	}
	s.s.SetCloseFunc(func() error {
		for i, v := range this.list {
			if v.(*Subscribe).k == k {
				this.list = append(this.list[:i], this.list[i+1:]...)
				break
			}
		}
		return nil
	})
	this.list = append(this.list, s)
	return s
}

// Subscribe 订阅对象,开放指定接口
type Subscribe struct {
	k string
	s *Safe
}

func (this *Subscribe) Chan() chan interface{} {
	return this.s.C
}

func (this *Subscribe) Close() error {
	return this.s.Close()
}

func (this *Subscribe) Closed() bool {
	return this.s.Closed()
}
