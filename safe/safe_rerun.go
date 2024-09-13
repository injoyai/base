package safe

import (
	"context"
	"errors"
	"io"
	"time"
)

type Dialer interface {
	Dial(ctx context.Context) error //建立连接,例通过tcp连接到服务器
	Run(ctx context.Context) error  //执行业务逻辑,例监听tcp服务的数据
	io.Closer                       //关闭此次连接,这个可以不用,为了兼容多种情况
}

func NewRerun() *Rerun {
	x := &Rerun{
		Interval: func(index int) time.Duration {
			if index == 0 {
				return time.Second
			}
			return time.Second * 10
		},
		Runner: NewRunner(nil),
		Closer: NewCloserErr(errors.New("未连接")),
	}
	x.Closer.SetCloseFunc(func(err error) error {
		x.Stop()
		return nil
	})
	return x
}

type Rerun struct {
	Dialer
	*Runner
	*Closer //单次生命周期状态,可复用

	//重试间隔函数,参数已经重试次数,为0是刚断开上次连接
	Interval func(retryNum int) time.Duration
}

func (this *Rerun) GetOnline() (online bool, reason string) {
	if this.Closer == nil {
		return false, "未连接"
	}
	if this.Closer.Closed() {
		return false, this.Closer.Err().Error()
	}
	return true, "连接成功"
}

func (this *Rerun) Stop(wait ...bool) {
	if this.Dialer != nil {
		this.Dialer.Close()
	}
	this.Runner.Stop(wait...)
}

func (this *Rerun) Restart() {
	this.Stop(true)
	this.Runner.SetFunc(this.runAfter(0))
	this.Runner.Start()
}

// Update 更新后会尝试连接(等待结果),
// 适用于手动修改配置,实时反馈配置执行结果
// 连接失败会返回错误,并在后台开始尝试,例如服务暂时不行,后续会正常
func (this *Rerun) Update(r Dialer) error {
	//关闭老的,如果存在的话
	this.Stop(true)
	this.Dialer = r
	//连接失败则退出,连接成功则循环执行(断线重连)
	this.Runner.SetFunc(func(ctx context.Context) error { return r.Dial(ctx) })
	err := this.Runner.Run()

	go func(err error) {
		if err == nil {
			//如果连接成功了,则等待此次运行结束
			this.Runner.SetFunc(this.runBefore(r.Run))
		} else {
			//连接失败了,则10秒之后开始重试
			this.Runner.SetFunc(this.runAfter(time.Second * 10))
		}
		this.Runner.Run()
	}(err)

	return err
}

// Must 一直重连直到成功
func (this *Rerun) Must(r Dialer) error {
	//关闭老的,如果存在的话
	this.Stop(true)
	this.Dialer = r
	this.Runner.SetFunc(this.runAfter(0))
	return this.Runner.Run()
}

func (this *Rerun) runBefore(run func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		if run != nil {
			this.Closer.Reset()
			err := run(ctx)
			this.Closer.CloseWithErr(err)
		}

		t := time.NewTimer(0)
		for i := 1; ; i++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-t.C:
				if err := this.Dialer.Dial(ctx); err == nil {
					this.Closer.Reset()
					err = this.Dialer.Run(ctx)
					this.Closer.CloseWithErr(err)
					i = 0
				}
				t.Reset(this.Interval(i))
			}
		}
	}
}

func (this *Rerun) runAfter(after time.Duration) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		t := time.NewTimer(after)
		for i := 1; ; i++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-t.C:
				if err := this.Dialer.Dial(ctx); err == nil {
					this.Closer.Reset()
					err = this.Dialer.Run(ctx)
					this.Closer.CloseWithErr(err)
					i = 0
				}
				t.Reset(this.Interval(i))
			}
		}
	}
}
