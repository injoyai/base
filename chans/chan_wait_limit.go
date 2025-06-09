package chans

import (
	"sync"
)

// NewWaitLimit sync.WaitGroup基础上加了同时释放的数量
func NewWaitLimit(limit int) *WaitLimit {
	return &WaitLimit{c: make(chan struct{}, limit)}
}

type WaitLimit struct {
	wg sync.WaitGroup
	c  chan struct{}
}

func (this *WaitLimit) Add() {
	this.wg.Add(1)
	this.c <- struct{}{}
}

func (this *WaitLimit) Done() {
	<-this.c
	this.wg.Done()
}

func (this *WaitLimit) Wait() {
	this.wg.Wait()
}
