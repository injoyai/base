package chans

import (
	"testing"
)

func TestUID(t *testing.T) {
	m := make(map[string]byte)
	for {
		uid := UID()
		if _, ok := m[uid]; ok {
			t.Log("重复:", uid)
			return
		} else {
			m[uid] = 0
		}
	}
}
