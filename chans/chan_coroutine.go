package chans

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Coroutine interface {
	Do(f func())
	DoWait(f func())
	DoRetry(f func() error, retry int)
	DoContext(f func(ctx context.Context) error)

	Wait()
	OnErr(func(err error))
	SetTimeout(timeout time.Duration)
}

func NewCoroutine(max int) Coroutine {
	return &coroutine{
		wg:  sync.WaitGroup{},
		c:   make(chan struct{}, max),
		ctx: context.Background(),
	}
}

type coroutine struct {
	wg     sync.WaitGroup
	c      chan struct{}
	ctx    context.Context
	cancel context.CancelFunc
	onErr  func(err error)
}

func (this *coroutine) OnErr(f func(err error)) {
	this.onErr = f
}

// SetTimeout 设置超时
func (this *coroutine) SetTimeout(timeout time.Duration) {
	if this.cancel != nil {
		this.cancel()
	}
	this.ctx, this.cancel = context.WithTimeout(context.Background(), timeout)
}

// Do 执行
func (this *coroutine) Do(f func()) {
	this.do(func(ctx context.Context) error {
		f()
		return nil
	})
}

// DoWait 执行并等待执行完成
func (this *coroutine) DoWait(f func()) {
	<-this.do(func(ctx context.Context) error {
		f()
		return nil
	})
}

// DoRetry 执行并重试
func (this *coroutine) DoRetry(f func() error, retry int) {
	this.do(func(ctx context.Context) (err error) {
		for i := 0; i <= retry; i++ {

			select {
			case <-ctx.Done():
				err = ctx.Err()
				return
			default:
			}

			if err = f(); err == nil {
				return
			}
		}
		return
	})
}

func (this *coroutine) DoContext(f func(ctx context.Context) error) {
	this.do(f)
}

func (this *coroutine) do(f func(ctx context.Context) error) chan struct{} {

	c := make(chan struct{})

	select {
	case <-this.ctx.Done():
		close(c)
		return c
	default:
	}

	this.wg.Add(1)
	this.c <- struct{}{}
	go func(ctx context.Context) {
		var err error
		defer func() {
			if e := recover(); e != nil {
				switch v := e.(type) {
				case error:
					err = v
				default:
					err = fmt.Errorf("%v", e)
				}
			}
			if this.onErr != nil {
				this.onErr(err)
			}
			<-this.c
			this.wg.Done()
			close(c)
		}()
		err = f(ctx)
	}(this.ctx)

	return c
}

func (this *coroutine) Wait() {
	this.wg.Wait()
}
