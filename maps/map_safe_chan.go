package maps

import (
	"fmt"
	"sync/atomic"
	"time"
)

const (
	chanTryInput       = "try"
	chanMustInput      = "must"
	chanGoMustInput    = "go-must"
	chanTimeoutInput   = "timeout"
	chanGoTimeoutInput = "go-timeout"
)

type Chan struct {
	key          interface{}
	inputType    string
	inputTimeout time.Duration
	C            chan interface{}
	closeFunc    func()
	closed       uint32
}

func (this *Chan) Close() (err error) {
	if atomic.CompareAndSwapUint32(&this.closed, 0, 1) {
		defer func() {
			if e := recover(); e != nil {
				err = fmt.Errorf("%v", e)
			}
		}()
		if this.closeFunc != nil {
			this.closeFunc()
		}
		close(this.C)
	}
	return
}

func (this *Chan) Closed() bool {
	return this.closed == 1
}

func (this *Chan) setCloseFunc(fn func()) *Chan {
	this.closeFunc = fn
	return this
}

func (this *Chan) add(value interface{}) {
	if this.Closed() {
		return
	}
	switch this.inputType {
	case chanTryInput:
		select {
		case this.C <- value:
		default:
		}
	case chanMustInput:
		this.C <- value
	case chanGoMustInput:
		go this.add(value)
	case chanTimeoutInput:
		timer := time.NewTimer(this.inputTimeout)
		defer timer.Stop()
		select {
		case this.C <- value:
		case <-timer.C:
		}
	case chanGoTimeoutInput:
		go this.add(value)
	default:
		this.inputType = chanTryInput
		this.add(value)
	}
}

func newChan(inputType string, key interface{}, cap ...uint) *Chan {
	c := &Chan{inputType: inputType, key: key}
	if len(cap) > 0 && cap[0] > 0 {
		c.C = make(chan interface{}, cap[0])
	} else {
		c.C = make(chan interface{})
	}
	return c
}
