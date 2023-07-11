package chans

import (
	"fmt"
	"github.com/injoyai/conv"
	"time"
)

type Listen interface {
	Publish(value interface{})
	Subscribe(cap ...uint) *Subscribe
}

func NewListen() Listen {
	return &listen{}
}

type listen struct {
	list []interface{}
}

func (this *listen) Len() int {
	return len(this.list)
}

func (this *listen) Cap() int {
	return cap(this.list)
}

func (this *listen) Publish(value interface{}) {
	for _, v := range this.list {
		v.(*Subscribe).s.Try(value)
	}
}

func (this *listen) TryPublish(value interface{}) {
	for _, v := range this.list {
		v.(*Subscribe).s.Try(value)
	}
}

func (this *listen) MustPublish(value interface{}) {
	for _, v := range this.list {
		v.(*Subscribe).s.Must(value)
	}
}

func (this *listen) TimeoutPublish(value interface{}, timeout time.Duration) {
	for _, v := range this.list {
		v.(*Subscribe).s.Timeout(value, timeout)
	}
}

func (this *listen) Subscribe(cap ...uint) *Subscribe {
	c := make(chan interface{}, conv.GetDefaultUint(0, cap...))
	k := fmt.Sprintf("%p", c)
	s := &Subscribe{
		C: c,
		k: k,
		s: NewSafe(c),
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
	C chan interface{}
	k string
	s *Safe
}

func (this *Subscribe) Close() error {
	return this.s.Close()
}

func (this *Subscribe) Closed() bool {
	return this.s.Closed()
}
