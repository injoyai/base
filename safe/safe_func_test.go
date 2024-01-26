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
