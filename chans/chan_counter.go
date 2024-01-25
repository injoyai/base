package chans

// Counter 计数器
type Counter struct {
	c    chan struct{}
	done chan struct{}
}

func NewCounter(times uint) *Counter {
	return &Counter{
		c:    make(chan struct{}, times),
		done: make(chan struct{}, 1),
	}
}

func (this *Counter) Reset() {
	select {
	case <-this.done:
	default:
	}
	this.c = make(chan struct{}, cap(this.c))
}

// Add 计数
func (this *Counter) Add() bool {
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

func (this *Counter) Wait() {
	<-this.done
}
