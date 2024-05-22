package safe

import (
	"context"
	"sync/atomic"
)

func NewRunner(fn func(ctx context.Context) error) *Runner {
	return NewRunnerWithContext(context.Background(), fn)
}

func NewRunnerWithContext(ctx context.Context, fn func(ctx context.Context) error) *Runner {
	c := make(chan struct{})
	close(c)
	return &Runner{
		fn:     fn,
		parent: ctx,
		stop:   c,
	}
}

type Runner struct {
	running uint32
	fn      func(ctx context.Context) error
	parent  context.Context
	cancel  context.CancelFunc
	stop    chan struct{}
}

// SetFunc 设置启用的函数
func (this *Runner) SetFunc(fn func(ctx context.Context) error) *Runner {
	this.fn = fn
	return this
}

// Done 结束信号,以Running为准,
// Done的时候可能还在执行,
// Done可能是发了个信号给函数
func (this *Runner) Done() <-chan struct{} {
	return this.stop
}

// Running 是否在运行
func (this *Runner) Running() bool {
	return atomic.LoadUint32(&this.running) == 1
}

func (this *Runner) Run() error {

	//判断是否已经启用
	if atomic.CompareAndSwapUint32(&this.running, 0, 1) {

		return func() (err error) {

			//设置未启用状态
			defer atomic.StoreUint32(&this.running, 0)
			defer Recover(&err)
			this.stop = make(chan struct{})
			defer func() { close(this.stop) }()

			//通过上下文来关闭进程
			ctx, cancel := context.WithCancel(this.parent)
			this.cancel = cancel
			if this.fn != nil {
				return this.fn(ctx)
			}
			return

		}()
	}
	return nil
}

// Start 启用,协程执行运行
func (this *Runner) Start() {
	go this.Run()
}

// Stop 结束,释放结束信号,可选是否等待结束完成
func (this *Runner) Stop(wait ...bool) {
	if this.cancel != nil {
		this.cancel()
	}
	if len(wait) > 0 && wait[0] {
		//等待结束
		<-this.Done()
	}
}

// Restart 重新执行函数
func (this *Runner) Restart() {
	this.Stop(true)
	this.Start()
}
