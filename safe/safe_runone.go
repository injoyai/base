package safe

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

/*
RunOne
最多只有一个在运行,需要等待上一个结束
*/
type RunOne interface {
	Run() error
	Running() bool
	Close() error
	SetHandler(fn func(ctx context.Context) error)
}

func NewRunOne(fn func(ctx context.Context) error) RunOne {
	return &runOne{
		fn:   fn,
		pool: sync.Pool{New: func() interface{} { return make(chan struct{}) }},
	}
}

type runOne struct {
	cancel  context.CancelFunc
	running uint32
	mu      sync.Mutex
	fn      func(ctx context.Context) error
	pool    sync.Pool
	done    chan struct{}
}

func (this *runOne) SetHandler(fn func(ctx context.Context) error) {
	this.fn = fn
}

func (this *runOne) Run() (err error) {

	if this.fn == nil {
		return errors.New("未设置函数")
	}

	this.mu.Lock()
	if this.cancel != nil {
		this.cancel()
	}
	if this.done != nil {
		//等待上次结束
		<-this.done
	}
	this.done = this.pool.Get().(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	this.cancel = cancel
	this.mu.Unlock()

	defer func() {
		atomic.StoreUint32(&this.running, 0)
		close(this.done)
	}()

	atomic.StoreUint32(&this.running, 1)
	return this.fn(ctx)

}

func (this *runOne) Running() bool {
	return atomic.LoadUint32(&this.running) == 1
}

func (this *runOne) Close() error {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.cancel != nil {
		this.cancel()
		this.cancel = nil
	}
	return nil
}
