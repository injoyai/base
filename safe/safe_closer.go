package safe

import (
	"errors"
	"sync/atomic"
)

func NewCloser() *Closer {
	return &Closer{}
}

type Closer struct {
	closed    uint32
	err       error
	closeFunc func() error
	done      chan struct{}
}

// Done 关闭信号
func (this *Closer) Done() <-chan struct{} {
	if this.done == nil {
		this.done = make(chan struct{})
	}
	return this.done
}

// Err 错误信息
func (this *Closer) Err() error {
	return this.err
}

// Closed 是否已关闭
func (this *Closer) Closed() bool {
	if this == nil {
		//方便业务逻辑 xxx==nil || xxx.Closed()
		return true
	}
	return this.Err() != nil
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
			if this.done != nil {
				close(this.done)
			}
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
