package safe

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"
)

func TestNewRerun(t *testing.T) {
	x := NewRerun()
	x.Update(&dial{name: "第一", num: 0})
	x.Update(&dial{name: "第二", num: 10})

	<-time.After(time.Second * 60)
}

type dial struct {
	name  string
	num   int
	retry int
}

func (this *dial) Dial(ctx context.Context) error {
	this.retry++
	log.Println(this.name, "dial", this.retry)
	if this.retry < 3 {
		return errors.New("err")
	}
	this.num++
	return nil
}

func (this *dial) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
			log.Println(this.num)
		}
	}
}

func (this *dial) Close() error {
	return nil
}
