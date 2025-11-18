package types

import (
	"testing"
)

func TestList_Split(t *testing.T) {
	ls := List[int]{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	t.Log(ls.Split(1)) //[[0] [1] [2] [3] [4] [5] [6] [7] [8] [9] [10]]
	t.Log(ls.Split(3)) //[[0 1 2] [3 4 5] [6 7 8] [9 10]]
	t.Log(ls.Split(4)) //[[0 1 2 3] [4 5 6 7] [8 9 10]]
	t.Log(ls.Split(5)) //[[0 1 2 3 4] [5 6 7 8 9] [10]]
	t.Log(ls.Split(6)) //[[0 1 2 3 4 5] [6 7 8 9 10]]
}

func TestList_Cut(t *testing.T) {
	ls := List[int]{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	t.Log(ls.Limit(1))      //[0]
	t.Log(ls.Limit(1, 0))   //[0]
	t.Log(ls.Limit(1, 9))   //[9]
	t.Log(ls.Limit(1, 10))  //[10]
	t.Log(ls.Limit(1, -1))  //[]
	t.Log(ls.Limit(2, 0))   //[0 1]
	t.Log(ls.Limit(2, -1))  //[0]
	t.Log(ls.Limit(1, -20)) //[]
}

func TestList_Cut1(t *testing.T) {
	ls := List[int]{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	t.Log(ls.Cut(0, 10))     //[0 1 2 3 4 5 6 7 8 9]
	t.Log(ls.Cut(0, -1))     //[0 1 2 3 4 5 6 7 8 9]
	t.Log(ls.Cut(0, 100))    //[0 1 2 3 4 5 6 7 8 9 10]
	t.Log(ls.Cut(-3, -1))    //[8 9]
	t.Log(ls.Cut(-3, 0))     //[]
	t.Log(ls.Cut(-1, 100))   //[10]
	t.Log(ls.Cut(5, -2))     //[5 6 7 8]
	t.Log(ls.Cut(-100, 3))   //[0 1 2]
	t.Log(ls.Cut(-100, 100)) //[0 1 2 3 4 5 6 7 8 9 10]
	t.Log(ls.Cut(-3))        //[8 9 10]
	t.Log(ls.Cut(-100))      //[0 1 2 3 4 5 6 7 8 9 10]
}

func TestList_IsBand(t *testing.T) {
	ls := List[int]{0, 1, 2, 3}
	t.Log(ls.IsBand(func(a, b int) bool { return a > b }))
	t.Log(ls.IsBand(func(a, b int) bool { return a < b }))

	ls = List[int]{0, 2, 1, 3}
	t.Log(ls.IsBand(func(a, b int) bool { return a > b }))
	t.Log(ls.IsBand(func(a, b int) bool { return a < b }))

	ls = List[int]{0, 2, 2, 3}
	t.Log(ls.IsBand(func(a, b int) bool { return a > b }))
	t.Log(ls.IsBand(func(a, b int) bool { return a < b }))
	t.Log(ls.IsBand(func(a, b int) bool { return a <= b }))
}

func TestList_IsSort(t *testing.T) {
	ls := List[int]{0, 2, 2, 3}
	t.Log(ls.IsSort(func(a, b int) bool { return a <= b }))
}

func TestList_MergeAlternate(t *testing.T) {
	{
		ls1 := List[int]{0, 1, 2, 3, 4, 5}
		ls2 := List[int]{4, 8, 6}
		res := ls1.MergeAlternate(ls2)
		t.Log(res)
	}
	{
		ls1 := List[int]{0, 1}
		ls2 := List[int]{4, 8, 6, 9, 7}
		res := ls1.MergeAlternate(ls2)
		t.Log(res)
	}
}
