package chans

import (
	"github.com/injoyai/base/safe"
	"github.com/injoyai/base/types"
	"github.com/injoyai/conv"
)

func NewSafe[T any](cap ...int) *Safe[T] {
	c := make(types.Chan[T], conv.Default(0, cap...))
	return &Safe[T]{
		Chan: c,
		Closer: safe.NewCloser().SetCloseFunc(func(err error) error {
			close(c)
			return nil
		}),
	}
}

type Safe[T any] struct {
	types.Chan[T]
	*safe.Closer
}

func (this *Safe[T]) Close() error {
	return this.Closer.Close()
}
