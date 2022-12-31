package chans

import (
	"sync"
)

// NewWaitLimit sync.WaitGroup基础上加了同时释放的数量
func NewWaitLimit(limit uint) *WaitLimit {
	return &WaitLimit{c: make(chan struct{}, limit)}
}

type WaitLimit struct {
	sync.WaitGroup
	c chan struct{}
}

func (this *WaitLimit) Add() {
	this.WaitGroup.Add(1)
	this.c <- struct{}{}
}

func (this *WaitLimit) Done() {
	<-this.c
	this.WaitGroup.Done()
}

func (this *WaitLimit) Wait() {
	this.WaitGroup.Wait()
}
