package chans

import (
	"errors"
	"fmt"
	"github.com/injoyai/conv"
	"io"
	"sync/atomic"
)

var _ io.ReadWriteCloser = new(IO)

func NewIO(cap ...uint) *IO {
	c := make(chan []byte, conv.GetDefaultUint(0, cap...))
	return &IO{C: c}
}

type IO struct {
	C      chan []byte
	cache  []byte
	closed uint32
}

// Write 实现io.Writer接口
func (this *IO) Write(p []byte) (n int, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	if atomic.LoadUint32(&this.closed) == 1 {
		return 0, errors.New("io closed")
	}
	this.C <- p
	return len(p), nil
}

func (this *IO) ReadMessage() ([]byte, error) {
	if atomic.LoadUint32(&this.closed) == 1 {
		return nil, io.EOF
	}

	bs, ok := <-this.C
	if !ok {
		atomic.StoreUint32(&this.closed, 1)
		//这个类型的目的就是为了控制EOF,
		//返回错误的话就不能达到目标效果
		//固这里返回EOF,下同
		return nil, io.EOF
	}
	return bs, nil
}

func (this *IO) Read(p []byte) (n int, err error) {

	if len(this.cache) == 0 {
		this.cache, err = this.ReadMessage()
		if err != nil {
			return
		}
	}

	//从缓存(上次剩余的字节)复制数据到p
	n = copy(p, this.cache)
	if n < len(this.cache) {
		this.cache = this.cache[n:]
		return
	}

	this.cache = nil
	return
}

func (this *IO) Close() error {
	if atomic.CompareAndSwapUint32(&this.closed, 0, 1) {
		close(this.C)
	}
	return nil
}
