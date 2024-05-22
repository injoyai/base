package list

import (
	"github.com/injoyai/conv"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	x := New(nil, nil, nil, nil)
	b := x.List()
	t.Log(b) //[<nil> <nil> <nil> <nil>]
	x.RemoveNil()
	x.Append(0)
	b = append(b, 10)
	t.Log(b) //[<nil> <nil> <nil> <nil>,10]
	x.Append("xx")
	t.Log(x.List())  //[0 xx]
	t.Log(x.Get(0))  //0 true
	t.Log(x.Get(-1)) //xx true

	//=====//

	x.Append(6, 7, 8, 9)
	t.Log(x.List()) //[0 xx 6 7 8 9]
	x.Replace(1, 5)
	t.Log(x.List()) //[0 5 6 7 8 9]
	x.Insert(1, 1, 2, 3)
	t.Log(x.List()) //[0 1 2 3 5 6 7 8 9]
	x.Insert(4, 4)
	t.Log(x.List()) //[0 1 2 3 4 5 6 7 8 9]
	t.Log(x.Sort(func(a, b interface{}) bool {
		return a.(int) > b.(int)
	})) //[9 8 7 6 5 3 2 1 0]
	t.Log(x.List())             //[9 8 7 6 5 3 2 1 0]
	t.Log(x.GetVar(0).String()) //9
	t.Log(x)                    //[9 8 7 6 5 4 3 2 1 0]
}

func TestCut(t *testing.T) {
	x := New()
	x.Append(6, 7, 8, 9)
	t.Log(x.Copy().Cut(-1, 0))  //[]
	t.Log(x.Copy().Cut(0, 1))   //[6]
	t.Log(x.Copy().Cut(1, 20))  //[7,8,9]
	t.Log(x.Copy().Cut(-2, -1)) //[8]
	t.Log(x.Copy().Cut(-4, 1))  //[6]
	t.Log(x.Copy().Cut(-4, -1)) //[6,7,8]
	t.Log(x.Copy().Cut(-4, -1)) //[6,7,8]
	t.Log(x.Copy().Cut(-4, 20)) //[6,7,8,9]
	t.Log(x.Copy().Cut(20, -5)) //[]
	x = New()
	t.Log(x.Cut(0, 1)) //[]
}

func TestList_Reverse(t *testing.T) {
	{
		x := New()
		x.Append(1, 2, 3)
		x.Reverse()
		l := x.List()
		t.Log(l) //[3 2 1]
	}
	{
		x := New()
		x.Append(1, 2, 3, 4)
		x.Reverse()
		l := x.List()
		t.Log(l) //[4 3 2 1]
	}
	{
		x := New()
		x.Append(1)
		x.Reverse()
		l := x.List()
		t.Log(l) //[1]
	}
	{
		x := New()
		x.Reverse()
		l := x.List()
		t.Log(l) //[]
	}
}

func TestSQL(t *testing.T) {
	type People struct {
		Name string
		Age  int
		Time time.Time
		Ok   bool
	}
	list := []People(nil)
	for i := 0; i < 20; i++ {
		list = append(list, People{
			Name: "xx" + conv.String(i),
			Age:  i,
			Time: time.Now(),
			Ok:   i%2 == 0,
		})
	}

	x := New(list)
	_ = x
	x.Where(func(i int, v interface{}) bool {
		return v.(People).Age > 5
	})
	x = x.OrderBy(func(v interface{}) interface{} {
		return v.(People).Age
	}).Reverse()
	y := x.Copy()

	t.Log(x.List())
	x.Reverse().Limit(4)
	t.Log(y)
	t.Log(x.List())
	maps := []map[string]interface{}(nil)
	x.Find(&maps)
	t.Log(maps)

}
