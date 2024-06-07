package timeout

import (
	"context"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	x := New()
	x.SetTimeout(time.Second * 10)
	x.SetInterval(time.Second)
	x.SetDealFunc(func(key interface{}) error {
		t.Log("超时: ", key)
		return nil
	})
	x.Keep(1)

	go x.Run(context.Background())
	<-time.After(time.Second * 12)
}
