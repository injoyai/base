package chans

import (
	"fmt"
	"github.com/injoyai/conv"
)

type Listen interface {
	Publish(value interface{})
	Subscribe(cap ...int) *Subscribe
}

func NewListen() Listen {
	return &listen{}
}

type listen struct {
	list []interface{}
}

func (this *listen) Publish(value interface{}) {
	for _, v := range this.list {
		v.(*Subscribe).s.Try(value)
	}
}

func (this *listen) Subscribe(cap ...int) *Subscribe {
	c := make(chan interface{}, conv.GetDefaultInt(0, cap...))
	k := fmt.Sprintf("%p", c)
	s := &Subscribe{
		C: c,
		k: k,
		s: NewSafe(c).SetCloseFunc(func() error {
			for i, v := range this.list {
				if v.(*Subscribe).k == k {
					this.list = append(this.list[:i], this.list[i+1:]...)
					break
				}
			}
			return nil
		}),
	}
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
