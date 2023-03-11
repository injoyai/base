package safe

import (
	"context"
	"errors"
	"sync/atomic"
)

func NewCloser() *Closer {
	return NewCloserWithContext(context.Background())
}

func NewCloserWithContext(ctx context.Context) *Closer {
	ctx, cancel := context.WithCancel(ctx)
	return &Closer{ctx: ctx, cancel: cancel}
}

type Closer struct {
	closed    uint32
	err       error
	closeFunc func() error
	ctx       context.Context
	cancel    context.CancelFunc
}

// Ctx 上下文
func (this *Closer) Ctx() context.Context {
	return this.ctx
}

// Done ctx.Done
func (this *Closer) Done() <-chan struct{} {
	return this.ctx.Done()
}

// Err 错误信息
func (this *Closer) Err() error {
	return this.err
}

// Closed 是否已关闭
func (this *Closer) Closed() bool {
	select {
	case <-this.Done():
		//确保错误信息closeErr已经赋值,不用this.closed==1
		return true
	default:
		return false
	}
}

// Close 关闭,实现io.Closer接口
func (this *Closer) Close() error {
	return this.CloseWithErr(errors.New("主动关闭"))
}

// CloseWithErr 根据错误关闭
func (this *Closer) CloseWithErr(err error) error {
	if err != nil {
		if atomic.CompareAndSwapUint32(&this.closed, 0, 1) {
			this.err = err
			this.cancel()
			if this.closeFunc != nil {
				return this.closeFunc()
			}
			return nil
		}
	}
	return nil
}

// SetCloseFunc 设置关闭函数
func (this *Closer) SetCloseFunc(fn func() error) {
	this.closeFunc = fn
}
