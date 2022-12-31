package chans

import (
	"errors"
)

type Limit struct {
	running chan byte //正在执行的接口
	limit   int       //最大执行数量
}

func NewLimit(limit int) *Limit {
	return &Limit{
		running: make(chan byte, limit),
		limit:   limit,
	}
}

func (this *Limit) Do() error {
	select {
	case this.running <- byte(0):
	default:
		return errors.New("过于频繁,请稍后再试")
	}
	return nil
}

func (this *Limit) Done() {
	select {
	case <-this.running:
	default:
	}
}
