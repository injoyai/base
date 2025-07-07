package types

import (
	"testing"
)

func TestSortMap_Sort(t *testing.T) {
	m := SortMap[string, int]{
		"d": 4,
		"c": 3,
		"a": 1,
		"b": 2,
	}
	t.Log(m.Sort())
	t.Log(m.Sort(false))
	t.Log(m.Sort(true))
}
