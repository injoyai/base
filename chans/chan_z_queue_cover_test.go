package chans

import "testing"

func TestNewQueue(t *testing.T) {
	q := NewQueueCover(10)
	t.Log(q.PopTry())
	for i := 0; i < 1000000; i++ {
		q.Append(i)
	}
	t.Log(q.PopMust())
	t.Log(q.PopMust())
}
