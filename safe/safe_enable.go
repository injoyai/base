package safe

import (
	"context"
	"sync/atomic"
)

type Enable struct {
	enable uint32
	fn     func(ctx context.Context) error
	cancel context.CancelFunc
}

func NewEnable(fn func(ctx context.Context) error) *Enable {
	return &Enable{
		enable: 0,
		fn:     fn,
		cancel: nil,
	}
}

func (this *Enable) Enabled() bool {
	return atomic.LoadUint32(&this.enable) == 1
}

func (this *Enable) SetFunc(fn func(ctx context.Context) error) *Enable {
	this.fn = fn
	return this
}

func (this *Enable) Enable(b ...bool) error {

	if len(b) > 0 && !b[0] {
		//设置禁用
		return this.Disable()
	}

	//判断是否已经启用
	if atomic.CompareAndSwapUint32(&this.enable, 0, 1) {
		ctx, cancel := context.WithCancel(context.Background())
		if this.cancel != nil {
			this.cancel()
		}
		this.cancel = cancel
		if this.fn != nil {
			err := this.fn(ctx)
			if err != nil {
				//启用失败,设置未启用
				atomic.StoreUint32(&this.enable, 0)
			}
			return err
		}
	}
	return nil
}

func (this *Enable) Disable() error {
	//设置未启用,并关闭上下文
	atomic.StoreUint32(&this.enable, 0)
	if this.cancel != nil {
		this.cancel()
	}
	return nil
}

// Restart 重新执行函数
func (this *Enable) Restart() error {
	if err := this.Disable(); err != nil {
		return err
	}
	return this.Enable()
}
