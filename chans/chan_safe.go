package chans

import (
	"context"
	"errors"
	"github.com/injoyai/base/safe"
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
		C:      ch,
		ctx:    ctx,
		Closer: safe.NewCloser(),
	}
}

type Safe struct {
	C   chan interface{}
	ctx context.Context
	*safe.Closer
}

func (this *Safe) Try(value interface{}) (bool, error) {
	if err := this.Err(); err != nil {
		return false, err
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
	if err := this.Err(); err != nil {
		return err
	}
	this.C <- value
	return nil
}

func (this *Safe) Timeout(value interface{}, timeout time.Duration) error {
	if err := this.Err(); err != nil {
		return err
	}
	timer := time.NewTimer(timeout)
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
