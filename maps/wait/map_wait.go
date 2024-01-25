package wait

import (
	"errors"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"sync"
	"time"
)

var (
	m    *maps.Safe
	once sync.Once
)

// Take 获取一个等待实例,不存在则新建实例
func Take(keys ...string) *Entity {
	once.Do(func() { m = maps.NewSafe() })
	key := conv.GetDefaultString("default", keys...)
	if val := m.GetVar(key); !val.IsNil() {
		return val.Val().(*Entity)
	}
	newWait := New(time.Second * 30)
	m.Set(key, newWait)
	return newWait
}

// Wait 等待
func Wait(key string, timeout ...time.Duration) (interface{}, error) {
	return Take().Wait(key, timeout...)
}

// SetTimeout 设置超时时间
func SetTimeout(t time.Duration) *Entity {
	return Take().SetTimeout(t)
}

// SetReuse 复用模式
// 如果前面有相同key在等待,结束后会直接将结果赋值给所有相同key的等待结果
// 否则相同key会排队挨个等待
func SetReuse(b ...bool) *Entity {
	return Take().SetReuse(b...)
}

// IsWait 是否有相同key在等待
func IsWait(key string) bool {
	return Take().IsWait(key)
}

// Done 完成等待,给结果赋值,返回key是否存在
func Done(key string, v interface{}, err ...error) bool {
	return Take().Done(key, v, err...)
}

// Entity 等待列表
type Entity struct {
	m       map[string]*wait //map
	mu      sync.RWMutex     //锁
	timeout time.Duration    //超时时间
	reuse   bool             //复用模式
}

func New(timeout time.Duration) *Entity {
	return &Entity{
		m:       make(map[string]*wait),
		timeout: timeout,
		reuse:   false,
	}
}

func (this *Entity) SetReuse(b ...bool) *Entity {
	this.reuse = !(len(b) > 0 && !b[0])
	return this
}

func (this *Entity) SetTimeout(t time.Duration) *Entity {
	this.timeout = t
	return this
}

func (this *Entity) IsWait(key string) bool {
	this.mu.Lock()
	_, ok := this.m[key]
	this.mu.Unlock()
	return ok
}

func (this *Entity) Wait(key string, timeouts ...time.Duration) (interface{}, error) {

	timeout := this.timeout
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}

	this.mu.RLock()
	w, ok := this.m[key]
	this.mu.RUnlock()

	if ok {
		if this.reuse {
			data := <-w.result()
			return data.data, data.err
		} else {
			timer := time.NewTimer(timeout)
			select {
			case <-w.finish():
				timer.Stop()
			case <-timer.C:
				return nil, errors.New("超时")
			}
		}
	}

	w = newWait(timeout)
	this.mu.Lock()
	this.m[key] = w
	this.mu.Unlock()

	data, err := w.wait()

	this.mu.Lock()
	defer this.mu.Unlock()
	delete(this.m, key)
	return data, err
}

func (this *Entity) Done(key string, v interface{}, err ...error) bool {
	this.mu.Lock()
	defer this.mu.Unlock()
	_, ok := this.m[key]
	if ok {
		if len(err) > 0 {
			this.m[key].done(err[0])
		} else {
			this.m[key].done(v)
		}
		delete(this.m, key)
	}
	return ok
}

//===================================================================//

//数据包
type data struct {
	data interface{} //数据
	err  error       //错误
}

func newData(v interface{}, err error) *data {
	return &data{
		data: v,
		err:  err,
	}
}

//等待机制,消费一次
type wait struct {
	c       chan interface{} //等待通道
	timeout time.Duration    //超时时间
	_finish chan uintptr     //结束信号
	_result []chan *data     //结果
}

func newWait(timeout time.Duration) *wait {
	return &wait{
		c:       make(chan interface{}, 1),
		timeout: timeout,
		_finish: make(chan uintptr, 1),
	}
}

func (this *wait) finish() <-chan uintptr {
	return this._finish
}

func (this *wait) result() <-chan *data {
	c := make(chan *data)
	this._result = append(this._result, c)
	return c
}

func (this *wait) reResult() {
	this._result = []chan *data{}
}

//等待回调
func (this *wait) wait() (data interface{}, err error) {
	timer := time.NewTimer(this.timeout)
	defer func() {
		timer.Stop()
		select {
		case this._finish <- 0:
		default:
		}
		for _, v := range this._result {
			v <- newData(data, err)
			close(v)
		}
		this.reResult()
	}()
	select {
	case v := <-this.c:
		switch val := v.(type) {
		case error:
			err = val
		default:
			data = v
		}
	case <-timer.C:
		err = errors.New("超时")
	}
	return
}

//收到回调,结束等待
func (this *wait) done(v interface{}) {
	select {
	case this.c <- v:
	default:
	}
}
