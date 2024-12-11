package chans

import "sync"

type LimitGo interface {
	Go(f func())
	Wait()
}

func NewLimitGo(limit int) LimitGo {
	return &limitGo{c: make(chan struct{}, limit)}
}

type limitGo struct {
	wg sync.WaitGroup
	c  chan struct{}
}

func (this *limitGo) Go(f func()) {
	this.wg.Add(1)
	this.c <- struct{}{}
	go func() {
		defer func() {
			<-this.c
			this.wg.Done()
		}()
		f()
	}()
}

func (this *limitGo) Wait() {
	this.wg.Wait()
}
