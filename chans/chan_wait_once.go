package chans

import (
	"github.com/injoyai/base/types"
	"sync"
)

var waitOncePool = sync.Pool{}

type WaitOnce[T any] types.Chan[T]

func (c WaitOnce[T]) Done(any ...T) {
	defer func() { recover() }()
	var v T
	if len(any) > 0 {
		v = any[0]
	}
	select {
	case c <- v:
	default:
	}
	close(c)
}

func (c WaitOnce[T]) Wait() T {
	v := <-c
	defer waitOncePool.Put(c)
	return v
}

// NewWaitOnce 只使用一次,
func NewWaitOnce[T any]() WaitOnce[T] {
	if waitOncePool.New == nil {
		waitOncePool.New = func() any {
			return make(WaitOnce[T], 1)
		}
	}
	return waitOncePool.Get().(WaitOnce[T])
}
