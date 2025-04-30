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
	return &Rerun{
		dialed:     false,
		dialErr:    errors.New("未连接"),
		onInterval: func(index int) time.Duration { return time.Second * 10 },
		firstErr:   make(chan error),
		OneRun:     NewOneRun(nil),
	}
}

type Rerun struct {
	dialed     bool
	dialErr    error
	onInterval func(retry int) time.Duration
	onDial     func(index, retry int, err error)
	firstErr   chan error
	OneRun
}

func (this *Rerun) OnDial(fn func(index, retry int, err error)) {
	this.onDial = fn
}

func (this *Rerun) OnInterval(fn func(retry int) time.Duration) {
	this.onInterval = fn
}

func (this *Rerun) DialRun(r Dialer) error {
	go this.Run(r)
	err := <-this.firstErr
	return err
}

func (this *Rerun) Run(r Dialer) error {
	this.OneRun.SetHandler(func(ctx context.Context) error {
		for index := 0; ; index++ {
			select {
			case <-ctx.Done():
				return ctx.Err()

			default:
				//等待连接成功
				for retry := 0; ; {
					select {
					case <-ctx.Done():
						return ctx.Err()
					default:
					}
					err := r.Dial(ctx)
					if index == 0 && retry == 0 {
						select {
						case this.firstErr <- err:
						default:
						}
					}
					if this.onDial != nil {
						this.onDial(index, retry, err)
					}
					if err == nil {
						break
					}
					retry++
					this.dialed, this.dialErr = false, err
					<-time.After(this.onInterval(retry))
				}

				this.dialed, this.dialErr = true, nil
				err := r.Run(ctx)
				this.dialed, this.dialErr = false, err

			}
		}
	})
	return this.OneRun.Run(context.Background())
}

func (this *Rerun) Status() (dialed bool, reason string) {
	dialed = this.dialed
	if this.dialErr != nil {
		reason = this.dialErr.Error()
	}
	return
}
