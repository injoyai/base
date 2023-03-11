package safe

import (
	"sync/atomic"
)

type Int64 int64

func NewInt64(n int64) *Int64 {
	return (*Int64)(&n)
}

func (this *Int64) Add(n int64) {
	atomic.AddInt64((*int64)(this), n)
}

func (this *Int64) Swap(new int64) int64 {
	return atomic.SwapInt64((*int64)(this), new)
}

func (this *Int64) CompareAndSwap(old, new int64) bool {
	return atomic.CompareAndSwapInt64((*int64)(this), old, new)
}

type Int32 int32

func NewInt32(n int32) *Int32 {
	return (*Int32)(&n)
}

func (this *Int32) Add(n int32) {
	atomic.AddInt32((*int32)(this), n)
}

func (this *Int32) Swap(new int32) int32 {
	return atomic.SwapInt32((*int32)(this), new)
}

func (this *Int32) CompareAndSwap(old, new int32) bool {
	return atomic.CompareAndSwapInt32((*int32)(this), old, new)
}

type Bool struct {
	n uint32
}

func (this *Bool) IsTrue() bool {
	return atomic.LoadUint32(&this.n) == 1
}

// ListenTrue 设置成true,发生变化并执行
func (this *Bool) ListenTrue(fn func()) {
	if atomic.CompareAndSwapUint32(&this.n, 0, 1) {
		fn()
	}
}

// ListenFalse 设置成false,发生变化并执行
func (this *Bool) ListenFalse(fn func()) {
	if atomic.CompareAndSwapUint32(&this.n, 1, 0) {
		fn()
	}
}

/*
Once 执行一次
区别于源码的sync.Once,源码等待函数执行完成
*/
type Once struct {
	n uint32
}

func (this *Once) Do(fn func()) bool {
	if atomic.CompareAndSwapUint32(&this.n, 0, 1) {
		fn()
		return true
	}
	return false
}
