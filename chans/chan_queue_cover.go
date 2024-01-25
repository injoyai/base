package chans

import (
	"errors"
	"time"
)

func NewQueueCover(cap int) *QueueCover {
	if cap <= 0 {
		cap = 1
	}
	return &QueueCover{
		c: make(chan interface{}, cap),
	}
}

// QueueCover 消息队列,百万次数据存入0.07s
// 如果队列已满,后进会把消息整体前移,最前的消息会丢弃
type QueueCover struct {
	c chan interface{}
}

func (this *QueueCover) Append(data interface{}) {
	select {
	case this.c <- data:
	default:
		select {
		case <-this.c:
		default:
		}
		this.Append(data)
	}
}

func (this *QueueCover) Len() int {
	return len(this.c)
}

func (this *QueueCover) Cap() int {
	return cap(this.c)
}

func (this *QueueCover) Pop() interface{} {
	data := <-this.c
	return data
}

// PopMust 取出 same <- this.c
func (this *QueueCover) PopMust() interface{} {
	data := <-this.c
	return data
}

// PopTimeout 取出数据,时间限制
func (this *QueueCover) PopTimeout(timeout time.Duration) (interface{}, error) {
	select {
	case data := <-this.c:
		return data, nil
	case <-time.After(timeout):
		return nil, errors.New("取出数据超时")
	}
}

// PopTry 尝试取出数据,返回数据和是否取出
func (this *QueueCover) PopTry() (interface{}, bool) {
	select {
	case data := <-this.c:
		return data, true
	default:
		return nil, false
	}
}
