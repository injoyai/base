package safe

import (
	"testing"
	"time"
)

func TestNewCloser(t *testing.T) {
	c := NewCloser()
	go func() { <-time.After(time.Second * 5); c.Close() }()
	<-c.Done()
}
