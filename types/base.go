package types

import (
	"context"
	"io"
	"time"
)

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Floater interface {
	~float32 | ~float64
}

type Number interface {
	Integer
	Floater
}

type Ctx interface {
	Done() <-chan struct{}
	Err() error
}

// Err 错误,好处是能定义在const
type Err string

func (this Err) Error() string { return string(this) }

type Debugger bool

func (this *Debugger) Debug(b ...bool) {
	*this = Debugger(len(b) == 0 || b[0])
}

type Sorter interface {
	Len() int
	Swap(i, j int)
}

type Closer interface {
	io.Closer
	Closed() bool
	CloseWithErr(err error) error
}

type Runner interface {
	Run(ctx context.Context) error
}

type Cacher[K comparable, V any] interface {
	Get(key K) (V, error)
	Set(key K, value V, expiration ...time.Duration) error
	Del(key K) error
}

type Mapper[K comparable, V any] interface {
	Get(K) (V, bool)
	Set(K, V)
	Del(K)
	Range(func(K, V) bool)
}
