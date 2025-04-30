package timeout

import (
	"context"
	"github.com/injoyai/base/maps"
	"time"
)

func NewDefault() *Timeout[any] {
	return NewAny()
}

func NewAny() *Timeout[any] {
	return &Timeout[any]{
		interval:    time.Second * 10,
		timeout:     time.Minute,
		timeoutFunc: nil,
		m:           maps.NewSafe[any, time.Time](),
	}
}

func New[K comparable]() *Timeout[K] {
	return &Timeout[K]{
		interval:    time.Second * 10,
		timeout:     time.Minute,
		timeoutFunc: nil,
		m:           maps.NewSafe[K, time.Time](),
	}
}

// Timeout 超时机制
type Timeout[K comparable] struct {
	interval    time.Duration            //检查间隔
	timeout     time.Duration            //超时时间
	timeoutFunc func(key K) error        //超时执行的函数
	m           *maps.Safe[K, time.Time] //
}

func (this *Timeout[K]) SetDealFunc(fn func(key K) error) *Timeout[K] {
	this.timeoutFunc = fn
	return this
}

func (this *Timeout[K]) SetTimeout(timeout time.Duration) *Timeout[K] {
	this.timeout = timeout
	return this
}

func (this *Timeout[K]) SetInterval(interval time.Duration) *Timeout[K] {
	this.interval = interval
	return this
}

func (this *Timeout[K]) Keep(key K) {
	this.m.Set(key, time.Now())
}

func (this *Timeout[K]) Del(key K) {
	this.m.Del(key)
}

func (this *Timeout[K]) Run(ctx context.Context) {
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
