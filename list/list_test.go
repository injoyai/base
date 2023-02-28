package list

import (
	"testing"
)

func TestNew(t *testing.T) {
	x := New()
	b := x.List()
	t.Log(b) //[<nil> <nil> <nil> <nil>]
	x.RemoveNil()
	x.Append(0)
	b = append(b, 10)
	t.Log(b) //[<nil> <nil> <nil> <nil>,10]
	x.Append("xx")
	t.Log(x.List()) //[xx]
	t.Log(x.Get(0))
	t.Log(x.Get(-1))

	//=====//

	x.Append(6, 7, 8, 9)
	t.Log(x.List()) //[0 xx 6 7 8 9]
	x.Replace(1, 5)
	t.Log(x.List()) //[0 5 6 7 8 9]
	x.Insert(1, 1, 2, 3)
	t.Log(x.List()) //[0 1 2 3 5 6 7 8 9]
	x.Insert(4, 4)
	t.Log(x.Sort(func(a, b interface{}) bool {
		return a.(int) > b.(int)
	}))
	t.Log(x.List())             //[9 8 7 6 5 3 2 1 0]
	t.Log(x.GetVar(0).String()) //"9"
	t.Log(x)
}

func TestNext(t *testing.T) {
	x := New()
	x.Append(6, 7, 8, 9)
	for i := 0; i < 20; i++ {
		t.Log(x.NextIdx(), x.Next())
	}
}

func TestCut(t *testing.T) {
	x := New()
	x.Append(6, 7, 8, 9)
	t.Log(x.Cut(-1, 0))  //[]
	t.Log(x.Cut(0, 1))   //[6]
	t.Log(x.Cut(1, 20))  //[7,8,9]
	t.Log(x.Cut(-2, -1)) //[8]
	t.Log(x.Cut(-4, 1))  //[6]
	t.Log(x.Cut(-4, -1)) //[6,7,8]
	t.Log(x.Cut(-5, -1)) //[6,7,8]
	t.Log(x.Cut(-5, 20)) //[6,7,8,9]
	t.Log(x.Cut(20, -5)) //[]
	x = New()
	t.Log(x.Cut(0, 1)) //[]
}

func TestList_Reverse(t *testing.T) {
	{
		x := New()
		x.Append(1, 2, 3)
		x.Reverse()
		l := x.List()
		t.Log(l)
	}
}
