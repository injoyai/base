package chans

import (
	"context"
	"sync/atomic"
)

type Rerun struct {
	enable atomic.Value
	e      *Entity
	fn     func(ctx context.Context)
	cancel context.CancelFunc
}

// NewRerun 执行一个函数,重新执行函数,会关掉老的函数
func NewRerun(fn func(ctx context.Context)) *Rerun {
	e := &Rerun{e: NewEntity(1)}
	e.e.SetHandler(func(ctx context.Context, no, num int, data interface{}) {
		fn(data.(context.Context))
	})
	return e
}

// Rerun 重新执行函数,关闭正在执行的函数(如果有),执行新的函数
func (this *Rerun) Rerun() error {
	ctx, cancel := context.WithCancel(context.Background())
	if this.cancel != nil {
		this.cancel()
	}
	this.cancel = cancel
	this.enable.Store(1)
	return this.e.Do(ctx)
}

// Close 关闭正在执行的函数(如果有)
func (this *Rerun) Close() error {
	this.enable.Store(0)
	if this.cancel != nil {
		this.cancel()
	}
	return nil
}

// RerunOrDisable 例如修改配置的时候,重启或者禁用
func (this *Rerun) RerunOrDisable(enable bool) error {
	if enable {
		return this.Rerun()
	}
	return this.Close()
}

// Enabled 是否启用
func (this *Rerun) Enabled() bool {
	return this.enable.Load() == 1
}

// Enable 启用/禁用
func (this *Rerun) Enable(b ...bool) error {
	if len(b) > 0 && !b[0] {
		return this.Disable()
	}
	return this.Rerun()
}

// Disable 禁用
func (this *Rerun) Disable() error {
	return this.Close()
}
