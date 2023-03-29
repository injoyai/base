package chans

import (
	"context"
	"github.com/injoyai/conv"
)

// QueueFunc 协程数量控制
type QueueFunc struct{ *Entity }

func NewQueueFunc(num int, cap ...int) *QueueFunc {
	e := NewEntity(num, cap...)
	e.SetHandler(func(ctx context.Context, no, num int, data interface{}) {
		data.(func(ctx context.Context, no int, num int))(ctx, no, num)
	})
	return &QueueFunc{e}
}

func (this *QueueFunc) Do(fn ...func(ctx context.Context, no int, num int)) error {
	return this.Entity.Do(conv.Interfaces(fn)...)
}

func (this *QueueFunc) Try(fn ...func(ctx context.Context, no int, num int)) error {
	return this.Entity.Try(conv.Interfaces(fn)...)
}
