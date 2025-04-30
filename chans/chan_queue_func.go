package chans

import (
	"context"
)

type QueueFuncHandler func(ctx context.Context, no int, num int)

// QueueFunc 协程数量控制
type QueueFunc struct {
	e *Entity[QueueFuncHandler]
}

func NewQueueFunc(num int, cap ...int) *QueueFunc {
	e := NewEntity[QueueFuncHandler](num, cap...)
	e.SetHandler(func(ctx context.Context, no, num int, fn QueueFuncHandler) {
		fn(ctx, no, num)
	})
	return &QueueFunc{e}
}

func (this *QueueFunc) Do(fn ...QueueFuncHandler) error {
	return this.e.Do(fn...)
}

func (this *QueueFunc) Try(fn ...QueueFuncHandler) (bool, error) {
	return this.e.Try(fn...)
}
