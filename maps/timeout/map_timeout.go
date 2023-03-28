package timeout

import (
	"context"
	"github.com/injoyai/base/maps"
	"time"
)

func New() *Timeout {
	return NewWithContext(context.Background())
}

func NewWithContext(ctx context.Context) *Timeout {
	ctx, cancel := context.WithCancel(ctx)
	return &Timeout{
		interval:    time.Second * 10,
		timeout:     time.Minute,
		timeoutFunc: nil,
		m:           maps.NewSafe(),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Timeout 超时机制
type Timeout struct {
	interval    time.Duration //检查间隔
	timeout     time.Duration //超时时间
	timeoutFunc func()        //超时执行的函数
	m           *maps.Safe    //
	ctx         context.Context
	cancel      context.CancelFunc
}

func (this *Timeout) SetDealFunc(fn func()) *Timeout {
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

func (this *Timeout) Run() {
	for {
		<-time.After(this.interval)
		select {
		case <-this.ctx.Done():
			return
		default:
			now := time.Now()
			this.m.Range(func(key, value interface{}) bool {
				if now.Sub(value.(time.Time)) > this.timeout {
					if this.timeoutFunc != nil {
						this.timeoutFunc()
						this.m.Del(key)
					}
				}
				return true
			})
		}

	}
}
