package safe

import (
	"github.com/injoyai/base/types"
	"github.com/injoyai/conv"
	"sync"
)

// NewChan 安全关闭的channel
func NewChan[T any](cap ...int) *Chan[T] {
	_cap := conv.Default(0, cap...)
	return &Chan[T]{
		C:    make(types.Chan[T], _cap),
		done: make(chan struct{}),
	}
}

type Chan[T any] struct {
	C    types.Chan[T]
	once sync.Once
	done chan struct{}
}

func (this *Chan[T]) Close() {
	this.once.Do(func() {
		close(this.C)
		close(this.done)
	})
}

func (this *Chan[T]) Done() <-chan struct{} {
	return this.done
}

func (this *Chan[T]) Closed() bool {
	select {
	case <-this.done:
		return true
	default:
		return false
	}
}
