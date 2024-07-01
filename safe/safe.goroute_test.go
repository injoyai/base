package safe

import (
	"context"
	"testing"
	"time"
)

func TestNewGoroute(t *testing.T) {
	g := NewGoroute()
	g.Go(func(ctx context.Context) {
		<-time.After(time.Second * 3)
	}, func(ctx context.Context) {
		<-time.After(time.Second * 4)
	}, func(ctx context.Context) {
		<-time.After(time.Second * 5)
	})
	t.Log(g.Total())
	<-time.After(time.Second * 3)
	t.Log(g.Active())
	<-g.Done()
	t.Log(g.Total(), g.Active())
}
