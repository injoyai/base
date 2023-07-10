package chans

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

func NewSafe(c ...chan interface{}) *Safe {
	return NewSafeWithContext(context.Background(), c...)
}

func NewSafeWithContext(ctx context.Context, c ...chan interface{}) *Safe {
	var ch chan interface{}
	if len(c) > 0 && c[0] != nil {
		ch = c[0]
	} else {
		ch = make(chan interface{}, 1)
	}
	return &Safe{
		C:     ch,
		close: nil,
		err:   atomic.Value{},
		once:  sync.Once{},
		ctx:   ctx,
	}
}

type Safe struct {
	C     chan interface{}
	close func() error
	err   atomic.Value
	once  sync.Once
	ctx   context.Context
}

func (this *Safe) SetCloseFunc(fn func() error) *Safe {
	this.close = fn
	return this
}

func (this *Safe) Close() (err error) {
	this.once.Do(func() {
		if this.close != nil {
			err = this.close()
		}
		this.err.Store(errors.New("主动关闭"))
		close(this.C)
	})
	return
}

func (this *Safe) Try(value interface{}) (bool, error) {
	if err := this.err.Load(); err != nil {
		return false, err.(error)
	}
	select {
	case <-this.ctx.Done():
		return false, errors.New("上下文关闭")
	case this.C <- value:
		return true, nil
	default:
		return false, nil
	}
}

func (this *Safe) Must(value interface{}) error {
	if err := this.err.Load(); err != nil {
		return err.(error)
	}
	this.C <- value
	return nil
}

func (this *Safe) Timeout(value interface{}, timeout ...time.Duration) error {
	if err := this.err.Load(); err != nil {
		return err.(error)
	}
	timer := time.NewTimer(timeout[0])
	defer timer.Stop()
	select {
	case <-this.ctx.Done():
		return errors.New("上下文关闭")
	case this.C <- value:
	case <-timer.C:
		return errors.New("超时")
	}
	return nil
}
