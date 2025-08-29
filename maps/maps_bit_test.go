package maps

import (
	"math/rand"
	"testing"
)

func TestNewBit(t *testing.T) {
	{
		m := NewBit()
		m.Set(1, true)
		m.Set(1, true)
		m.Set(1, true)
	}
	{
		m := NewBit()
		for i := uint64(0); i < 100; i++ {
			m.Set(i, true)
			if !m.Get(i) {
				t.Log("错误")
			}
		}
	}
	{
		m := NewBit()
		max := 1000
		list := map[uint64]bool{}
		for i := uint64(0); i < 100; i++ {
			x := uint64(rand.Intn(max))
			m.Set(x, true)
			list[x] = true
		}
		for i := uint64(0); i < uint64(max); i++ {
			if m.Get(i) != list[i] {
				t.Log("错误:", i, m.Get(i), list[i])
			}
		}
	}
}
