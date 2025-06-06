package chans

import (
	"errors"
	"fmt"
	"github.com/injoyai/base/types"
	"github.com/injoyai/conv"
	"io"
	"sync/atomic"
	"time"
)

var _ io.ReadWriteCloser = new(IO)

func NewIO[T conv.Integer](cap T, timeout ...time.Duration) *IO {
	return &IO{
		C:       make(types.Chan[[]byte], cap),
		Timeout: conv.Default[time.Duration](-1, timeout...),
	}
}

type IO struct {
	C         types.Chan[[]byte]
	Timeout   time.Duration
	readCache []byte
	closed    uint32
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
	if !this.C.Add(p, this.Timeout) {
		return 0, nil
	}
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

	if len(this.readCache) == 0 {
		this.readCache, err = this.ReadMessage()
		if err != nil {
			return
		}
	}

	//从缓存(上次剩余的字节)复制数据到p
	n = copy(p, this.readCache)
	if n < len(this.readCache) {
		this.readCache = this.readCache[n:]
		return
	}

	this.readCache = nil
	return
}

func (this *IO) Close() error {
	if atomic.CompareAndSwapUint32(&this.closed, 0, 1) {
		close(this.C)
	}
	return nil
}

func (this *IO) Closed() bool {
	return atomic.LoadUint32(&this.closed) == 1
}
