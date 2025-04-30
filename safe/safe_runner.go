package safe

import (
	"context"
	"sync/atomic"
)

func NewRunner(fn func(ctx context.Context) error) *Runner {
	//去除一开始就返回一个Done,否则可能出现还没Run起来就Done退出的情况
	return &Runner{
		fn:   fn,
		stop: make(chan struct{}),
	}
}

type Runner struct {
	running uint32
	fn      func(ctx context.Context) error
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

func (this *Runner) Run(ctx context.Context) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	//判断是否已经启用
	if atomic.CompareAndSwapUint32(&this.running, 0, 1) {

		return func() (err error) {

			//设置未启用状态
			defer atomic.StoreUint32(&this.running, 0)
			defer Recover(&err)
			select {
			case <-this.stop:
				//已关闭,则新建
				this.stop = make(chan struct{})
			default:
			}
			defer func() { close(this.stop) }()

			//通过上下文来关闭进程
			ctx, cancel := context.WithCancel(ctx)
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
	go this.Run(context.Background())
}

// Stop 结束,释放结束信号,可选是否等待结束完成
func (this *Runner) Stop(wait ...bool) {
	if this.cancel != nil {
		this.cancel()
		if this.Running() && len(wait) > 0 && wait[0] {
			//如果需要等待,并且在运行中,则等待结束
			<-this.Done()
		}
	}
}

// Restart 重新执行函数
func (this *Runner) Restart() {
	this.Stop(true)
	this.Start()
}
