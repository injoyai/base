package chans

import (
	"context"
	"errors"
	"time"
)

type Last struct {
	*Entity
	data chan []byte
}

func NewLast() *Last {
	c := make(chan []byte)
	e := NewEntity(1)
	e.SetHandler(func(ctx context.Context, no, num int, data interface{}) {
		select {
		case <-c:
		default:
		}
		select {
		case c <- data.([]byte):
		default:
		}
	})
	return &Last{
		Entity: e,
		data:   c,
	}
}

func (this *Last) ReadWithTimeout(timeout time.Duration) ([]byte, error) {
	select {
	case data := <-this.data:
		return data, nil
	case <-time.After(timeout):
		return nil, errors.New("超时")
	}
}

func (this *Last) Read() []byte {
	return <-this.data
}
