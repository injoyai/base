package timeout

import (
	"context"
	"github.com/injoyai/base/maps"
	"time"
)

func New() *Timeout {
	return NewGeneric[any]()
}

func NewGeneric[K comparable]() *Generic[K] {
	return &Generic[K]{
		interval:    time.Second * 10,
		timeout:     time.Minute,
		timeoutFunc: nil,
		m:           maps.NewGeneric[K, time.Time](),
	}
}

type Timeout = Generic[any]

// Generic 超时机制
type Generic[K comparable] struct {
	interval    time.Duration               //检查间隔
	timeout     time.Duration               //超时时间
	timeoutFunc func(key K) error           //超时执行的函数
	m           *maps.Generic[K, time.Time] //
}

func (this *Generic[K]) SetDealFunc(fn func(key K) error) *Generic[K] {
	this.timeoutFunc = fn
	return this
}

func (this *Generic[K]) SetTimeout(timeout time.Duration) *Generic[K] {
	this.timeout = timeout
	return this
}

func (this *Generic[K]) SetInterval(interval time.Duration) *Generic[K] {
	this.interval = interval
	return this
}

func (this *Generic[K]) Keep(key K) {
	this.m.Set(key, time.Now())
}

func (this *Generic[K]) Del(key K) {
	this.m.Del(key)
}

func (this *Generic[K]) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(this.interval):
			now := time.Now()
			this.m.Range(func(key K, value time.Time) bool {
				if now.Sub(value) > this.timeout {
					if this.timeoutFunc != nil && this.timeoutFunc(key) == nil {
						//超时函数未设置,或者执行成功,则删除缓存
						this.m.Del(key)
					}
				}
				return true
			})
		}
	}
}
