package chans

import (
	"github.com/injoyai/conv"
	"time"
)

// Range 仿python的range 参数1-3个
// 当1个参数, 例 Range(5) 输出 0,1,2,3,4
// 当2个参数, 例 Range(1,5) 输出 1,2,3,4
// 当3个参数, 例 Range(0,5,2) 输出 0,2,4
func Range[T conv.Integer](n T, m ...T) <-chan T {
	start, end, step := T(0), n, T(1)
	switch len(m) {
	case 0:
	case 1:
		start, end = n, m[0]
	default:
		start, end, step = n, m[0], m[1]
	}
	c := make(chan T)
	go func() {
		for i := start; i < end; i += step {
			c <- i
		}
		close(c)
	}()
	return c
}

// TraverseInterval 按时间间隔遍历
func TraverseInterval(interval time.Duration) <-chan int {
	return Traverse(-1, interval)
}

// TraverseCount 按次数遍历
func TraverseCount[T conv.Integer](num T) <-chan T {
	return Traverse(num)
}

// Traverse 遍历 range each traverse
// @num,数量,-1为死循环
// @interval,间隔
func Traverse[T conv.Integer](num T, interval ...time.Duration) <-chan T {
	c := make(chan T)
	go func() {
		for i := T(0); ; i++ {
			if num >= 0 && i >= num {
				break
			}
			if len(interval) > 0 && interval[0] > 0 {
				time.Sleep(interval[0])
			}
			c <- i
		}
		close(c)
	}()
	return c
}

func Count[T conv.Integer](num T, interval ...time.Duration) <-chan T {
	return Traverse(num, interval...)
}
