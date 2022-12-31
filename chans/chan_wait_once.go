package chans

import "sync"

var waitOncePool = sync.Pool{}

type WaitOnce struct {
	c chan interface{}
}

func (this *WaitOnce) Done(any ...interface{}) {
	defer recover()
	var v interface{}
	if len(any) > 0 {
		v = any[0]
	}
	select {
	case this.c <- v:
	default:
	}
	close(this.c)
}

func (this *WaitOnce) Wait() interface{} {
	v := <-this.c
	defer waitOncePool.Put(this.c)
	return v
}

// NewWaitOnce 只使用一次,
func NewWaitOnce() *WaitOnce {
	if waitOncePool.New == nil {
		waitOncePool.New = func() interface{} {
			return make(chan interface{}, 1)
		}
	}
	return &WaitOnce{c: waitOncePool.Get().(chan interface{})}
}
