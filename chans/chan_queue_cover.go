package chans

import (
	"errors"
	"github.com/injoyai/base/types"
	"time"
)

func NewQueueCover[T any](cap int) QueueCover[T] {
	if cap <= 0 {
		cap = 1
	}
	return make(QueueCover[T], cap)
}

// QueueCover 消息队列,百万次数据存入0.07s
// 如果队列已满,后进会把消息整体前移,最前的消息会丢弃
type QueueCover[T any] types.Chan[T]

func (c QueueCover[T]) Append(data T) {
	select {
	case c <- data:
	default:
		select {
		case <-c:
		default:
		}
		c.Append(data)
	}
}

func (c QueueCover[T]) Len() int {
	return len(c)
}

func (c QueueCover[T]) Cap() int {
	return cap(c)
}

func (c QueueCover[T]) Pop() any {
	data := <-c
	return data
}

// PopMust 取出 same <- this.c
func (c QueueCover[T]) PopMust() any {
	data := <-c
	return data
}

// PopTimeout 取出数据,时间限制
func (c QueueCover[T]) PopTimeout(timeout time.Duration) (T, error) {
	select {
	case data := <-c:
		return data, nil
	case <-time.After(timeout):
		var zero T
		return zero, errors.New("取出数据超时")
	}
}

// PopTry 尝试取出数据,返回数据和是否取出
func (c QueueCover[T]) PopTry() (T, bool) {
	select {
	case data := <-c:
		return data, true
	default:
		var zero T
		return zero, false
	}
}
