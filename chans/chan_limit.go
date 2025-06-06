package chans

import "github.com/injoyai/conv"

type Limit struct {
	c chan struct{}
}

// NewLimit 同时执行最大数
func NewLimit[T conv.Integer](limit T) *Limit {
	return &Limit{
		c: make(chan struct{}, limit),
	}
}

// Try 尝试,返回是否到达最大限制
func (this *Limit) Try() bool {
	select {
	case this.c <- struct{}{}:
	default:
		return false
	}
	return true
}

// Add 等待加入成功
func (this *Limit) Add() {
	this.c <- struct{}{}
}

// Done 释放,执行完成
func (this *Limit) Done() {
	select {
	case <-this.c:
	default:
	}
}
