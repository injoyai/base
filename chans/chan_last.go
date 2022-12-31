package chans

import (
	"errors"
	"time"
)

type Last struct {
	*Entity
	data    chan []byte
	timeout time.Duration
}

func NewLast() *Last {
	c := make(chan []byte)
	e := NewEntity(1)
	e.SetHandler(func(no, num int, data interface{}) {
		select {
		case c <- data.([]byte):
		default:
		}
	})
	return &Last{
		Entity:  e,
		data:    c,
		timeout: time.Second * 10,
	}
}

func (this *Last) SetTimeout(timeout time.Duration) {
	this.timeout = timeout
}

func (this *Last) Read() ([]byte, error) {
	select {
	case data := <-this.data:
		return data, nil
	case <-time.After(this.timeout):
		return nil, errors.New("超时")
	}
}
