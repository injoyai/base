package safe

import (
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
	Try(func() error {
		panic(66)
	}).Catch(func(err error) {
		t.Log(err)
	}).Finally(func(err error) {
		t.Log("done")
	})
}
