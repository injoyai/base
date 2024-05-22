package safe

import (
	"context"
	"testing"
	"time"
)

func TestNewEnable(t *testing.T) {
	e := NewRunner(func(ctx context.Context) error {
		t.Log("启用")
		defer t.Log("结束")
		for i := 0; i < 20; i++ {
			select {
			case <-ctx.Done():
				return nil
			default:
				time.Sleep(time.Second)
			}
		}
		return nil
	})

	t.Log("running: ", e.Running())
	e.Start()
	e.Start()
	<-time.After(time.Second)
	t.Log("running: ", e.Running())
	<-time.After(time.Second)
	t.Log("restart")
	e.Restart()

	<-time.After(time.Second * 30)
	t.Log("完成")
}
