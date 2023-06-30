package chans

import "context"

type Rerun struct {
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
	return this.e.Do(ctx)
}

// Close 关闭正在执行的函数(如果有)
func (this *Rerun) Close() error {
	if this.cancel != nil {
		this.cancel()
	}
	return nil
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
