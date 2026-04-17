package timeout

import (
	"context"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	x := New()

	x.SetDealFunc(func(key any) error {
		t.Log("超时: ", key)
		return nil
	})
	x.SetTimeout(time.Second * 10)
	x.SetInterval(time.Second)
	x.Keep(1)

	go func() {
		<-time.After(time.Second * 15)
		x.Keep(1)
		x.Keep(2)
	}()

	x.Run(context.Background())

}
