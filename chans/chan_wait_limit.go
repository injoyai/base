package chans

import (
	"sync"
)

type WaitLimit interface {
	Add()
	Done()
	Wait()
}

// NewWaitLimit sync.WaitGroup基础上加了同时释放的数量
func NewWaitLimit(limit int) WaitLimit {
	return &waitLimit{c: make(chan struct{}, limit)}
}

type waitLimit struct {
	wg sync.WaitGroup
	c  chan struct{}
}

func (this *waitLimit) Add() {
	this.wg.Add(1)
	this.c <- struct{}{}
}

func (this *waitLimit) Done() {
	<-this.c
	this.wg.Done()
}

func (this *waitLimit) Wait() {
	this.wg.Wait()
}
