package chans

import (
	"testing"
	"time"
)

func TestNewQueueFunc(t *testing.T) {
	x := NewQueueFunc(3)
	for i := range Count(100) {
		go func(i int) {
			x.Do(func(no int, num int) {
				time.Sleep(time.Second)
				t.Log(no, num)
			})
		}(i)
	}
	<-time.After(time.Second * 30)
}
