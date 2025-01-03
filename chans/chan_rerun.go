package chans

import (
	"context"
	"sync/atomic"
)

type Rerun struct {
	ch      chan context.Context
	fn      func(ctx context.Context)
	enable  *uint32
	running *uint32
	cancel  context.CancelFunc
}

// NewRerun 执行一个函数,重新执行函数,会关掉老的函数
func NewRerun(fn func(ctx context.Context)) *Rerun {
	e := &Rerun{
		ch:      make(chan context.Context),
		fn:      fn,
		enable:  new(uint32),
		running: new(uint32),
	}
	go func() {
		for ctx := range e.ch {
			if e.fn != nil {
				atomic.StoreUint32(e.running, 1)
				e.fn(ctx)
				atomic.StoreUint32(e.running, 0)
			}
		}
	}()
	return e
}

// SetHandler 重新设置处理函数
func (this *Rerun) SetHandler(fn func(ctx context.Context)) {
	this.fn = fn
}

// Running 是否在运行
func (this *Rerun) Running() bool {
	return atomic.LoadUint32(this.running) == 1
}

// Rerun 重新执行函数,关闭正在执行的函数(如果有),执行新的函数
func (this *Rerun) Rerun() {
	if atomic.CompareAndSwapUint32(this.enable, 0, 1) {
		ctx, cancel := context.WithCancel(context.Background())
		if this.cancel != nil {
			this.cancel()
		}
		this.cancel = cancel
		this.ch <- ctx
	}
}

// Enabled 是否启用
func (this *Rerun) Enabled() bool {
	return atomic.LoadUint32(this.enable) == 1
}

// Enable 启用/禁用
func (this *Rerun) Enable(b ...bool) {
	if len(b) > 0 && !b[0] {
		this.Disable()
		return
	}
	this.Rerun()
}

// Disable 禁用
func (this *Rerun) Disable() {
	if atomic.CompareAndSwapUint32(this.enable, 1, 0) {
		if this.cancel != nil {
			this.cancel()
		}
	}
}
