package chans

import (
	"context"
	"testing"
	"time"
)

func TestNewQueueFunc(t *testing.T) {
	x := NewQueueFunc(3)
	for i := range Traverse(100) {
		go func(i int) {
			x.Do(func(ctx context.Context, no int, num int) {
				time.Sleep(time.Second)
				t.Log(no, num)
			})
		}(i)
	}
	<-time.After(time.Second * 30)
}
