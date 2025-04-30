package types

import (
	"github.com/injoyai/conv"
	"time"
)

type (
	ChanAny    = Chan[any]
	ChanBytes  = Chan[[]byte]
	ChanStruct = Chan[struct{}]
)

type Chan[T any] chan T

func (this Chan[T]) Close() error {
	defer func() { recover() }()
	close(this)
	return nil
}

func (this Chan[T]) Must(i T) {
	this <- i
}

func (this Chan[T]) Try(i T) bool {
	select {
	case this <- i:
		return true
	default:
		return false
	}
}

func (this Chan[T]) Timeout(i T, timeout time.Duration) bool {
	select {
	case this <- i:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (this Chan[T]) Add(i T, timeout ...time.Duration) bool {

	t := conv.Default[time.Duration](-1, timeout...)

	switch {
	case t < 0:
		this.Must(i)
		return true

	case t == 0:
		return this.Try(i)

	default:
		return this.Timeout(i, t)

	}
}
