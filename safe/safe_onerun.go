package safe

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

/*
OneRun
最多只有一个在运行,需要等待上一个结束
*/
type OneRun interface {
	Run() error
	Running() bool
	Close() error
	SetHandler(fn func(ctx context.Context) error)
}

func NewOneRun(fn func(ctx context.Context) error) OneRun {
	return &oneRun{
		fn:   fn,
		pool: sync.Pool{New: func() interface{} { return make(chan struct{}) }},
	}
}

type oneRun struct {
	cancel  context.CancelFunc
	running uint32
	mu      sync.Mutex
	fn      func(ctx context.Context) error
	pool    sync.Pool
	done    chan struct{}
}

func (this *oneRun) SetHandler(fn func(ctx context.Context) error) {
	this.fn = fn
}

func (this *oneRun) Run() (err error) {

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

func (this *oneRun) Running() bool {
	return atomic.LoadUint32(&this.running) == 1
}

func (this *oneRun) Close() error {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.cancel != nil {
		this.cancel()
		this.cancel = nil
	}
	return nil
}
