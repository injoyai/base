package heart

import (
	"github.com/injoyai/conv"
	"testing"
	"time"
)

func TestNewHeart(t *testing.T) {
	a := NewHeart(time.Minute/6, func(strings []*Info) {
		t.Log("在线数量:", len(strings))
		t.Log("在线:", strings)
	}, func(strings []string) {
		t.Log("离线数量:", len(strings))
		t.Log("离线:", strings)
	})

	i := 10
	for {
		time.Sleep(time.Millisecond * 100)
		for i2 := 0; i2 < i; i2++ {
			time.Sleep(time.Second)
			a.Keep(conv.String(i2))

		}
		i--
		if i <= 0 {
			time.Sleep(time.Second * 20)
			i = 10
		}
	}

}

func TestNewHeart2(t *testing.T) {

}
