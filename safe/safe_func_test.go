package safe

import (
	"errors"
	"testing"
)

func TestRecover(t *testing.T) {

	t.Log(func() (err error) {
		defer Recover(&err)
		panic(666)
		return
	}())

	t.Log(func() (err error) {
		defer Recover(&err, true)
		panic(667)
		return
	}())
}

func TestTry(t *testing.T) {
	if Try(func() error {
		return nil
	}) != nil {
		t.Error("返回错误")
		return
	}
	t.Log(Try(func() error {
		panic("只会打印panic")
	}))
	Try(func() error {
		return errors.New("不会打印正常错误")
	}, func(err error) {
		t.Log(err)
	})
}

func TestIndex(t *testing.T) {
	ls := []int{1, 2, 3, 4, 5}
	t.Log(ls[Index(len(ls), -1)])        //5
	t.Log(ls[Index(len(ls), -2)])        //4
	t.Log(ls[Index(len(ls), -3)])        //3
	t.Log(ls[Index(len(ls), -4)])        //2
	t.Log(ls[Index(len(ls), -5)])        //1
	t.Log(ls[Index(len(ls), 0)])         //1
	t.Log(ls[Index(len(ls), 1)])         //2
	t.Log(ls[Index(len(ls), 2)])         //3
	t.Log(ls[Index(len(ls), 3)])         //4
	t.Log(ls[Index(len(ls), 4)])         //5
	t.Log(ls[Index(len(ls), 5, true)])   //1
	t.Log(ls[Index(len(ls), 6, true)])   //2
	t.Log(ls[Index(len(ls), 7, true)])   //3
	t.Log(ls[Index(len(ls), 8, true)])   //4
	t.Log(ls[Index(len(ls), 9, true)])   //5
	t.Log(ls[Index(len(ls), 10, true)])  //1
	t.Log(ls[Index(len(ls), -6, true)])  //5
	t.Log(ls[Index(len(ls), -7, true)])  //4
	t.Log(ls[Index(len(ls), -8, true)])  //3
	t.Log(ls[Index(len(ls), -9, true)])  //2
	t.Log(ls[Index(len(ls), -10, true)]) //1
	t.Log(ls[Index(len(ls), -11, true)]) //5
}
