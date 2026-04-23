package chans

import (
	"io"
	"sync"
)

var _ io.WriteCloser = new(IO)

func NewIO(cap int) *IO {
	return &IO{
		ch:        make(chan []byte, cap),
		closeSign: make(chan struct{}),
	}
}

type IO struct {
	ch        chan []byte
	once      sync.Once
	closeSign chan struct{}
}

// Write 实现io.Writer接口
func (this *IO) Write(p []byte) (n int, err error) {
	b := append([]byte(nil), p...)
	select {
	case <-this.closeSign:
		return 0, io.ErrClosedPipe
	case this.ch <- b:
		return len(b), nil
	}
}

// ReadMessage 兼容老版ios.ReadMessage
func (this *IO) ReadMessage() ([]byte, error) {
	return this.ReadBytes()
}

func (this *IO) ReadBytes() ([]byte, error) {
	select {
	case bs := <-this.ch:
		return bs, nil
	default:
	}

	select {
	case <-this.closeSign:
		return nil, io.EOF
	case bs := <-this.ch:
		return bs, nil
	}
}

func (this *IO) Close() error {
	this.once.Do(func() {
		close(this.closeSign)
	})
	return nil
}

func (this *IO) Closed() bool {
	select {
	case <-this.closeSign:
		return true
	default:
		return false
	}
}

func (this *IO) Len() int { return len(this.ch) }
func (this *IO) Cap() int { return cap(this.ch) }
