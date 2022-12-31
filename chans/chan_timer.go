package chans

// Timer 计数器
type Timer struct {
	c    chan struct{}
	done chan struct{}
}

func NewTimer(times uint) *Timer {
	return &Timer{
		c:    make(chan struct{}, times),
		done: make(chan struct{}, 1),
	}
}

func (this *Timer) Reset() {
	select {
	case <-this.done:
	default:
	}
	this.c = make(chan struct{}, cap(this.c))
}

// Add 计数
func (this *Timer) Add() bool {
	select {
	case this.c <- struct{}{}:
	default:
		select {
		case this.done <- struct{}{}:
		default:
		}
		return true
	}
	if cap(this.c) == len(this.c) {
		return this.Add()
	}
	return false
}

func (this *Timer) Wait() {
	<-this.done
}
