package chans

import (
	"errors"
	"fmt"
	"github.com/injoyai/conv"
	"net"
	"time"
)

var _ net.Conn = new(Net)

func NewNet[T conv.Integer](cap T, timeout ...time.Duration) *Net {
	i := NewIO(cap, timeout...)
	return &Net{
		IO:   i,
		addr: netAddr(fmt.Sprintf("%p", i)),
	}
}

type Net struct {
	IO            *IO
	readDeadline  time.Time
	writeDeadline time.Time
	addr          net.Addr
}

func (this *Net) Read(b []byte) (n int, err error) {
	if !this.readDeadline.IsZero() && time.Now().Sub(this.readDeadline) > 0 {
		return 0, errors.New("read timeout")
	}
	return this.IO.Read(b)
}

func (this *Net) Write(b []byte) (n int, err error) {
	if !this.writeDeadline.IsZero() && time.Now().Sub(this.writeDeadline) > 0 {
		return 0, errors.New("write timeout")
	}
	//this.IO.Timeout = time.Now().Sub(this.writeDeadline)
	return this.IO.Write(b)
}

func (this *Net) Close() error {
	return this.IO.Close()
}

func (this *Net) LocalAddr() net.Addr {
	return this.addr
}

func (this *Net) RemoteAddr() net.Addr {
	return this.addr
}

func (this *Net) SetDeadline(t time.Time) error {
	if err := this.SetReadDeadline(t); err != nil {
		return err
	}
	return this.SetWriteDeadline(t)
}

func (this *Net) SetReadDeadline(t time.Time) error {
	this.readDeadline = t
	return nil
}

func (this *Net) SetWriteDeadline(t time.Time) error {
	this.writeDeadline = t
	return nil
}

type netAddr string

func (this netAddr) Network() string {
	return "memory"
}

func (this netAddr) String() string {
	return string(this)
}
