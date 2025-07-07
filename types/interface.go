package types

import (
	"context"
	"io"
	"time"
)

// Comparable 可比较,同1.21版本的cmp.Ordered
type Comparable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		~string
}

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Floater interface {
	~float32 | ~float64
}

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

type Signaler interface {
	Done() <-chan struct{}
	Err() error
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
	Running() bool
}

type Cacher[K comparable, V any] interface {
	// Get retrieves a value by key. Returns error if not found or expired.
	Get(key K) (V, error)

	// Set stores a value with an optional expiration.
	Set(key K, value V, expiration ...time.Duration) error

	// Del deletes a value by key.
	Del(key K) error
}

type Mapper[K comparable, V any] interface {
	// Get retrieves a value for a key. Returns false if the key is not present.
	Get(key K) (V, bool)

	// Set stores or overwrites the value for a key.
	Set(key K, value V)

	// Del deletes the value associated with a key.
	Del(key K)

	// Range iterates over all key-value pairs.
	// If the function returns false, iteration stops.
	Range(func(K, V) bool)
}
