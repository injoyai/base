package maps

import (
	"fmt"
	"sync/atomic"
)

type Chan struct {
	key       interface{}
	C         chan interface{}
	closeFunc func()
	closed    uint32
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

func (this *Chan) SetCloseFunc(fn func()) *Chan {
	this.closeFunc = fn
	return this
}

func (this *Chan) TryAdd(value interface{}) {
	if !this.Closed() {
		select {
		case this.C <- value:
		default:
		}
	}
}

func newChan(key interface{}, cap ...uint) *Chan {
	c := &Chan{key: key}
	if len(cap) > 0 && cap[0] > 0 {
		c.C = make(chan interface{}, cap[0])
	} else {
		c.C = make(chan interface{})
	}
	return c
}
