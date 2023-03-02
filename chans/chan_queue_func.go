package chans

import (
	"github.com/injoyai/conv"
)

type QueueFunc struct{ *Entity }

func NewQueueFunc(num int, cap ...int) *QueueFunc {
	e := NewEntity(num, cap...)
	e.SetHandler(func(no, num int, data interface{}) {
		data.(func(no int, num int))(no, num)
	})
	return &QueueFunc{e}
}

func (this *QueueFunc) Do(fn ...func(no int, num int)) error {
	return this.Entity.Do(conv.Interfaces(fn)...)
}

func (this *QueueFunc) Try(fn ...func(no int, num int)) error {
	return this.Entity.Try(conv.Interfaces(fn)...)
}
