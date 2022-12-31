package chans

import (
	"context"
	"errors"
	"sync"
)

type Err struct {
	mu     sync.Mutex
	err    error
	ctx    context.Context
	cancel context.CancelFunc
}

func NewErr() *Err {
	e := &Err{}
	e.ctx, e.cancel = context.WithCancel(context.Background())
	return e
}

func (this *Err) Done() <-chan struct{} {
	return this.ctx.Done()
}

func (this *Err) Closed() bool {
	return this.err != nil
}

func (this *Err) Err() error {
	return this.err
}

func (this *Err) Close() error {
	this.CloseWithErr(errors.New("主动关闭"))
	return nil
}

func (this *Err) CloseWithErr(err error) {
	if err == nil || this.err != nil {
		return
	}
	this.mu.Lock()
	defer this.mu.Unlock()
	this.err = err
	if this.cancel != nil {
		this.cancel()
	}
}
