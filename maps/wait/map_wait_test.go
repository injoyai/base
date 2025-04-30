package wait

import (
	"errors"
	"strconv"
	"testing"
	"time"
)

func TestNew2(t *testing.T) {
	SetTimeout(time.Second * 5).SetReuse(false)
	key := "test"
	go func() {
		time.Sleep(time.Millisecond)
		Done(key, "1")
	}()
	t.Log(Wait(key))
	go func() {
		time.Sleep(time.Millisecond)
		Done(key, "2")
	}()
	t.Log(Wait(key))
	go func() {
		time.Sleep(time.Millisecond)
		Done(key, "3")
	}()
	t.Log(Wait(key))
}

func TestNew(t *testing.T) {

	SetTimeout(time.Second * 5).SetReuse(false)
	key := "test"
	go func() {
		for i := 0; i < 10; i++ {

			go func() {
				t.Log(Wait(key))
			}()
			time.Sleep(time.Millisecond * 1000)
		}
	}()
	go func() {
		for i := 0; i < 10; i++ {

			go func() {
				t.Log(Wait(key))
			}()
			time.Sleep(time.Millisecond * 1000)
		}
	}()

	go func() {
		for i := 0; ; i++ {
			time.Sleep(time.Second * 3)
			Done(key, "钱测试"+strconv.Itoa(i))

		}
	}()

	select {}
}

func TestErr(t *testing.T) {
	{
		go func() {
			<-time.After(time.Millisecond * 200)
			Done("key", errors.New("错误"))
		}()
		t.Log(Wait("key"))
	}
	{
		go func() {
			<-time.After(time.Millisecond * 200)
			Done("key", nil, errors.New("错误"))
		}()
		t.Log(Wait("key"))
	}
	{
		go func() {
			<-time.After(time.Millisecond * 200)
			Done("key", nil)
		}()
		t.Log(Wait("key"))
	}
	{
		go func() {
			<-time.After(time.Millisecond * 200)
			Done("key", "000")
		}()
		t.Log(Wait("key"))
	}
}
