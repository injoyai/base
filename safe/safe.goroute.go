package safe

import (
	"context"
	"sync/atomic"
)

func NewGoroute() *Goroute {
	return NewGorouteWithContext(context.Background())
}

func NewGorouteWithContext(ctx context.Context) *Goroute {
	ctx, cancel := context.WithCancel(ctx)
	return &Goroute{
		ctx:    ctx,
		cancel: cancel,
	}
}

type Goroute struct {
	total  uint64
	active int32
	done   chan struct{}
	ctx    context.Context
	cancel context.CancelFunc
}

func (this *Goroute) Go(fs ...func(ctx context.Context)) {
	for _, f := range fs {
		atomic.AddUint64(&this.total, 1)
		if atomic.AddInt32(&this.active, 1) == 1 {
			this.done = make(chan struct{})
		}
		go func(f func(ctx context.Context)) {
			defer func() {
				recover()
				if atomic.AddInt32(&this.active, -1) == 0 {
					close(this.done)
				}
			}()
			f(this.ctx)
		}(f)
	}
}

// Done 等待所有协程执行结束
func (this *Goroute) Done() <-chan struct{} {
	return this.done
}

// Wait 等待所有协程执行结束
func (this *Goroute) Wait() {
	<-this.Done()
}

// Active 执行中的协程数量
func (this *Goroute) Active() int {
	return int(atomic.LoadInt32(&this.active))
}

// Total 总执行协程数量
func (this *Goroute) Total() uint64 {
	return atomic.LoadUint64(&this.total)
}

// Stop 结束所有协程,通过上下文的方式,固需要处理函数中的ctx
func (this *Goroute) Stop(wait ...bool) {
	this.cancel()
	if len(wait) > 0 && wait[0] {
		this.Done()
	}
}
