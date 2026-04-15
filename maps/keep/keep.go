package keep

import (
	"sync"
	"time"

	"github.com/injoyai/conv"
)

type Keeper[K any] interface {
	Keep(key K)
	Keeping(key K, n ...int) bool
	SetTimeout(timeout time.Duration)
	OnKeep(f func(K, time.Time))
}

var _ Keeper[int] = (*Keep[int])(nil)

func New[K comparable](timeout time.Duration) Keeper[K] {
	return &Keep[K]{m: map[K]*info{}, timeout: timeout}
}

type Keep[K comparable] struct {
	m       map[K]*info
	mu      sync.RWMutex
	timeout time.Duration
	onKeep  func(K, time.Time)
}

func (this *Keep[K]) SetTimeout(timeout time.Duration) {
	this.timeout = timeout
}

func (this *Keep[K]) OnKeep(f func(K, time.Time)) {
	this.onKeep = f
}

func (this *Keep[K]) Keeping(key K, n ...int) bool {
	this.mu.RLock()
	defer this.mu.RUnlock()
	_info, ok := this.m[key]
	if !ok || _info == nil {
		return false
	}
	return _info.keeping(this.timeout, n...)
}

func (this *Keep[T]) Keep(key T) {
	this.mu.Lock()
	_info, ok := this.m[key]
	if !ok {
		_info = &info{}
		this.m[key] = _info
	}
	this.mu.Unlock()

	now := time.Now()
	_info.keep(now)
	if this.onKeep != nil {
		this.onKeep(key, now)
	}
}

const _len = 6

type info struct {
	nodes [_len]time.Time
}

func (this *info) keep(t time.Time) {
	copy(this.nodes[:], this.nodes[1:])
	this.nodes[_len-1] = t
}

func (this *info) keeping(timeout time.Duration, n ...int) bool {
	_n := conv.Default(1, n...)
	for i := 0; i < _n && i < _len; i++ {
		if i == 0 {
			if time.Since(this.nodes[_len-1]) > timeout {
				return false
			}
		} else {
			if this.nodes[_len-i].Sub(this.nodes[_len-i-1]) > timeout {
				return false
			}
		}
	}
	return true
}
