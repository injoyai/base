package safe

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewRunOne(t *testing.T) {
	n := uint32(0)
	x := NewRunOne(func(ctx context.Context) error {
		atomic.AddUint32(&n, 1)
		for {
			select {
			case <-ctx.Done():
				t.Log("close", n)
				return ctx.Err()
			case <-time.After(time.Second):
				t.Log("run", n)
			}
		}
	})
	for i := 0; i < 10; i++ {
		<-time.After(time.Millisecond * 100)
		go x.Run()
	}
	time.Sleep(time.Second * 10)
}
