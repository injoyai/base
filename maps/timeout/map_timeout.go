package timeout

import (
	"context"
	"github.com/injoyai/base/maps"
	"time"
)

func New() *Timeout {
	return &Timeout{
		interval:    time.Second * 10,
		timeout:     time.Minute,
		timeoutFunc: nil,
		m:           maps.NewSafe(),
	}
}

// Timeout 超时机制
type Timeout struct {
	interval    time.Duration               //检查间隔
	timeout     time.Duration               //超时时间
	timeoutFunc func(key interface{}) error //超时执行的函数
	m           *maps.Safe                  //
}

func (this *Timeout) SetDealFunc(fn func(key interface{}) error) *Timeout {
	this.timeoutFunc = fn
	return this
}

func (this *Timeout) SetTimeout(timeout time.Duration) *Timeout {
	this.timeout = timeout
	return this
}

func (this *Timeout) SetInterval(interval time.Duration) *Timeout {
	this.interval = interval
	return this
}

func (this *Timeout) Keep(key interface{}) {
	this.m.Set(key, time.Now())
}

func (this *Timeout) Del(key interface{}) {
	this.m.Del(key)
}

func (this *Timeout) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(this.interval):
			now := time.Now()
			this.m.Range(func(key, value interface{}) bool {
				if now.Sub(value.(time.Time)) > this.timeout {
					if this.timeoutFunc == nil && this.timeoutFunc(key) == nil {
						//超时函数未设置,或者执行成功,则删除缓存
						this.m.Del(key)
					}
				}
				return true
			})
		}
	}
}
