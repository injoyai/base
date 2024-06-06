package chans

import (
	"time"
)

type Chan chan interface{}

func (this Chan) Must(i interface{}) { this.Add(i) }

func (this Chan) Try(i interface{}) { this.Add(i, 0) }

func (this Chan) Timeout(i interface{}, timeout time.Duration) { this.Add(i, timeout) }

func (this Chan) Add(i interface{}, timeout ...time.Duration) {
	t := time.Duration(-1)
	if len(timeout) > 0 {
		t = timeout[0]
	}
	switch {
	case t < 0:
		this <- i
	case t == 0:
		select {
		case this <- i:
		default:
		}
	case t > 0:
		select {
		case this <- i:
		case <-time.After(t):
		}
	}
}
