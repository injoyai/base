package chans

import "context"

type Rerun struct {
	e      *Entity
	fn     func(ctx context.Context)
	cancel context.CancelFunc
}

// NewRerun 执行一个函数,重新执行函数
func NewRerun(fn func(ctx context.Context)) *Rerun {
	e := &Rerun{e: NewEntity(1)}
	e.e.SetHandler(func(ctx context.Context, no, num int, data interface{}) {
		fn(data.(context.Context))
	})
	return e
}

// Rerun 重新执行函数
func (this *Rerun) Rerun() {
	ctx, cancel := context.WithCancel(context.Background())
	if this.cancel != nil {
		this.cancel()
	}
	this.cancel = cancel
	this.e.Do(ctx)
}
