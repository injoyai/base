package wait

import (
	"errors"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	takeMap        *maps.Safe
	takeOnce       sync.Once
	defaultTimeout = time.Second * 30
	clearNum       = 100
)

// Take 获取一个等待实例,不存在则新建实例
func Take(keys ...string) *Entity {
	takeOnce.Do(func() { takeMap = maps.NewSafe() })
	key := conv.GetDefaultString("default", keys...)
	w, _ := takeMap.GetOrSetByHandler(key, func() (interface{}, error) {
		return New(defaultTimeout), nil
	})
	return w.(*Entity)
}

// Wait 同步等待响应 等同于Sync
func Wait(key string, timeout ...time.Duration) (interface{}, error) {
	return Take().Wait(key, timeout...)
}

// Sync 同步等待响应
func Sync(key string, timeout ...time.Duration) (interface{}, error) {
	return Take().Sync(key, timeout...)
}

// Async 设置异步执行函数,
// 需要注意的是,不一定会有回调(超时的情况)
// 可以设置超时时间为0,则表示一直等待回调
// 对于异步而言,超时就代表回调函数的生命周期
func Async(key string, f Handler, num int, timeout ...time.Duration) {
	Take().Async(key, f, num, timeout...)
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

/*



 */

func New(timeout time.Duration) *Entity {
	return &Entity{
		m:       maps.NewSafe(),
		timeout: timeout,
		reuse:   false,
	}
}

// Entity 等待列表
type Entity struct {
	m       *maps.Safe
	timeout time.Duration //超时时间
	reuse   bool          //复用模式,相同的ke可以使用同一结果
}

// SetReuse 设置数据复用,例如同时下发了几个相同的任务,只会下发一个命令,结果由几个任务共享
func (this *Entity) SetReuse(b ...bool) *Entity {
	this.reuse = !(len(b) > 0 && !b[0])
	return this
}

// SetTimeout 设置全局等待时间
func (this *Entity) SetTimeout(t time.Duration) *Entity {
	this.timeout = t
	return this
}

func (this *Entity) IsWait(key string) bool {
	return this.m.Exist(key)
}

// Wait 同步等待数据响应
func (this *Entity) Wait(key string, timeouts ...time.Duration) (interface{}, error) {
	return this.Sync(key, timeouts...)
}

// Sync 同步等待数据响应,可以设置单次等待时间
func (this *Entity) Sync(key string, timeouts ...time.Duration) (interface{}, error) {
	timeout := conv.GetDefaultDuration(this.timeout, timeouts...)
	w, _ := this.m.GetOrSetByHandler(key, func() (interface{}, error) {
		return newAsync(), nil
	})
	return w.(*async).sync(timeout, this.reuse)
}

// Async 异步执行函数
func (this *Entity) Async(key string, f Handler, num int, timeouts ...time.Duration) {
	timeout := conv.GetDefaultDuration(this.timeout, timeouts...)
	w, _ := this.m.GetOrSetByHandler(key, func() (interface{}, error) {
		return newAsync(), nil
	})
	w.(*async).async(f, timeout, num)
}

// Done 设置回调数据
func (this *Entity) Done(key string, v interface{}, err ...error) bool {
	w, ok := this.m.Get(key)
	if ok {
		w.(*async).done(v, err...)
	}
	return ok
}

//===================================================================//

//数据包
type data struct {
	data interface{} //数据
	err  error       //错误
}

type Handler func(v interface{}, e error)

func newAsync() *async {
	a := &async{syncDone: make(chan struct{}, 1)}
	a.syncDone <- struct{}{}
	return a
}

// async 异步
type async struct {
	list       []*asyncItem
	syncResult []chan *data
	syncDone   chan struct{}
	number     int32
}

// _syncResult
func (this *async) result() <-chan *data {
	c := make(chan *data, 1)
	this.syncResult = append(this.syncResult, c)
	return c
}

// sync 同步等待回调
func (this *async) sync(timeout time.Duration, reuse bool) (v interface{}, e error) {

	timer := time.NewTimer(timeout)
	if reuse {
		select {
		case <-this.syncDone: //没有任务在执行
			timer.Stop()
		case <-timer.C:
			return nil, errors.New("超时")
		case result := <-this.result():
			timer.Stop()
			//复用相同key的响应数据
			return result.data, result.err
		}
	} else {
		select {
		case <-this.syncDone:
			timer.Stop()
		case <-timer.C:
			return nil, errors.New("超时")
		}
	}

	defer func() {
		for _, c := range this.syncResult {
			select {
			case c <- &data{data: v, err: e}:
			}
			close(c)
		}
		this.syncResult = []chan *data(nil)
		//需要在结果通道的后面执行
		this.syncDone <- struct{}{}
	}()

	c := make(chan *data, 1)
	this.async(func(v interface{}, e error) {
		select {
		case c <- &data{data: v, err: e}:
		default:
		}
	}, timeout, 1)

	select {
	case result := <-c:
		return result.data, result.err
	case <-time.After(timeout):
		return nil, errors.New("同步超时")
	}
}

// async 设置异步函数
func (this *async) async(f Handler, timeout time.Duration, num int) {
	this.list = append(this.list, &asyncItem{f: f, timeout: timeout, start: time.Now(), number: int32(num)})
}

// 数据回调,执行异步函数,现在异步超时是过滤了,todo 待实现超时异步回调
func (this *async) done(v interface{}, err ...error) {
	e := conv.GetDefaultErr(nil, err...)
	invalid := map[int]bool{}
	for i, t := range this.list {
		if !t.do(v, e) {
			invalid[i] = true
		}
	}
	//当无效回调函数过多时,清理一次
	if len(invalid) > clearNum {
		list := []*asyncItem(nil)
		for i, l := range this.list {
			if !invalid[i] {
				list = append(list, l)
			}
		}
		this.list = list
	}
}

type asyncItem struct {
	start   time.Time
	timeout time.Duration
	number  int32
	f       Handler
}

func (this *asyncItem) do(v interface{}, err error) bool {
	if this.timeout > 0 && time.Since(this.start) > this.timeout {
		return false
	}
	//负数表示一直有效
	if atomic.LoadInt32(&this.number) == 0 {
		return false
	}
	this.f(v, err)
	atomic.AddInt32(&this.number, -1)
	return true
}
