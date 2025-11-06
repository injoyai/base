package safe

import (
	"testing"
	"time"
)

func TestNewChan(t *testing.T) {
	c := NewChan[int]()
	go func() {
		<-time.After(time.Second * 3)
		c.Close()
	}()
	<-c.Done()
	t.Log(c.Closed())
	c.Close()
	t.Log("完成")
}
