package chans

import (
	"time"
)

// Count 遍历
// @num,数量,-1为死循环
// @interval,间隔
func Count(num int, interval ...time.Duration) <-chan int {
	a := make(chan int)
	go func() {
		for i := 0; ; i++ {
			if num >= 0 && i >= num {
				break
			}
			if len(interval) > 0 && interval[0] > 0 {
				time.Sleep(interval[0])
			}
			a <- i
		}
		close(a)
	}()
	return a
}
