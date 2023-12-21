package chans

import (
	"testing"
	"time"
)

func TestNewTimer(t *testing.T) {
	{
		x := NewCounter(0)
		for range Traverse(2, time.Second) {
			t.Log(x.Add())
		}
	}
	{
		x := NewCounter(10)
		for range Traverse(10, time.Second) {
			t.Log(x.Add())
		}
	}

}
