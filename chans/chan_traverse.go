package chans

import (
	"time"
)

// TraverseInterval 按时间间隔遍历
func TraverseInterval(interval time.Duration) <-chan int {
	return Traverse(-1, interval)
}

// TraverseCount 按次数遍历
func TraverseCount(num int) <-chan int {
	return Traverse(num)
}

// Traverse 遍历 range each traverse
// @num,数量,-1为死循环
// @interval,间隔
func Traverse(num int, interval ...time.Duration) <-chan int {
	c := make(chan int)
	go func() {
		for i := 0; ; i++ {
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

var Count = Traverse
