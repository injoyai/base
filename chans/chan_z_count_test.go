package chans

import (
	"testing"
	"time"
)

func TestTraverse(t *testing.T) {
	for a := range Traverse(10, time.Second) {
		t.Log(a)
	}
}
