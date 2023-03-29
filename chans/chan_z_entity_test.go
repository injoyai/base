package chans

import (
	"context"
	"errors"

	"testing"
	"time"
)

func TestNewEntity(t *testing.T) {
	type TestEntity struct {
		write *Entity
	}
	a := &TestEntity{
		write: NewEntity(1, 3),
	}

	time.Sleep(time.Second)
	a.write.SetHandler(func(ctx context.Context, no, num int, i interface{}) {
		t.Log("序号:", no)
		if i == nil {
			t.Log("is nil")
		}
		if fn, ok := i.(func()); ok {
			fn()
		}
		time.Sleep(time.Second)
		//t.Log(i)
	})

	go func() {
		t2 := time.Now()
		for i := 0; i < 10000; i++ {
			go func(i int) {

				//time.Sleep(time.Second)
				fn := func() {
					t.Log("钱测试", i)
				}
				a.write.Do(fn, fn, fn, fn, fn, fn)
			}(i)
		}

		t.Log(time.Now().Sub(t2))
		//return
	}()
	time.Sleep(time.Second * 5)
	for {
		a.write.CloseWithErr(errors.New("关闭"))
	}
	//a.write.Close("ces")
	//a.write.Close("ces2")
	select {}
}
