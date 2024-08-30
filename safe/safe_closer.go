package safe

import (
	"errors"
	"sync/atomic"
)

func NewCloser() *Closer {
	return &Closer{
		done: make(chan struct{}),
	}
}

type Closer struct {
	closed    uint32                //关闭状态
	err       error                 //错误信息
	closeFunc func(err error) error //关闭执行的函数
	done      chan struct{}         //结束信号
}

// Done 关闭信号
func (this *Closer) Done() <-chan struct{} {
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
			close(this.done)
			if this.closeFunc != nil {
				return this.closeFunc(err)
			}
			return nil
		}
	}
	return nil
}

// SetCloseFunc 设置关闭函数
func (this *Closer) SetCloseFunc(fn func(err error) error) *Closer {
	this.closeFunc = fn
	return this
}
