package chans

import "sync"

func NewSignal() *Signal {
	return &Signal{
		c:    make(chan struct{}),
		once: sync.Once{},
	}
}

type Signal struct {
	c    chan struct{}
	once sync.Once
}

func (this *Signal) Close() error {
	this.once.Do(func() {
		close(this.c)
	})
	return nil
}

func (this *Signal) Done() <-chan struct{} {
	return this.c
}

func (this *Signal) Wait() {
	<-this.Done()
}
