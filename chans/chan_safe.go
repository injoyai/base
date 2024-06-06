package chans

import (
	"github.com/injoyai/base/safe"
	"github.com/injoyai/conv"
)

func NewSafe(cap ...uint) *Safe {
	c := make(chan interface{}, conv.GetDefaultUint(0, cap...))
	return &Safe{
		C: c,
		Closer: safe.NewCloser().SetCloseFunc(func() error {
			close(c)
			return nil
		}),
	}
}

type Safe struct {
	C
	*safe.Closer
}

func (this *Safe) Close() error {
	return this.Closer.Close()
}
