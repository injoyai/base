package chans

import (
	"testing"
	"time"
)

func TestNewWait(t *testing.T) {
	x := NewWaitLimit(2)
	for i := 0; i < 10; i++ {
		x.Add()
		go func(i int) {
			defer x.Done()
			time.Sleep(time.Second)
			t.Log(i)
		}(i)

	}
	x.Wait()
}
