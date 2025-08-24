package safe

import (
	"context"
	"github.com/injoyai/conv"
	"sync/atomic"
)

func NewRunner2(fn func(ctx context.Context) error) *Runner2 {
	//去除一开始就返回一个Done,否则可能出现还没Run起来就Done退出的情况
	return &Runner2{
		fn:   fn,
		stop: make(chan struct{}),
	}
}

type Runner2 struct {
	running uint32
	fn      func(ctx context.Context) error
	stop    chan struct{}
	cancel  context.CancelFunc
}

// SetFunc 设置启用的函数
func (this *Runner2) SetFunc(fn func(ctx context.Context) error) *Runner2 {
	this.fn = fn
	return this
}

// Done 结束信号,以Running为准,
// Done的时候可能还在执行,
// Done可能是发了个信号给函数
func (this *Runner2) Done() <-chan struct{} {
	return this.stop
}

// Running 是否在运行
func (this *Runner2) Running() bool {
	return atomic.LoadUint32(&this.running) == 1
}

func (this *Runner2) Run(ctx ...context.Context) (err error) {

	//可选自定义context
	_ctx := conv.Default(nil, ctx...)
	if _ctx == nil {
		_ctx = context.Background()
	}

	//判断是否已经启用
	if atomic.CompareAndSwapUint32(&this.running, 0, 1) {

		//设置未启用状态
		defer atomic.StoreUint32(&this.running, 0)

		if this.fn != nil {
			defer Recover(&err)
			select {
			case <-this.stop:
				//已关闭,则新建
				this.stop = make(chan struct{})
			default:
			}
			defer func() { close(this.stop) }()

			//通过上下文来关闭进程
			ctx2, cancel := context.WithCancel(_ctx)
			defer cancel()
			this.cancel = cancel
			return this.fn(ctx2)
		}
	}

	return
}

// Start 启用,协程执行运行
func (this *Runner2) Start(ctx context.Context) {
	go this.Run(ctx)
}

// Stop 结束,释放结束信号,可选是否等待结束完成
func (this *Runner2) Stop(wait ...bool) {
	if this.cancel != nil {
		this.cancel()
		if this.Running() && len(wait) > 0 && wait[0] {
			//如果需要等待,并且在运行中,则等待结束
			<-this.Done()
		}
	}
}

// Restart 重新执行函数
func (this *Runner2) Restart(ctx context.Context) {
	this.Stop(true)
	this.Start(ctx)
}
