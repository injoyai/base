package chans

import (
	"time"
)

type C = Chan

type Chan chan interface{}

func (this Chan) Close() error {
	defer func() { recover() }()
	close(this)
	return nil
}

func (this Chan) Chan() chan interface{} {
	return this
}

func (this Chan) Must(i interface{}) {
	_ = this.Add(i)
}

func (this Chan) Try(i interface{}) bool {
	return this.Add(i, 0)
}

func (this Chan) Timeout(i interface{}, timeout time.Duration) bool {
	return this.Add(i, timeout)
}

func (this Chan) Add(i interface{}, timeout ...time.Duration) bool {
	t := time.Duration(-1)
	if len(timeout) > 0 {
		t = timeout[0]
	}

	switch {
	case t < 0:
		this <- i
		return true

	case t == 0:
		select {
		case this <- i:
			return true
		default:
			return false
		}

	default:
		select {
		case this <- i:
			return true
		case <-time.After(t):
			return false
		}

	}
}
