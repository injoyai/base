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
	once           sync.Once
	Default        *Entity
	defaultTimeout = time.Second * 30
)

// Default 获取一个等待实例,不存在则新建实例
func _default() *Entity {
	once.Do(func() {
		Default = New(defaultTimeout)
	})
	return Default
}

// Wait 同步等待响应 等同于Sync
func Wait(key string, timeout ...time.Duration) (any, error) {
	return _default().Wait(key, timeout...)
}

// Sync 同步等待响应
func Sync(key string, timeout ...time.Duration) (any, error) {
	return _default().Sync(key, timeout...)
}

// Async 设置异步执行函数,
// 需要注意的是,不一定会有回调(超时的情况)
// 可以设置超时时间为0,则表示一直等待回调
// 对于异步而言,超时就代表回调函数的生命周期
func Async(key string, f Handler[any], num int, timeout ...time.Duration) {
	_default().Async(key, f, num, timeout...)
}

// SetTimeout 设置超时时间
func SetTimeout(t time.Duration) *Entity {
	return _default().SetTimeout(t)
}

// SetReuse 复用模式
// 如果前面有相同key在等待,结束后会直接将结果赋值给所有相同key的等待结果
// 否则相同key会排队挨个等待
func SetReuse(b ...bool) *Entity {
	return _default().SetReuse(b...)
}

// IsWait 是否有相同key在等待
func IsWait(key string) bool {
	return _default().IsWait(key)
}

// Done 完成等待,给结果赋值,返回key是否存在
func Done(key string, v any, err ...error) bool {
	return _default().Done(key, v, err...)
}

/*



 */

// New 获取一个等待实例,对泛型封装太多层会编译失败
// New(timeout time.Duration) *Entity{return NewGeneric[any, any](timeout)} 这样就会编译失败
func New(timeout time.Duration) *Entity {
	return &Generic[any, any]{
		m:        maps.NewGeneric[any, *async[any]](),
		timeout:  timeout,
		reuse:    false,
		clearNum: 100,
	}
}

func NewGeneric[K comparable, V any](timeout time.Duration) *Generic[K, V] {
	return &Generic[K, V]{
		m:        maps.NewGeneric[K, *async[V]](),
		timeout:  timeout,
		reuse:    false,
		clearNum: 100,
	}
}

type Entity = Generic[any, any]

// Generic 等待列表
type Generic[K comparable, V any] struct {
	m        *maps.Generic[K, *async[V]]
	timeout  time.Duration //超时时间
	reuse    bool          //复用模式,相同的ke可以使用同一结果
	clearNum int           //清理数量
}

// SetReuse 设置数据复用,例如同时下发了几个相同的任务,只会下发一个命令,结果由几个任务共享
func (this *Generic[K, V]) SetReuse(b ...bool) *Generic[K, V] {
	this.reuse = len(b) == 0 || b[0]
	return this
}

// SetTimeout 设置全局等待时间
func (this *Generic[K, V]) SetTimeout(t time.Duration) *Generic[K, V] {
	this.timeout = t
	return this
}

// SetClearNum 设置清理数量
func (this *Generic[K, V]) SetClearNum(num int) *Generic[K, V] {
	this.clearNum = num
	return this
}

func (this *Generic[K, V]) IsWait(key K) bool {
	return this.m.Exist(key)
}

// Wait 同步等待数据响应
func (this *Generic[K, V]) Wait(key K, timeouts ...time.Duration) (V, error) {
	return this.Sync(key, timeouts...)
}

// Sync 同步等待数据响应,可以设置单次等待时间
func (this *Generic[K, V]) Sync(key K, timeouts ...time.Duration) (V, error) {
	timeout := conv.Default[time.Duration](this.timeout, timeouts...)
	w, _ := this.m.GetOrSetByHandler(key, func() (*async[V], error) {
		return newAsync[V](this.clearNum), nil
	})
	return w.sync(timeout, this.reuse)
}

// Async 异步执行函数
func (this *Generic[K, V]) Async(key K, f Handler[V], num int, timeouts ...time.Duration) {
	timeout := conv.Default[time.Duration](this.timeout, timeouts...)
	w, _ := this.m.GetOrSetByHandler(key, func() (*async[V], error) {
		return newAsync[V](this.clearNum), nil
	})
	w.async(f, timeout, num)
}

// Done 设置回调数据
func (this *Generic[K, V]) Done(key K, v V, err ...error) bool {
	w, ok := this.m.GetAndDel(key)
	if ok {
		w.done(v, err...)
	}
	return ok
}

//===================================================================//

// 数据包
type data[V any] struct {
	data V     //数据
	err  error //错误
}

type Handler[V any] func(v V, e error)

func newAsync[V any](clearNum int) *async[V] {
	a := &async[V]{
		syncDone: make(chan struct{}, 1),
		clearNum: clearNum,
	}
	a.syncDone <- struct{}{}
	return a
}

// async 异步
type async[V any] struct {
	list       []*asyncItem[V]
	syncResult []chan *data[V]
	syncDone   chan struct{}
	number     int32
	clearNum   int
}

// _syncResult
func (this *async[V]) result() <-chan *data[V] {
	c := make(chan *data[V], 1)
	this.syncResult = append(this.syncResult, c)
	return c
}

// sync 同步等待回调
func (this *async[V]) sync(timeout time.Duration, reuse bool) (v V, e error) {

	timer := time.NewTimer(timeout)
	if reuse {
		select {
		case <-this.syncDone: //没有任务在执行
			timer.Stop()
		case <-timer.C:
			var zero V
			return zero, errors.New("超时")
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
			var zero V
			return zero, errors.New("超时")
		}
	}

	defer func() {
		for _, c := range this.syncResult {
			select {
			case c <- &data[V]{data: v, err: e}:
			}
			close(c)
		}
		this.syncResult = []chan *data[V](nil)
		//需要在结果通道的后面执行
		this.syncDone <- struct{}{}
	}()

	c := make(chan *data[V], 1)
	this.async(func(v V, e error) {
		select {
		case c <- &data[V]{data: v, err: e}:
		default:
		}
	}, timeout, 1)

	select {
	case result := <-c:
		return result.data, result.err
	case <-time.After(timeout):
		var zero V
		return zero, errors.New("超时")
	}
}

// async 设置异步函数
func (this *async[V]) async(f Handler[V], timeout time.Duration, num int) {
	this.list = append(this.list, &asyncItem[V]{f: f, timeout: timeout, start: time.Now(), number: int32(num)})
}

// 数据回调,执行异步函数,现在异步超时是过滤了,todo 待实现超时异步回调
func (this *async[V]) done(v V, errs ...error) {
	err := conv.Default[error](nil, errs...)
	invalid := map[int]bool{}
	for i, t := range this.list {
		if !t.do(v, err) {
			invalid[i] = true
		}
	}
	//当无效回调函数过多时,清理一次
	if len(invalid) > this.clearNum {
		list := []*asyncItem[V](nil)
		for i, l := range this.list {
			if !invalid[i] {
				list = append(list, l)
			}
		}
		this.list = list
	}
}

type asyncItem[V any] struct {
	start   time.Time     //开始时间,配合有效时间使用
	timeout time.Duration //超时时间,有效期
	number  int32         //执行次数
	current int32         //当前次数
	f       Handler[V]    //回调函数
}

// do 执行回调函数,返回是否有效
func (this *asyncItem[V]) do(v V, err error) bool {
	if this == nil {
		return false
	}
	if this.timeout > 0 && time.Since(this.start) > this.timeout {
		return false
	}
	//负数表示一直有效
	if this.number < 0 || this.current < this.number {
		current := atomic.AddInt32(&this.current, 1)
		if this.number < 0 || current <= this.number {
			if this.f != nil {
				this.f(v, err)
			}
		}
	}
	return true
}
