package chans

type Limit struct {
	running chan struct{} //正在执行的接口
}

// NewLimit 同时执行最大数
func NewLimit(limit int) *Limit {
	return &Limit{
		running: make(chan struct{}, limit),
	}
}

// Do 执行,返回是否到达最大限制
func (this *Limit) Do() bool {
	select {
	case this.running <- struct{}{}:
	default:
		return true
	}
	return false
}

// Done 释放,执行完成
func (this *Limit) Done() {
	select {
	case <-this.running:
	default:
	}
}
