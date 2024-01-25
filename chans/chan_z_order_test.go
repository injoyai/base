package chans

import (
	"testing"
	"time"
)

func TestNewOrder(t *testing.T) {
	o := NewOrder()
	x1 := o.New()
	x2 := o.New()
	x3 := o.New()
	go func() {
		for i := 0; i < 3; i++ {
			x1.Do(func() {
				t.Log(1)
			})
		}
	}()
	go func() {
		for i := 0; i < 3; i++ {
			x2.Do(func() {
				t.Log(2)
			})
		}
	}()
	go func() {
		for i := 0; i < 3; i++ {
			x3.Do(func() {
				t.Log(3)
			})
		}
	}()
	<-time.After(time.Second * 5)
}
