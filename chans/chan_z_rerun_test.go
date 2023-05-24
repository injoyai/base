package chans

import (
	"context"
	"testing"
	"time"
)

func TestNewOneFunc(t *testing.T) {
	co := 0
	x := NewRerun(func(ctx context.Context) {
		co++
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
				t.Log(co)
			}
		}
	})
	for {
		<-time.After(time.Second * 3)
		x.Rerun()
	}
}
