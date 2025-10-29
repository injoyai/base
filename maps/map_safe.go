package maps

import (
	"github.com/injoyai/base/chans"
	"github.com/injoyai/base/types"
	"github.com/injoyai/conv"
	"sync"
	"time"
)

type (
	Safe = Generic[any, any]
)

// NewSafe 新建
func NewSafe() *Safe {
	return NewGeneric[any, any]()
}

func NewGeneric[K comparable, V any]() *Generic[K, V] {
	e := &Generic[K, V]{
		m: WithMutex[K, *Value[V]](),
	}
	e.ExtendGeneric = conv.NewExtendGeneric[K](e)
	return e
}

// Generic 读写分离,适合读多写少
// 千万次写速度 8.3s,
// 千万次读速度 2.3s
type Generic[K comparable, V any] struct {
	m                     types.Mapper[K, *Value[V]] //接口
	hmu                   sync.Map                   //函数锁
	clearOnce             sync.Once                  //清理过期数据,单次执行
	conv.ExtendGeneric[K]                            //接口

	onGet []func(K, V, bool) //读取数据事件
	onSet []func(K, V)       //设置数据事件
	onDel []func(K)          //删除数据事件

	subscribe     *chans.Subscribe[K, V]
	subscribeOnce sync.Once
}

// Exist 是否存在
func (this *Generic[K, V]) Exist(key K) bool {
	_, has := this.Get(key)
	return has
}

// Has 是否存在
func (this *Generic[K, V]) Has(key K) bool {
	_, has := this.Get(key)
	return has
}

// Get 获取数据
func (this *Generic[K, V]) Get(key K) (V, bool) {
	val, has := this.m.Get(key)
	if has {
		v, exist := val.Val()
		this._onGet(key, v, exist)
		return v, exist
	}
	var zero V
	this._onGet(key, zero, false)
	return zero, false
}

func (this *Generic[K, V]) _onGet(key K, value V, exist bool) {
	for _, f := range this.onGet {
		if f != nil {
			f(key, value, exist)
		}
	}
}

// MustGet 获取数据,不管是否存在
func (this *Generic[K, V]) MustGet(key K) V {
	value, _ := this.Get(key)
	return value
}

// GetVar 实现conv.Extend的接口
func (this *Generic[K, V]) GetVar(key K) *conv.Var {
	val, _ := this.Get(key)
	return conv.New(val)
}

// Set 设置数据,可选有效期
func (this *Generic[K, V]) Set(key K, value V, expiration ...time.Duration) {
	this.m.Set(key, NewValue[V](value, expiration...))
	//设置数据事件
	for _, f := range this.onSet {
		if f != nil {
			f(key, value)
		}
	}
}

// SetExpiration 设置有效期
func (this *Generic[K, V]) SetExpiration(key K, expiration time.Duration) bool {
	val, has := this.m.Get(key)
	if has {
		val.SetExpiration(expiration)
	}
	return has
}

// Del 删除键
func (this *Generic[K, V]) Del(key K) {
	this.m.Del(key)
	this.hmu.Delete(key)
	//删除事件
	for _, f := range this.onDel {
		if f != nil {
			f(key)
		}
	}
}

// GetAndDel 获取数据,并删除数据
func (this *Generic[K, V]) GetAndDel(key K) (V, bool) {
	value, has := this.Get(key)
	if !has {
		return value, has
	}
	this.Del(key)
	return value, has
}

// GetAndSet 获取老数据,并设置新数据
func (this *Generic[K, V]) GetAndSet(key K, value V, expiration ...time.Duration) (V, bool) {
	val, has := this.Get(key)
	this.Set(key, value, expiration...)
	return val, has
}

// GetOrSet 尝试获取数据,存在则返回数据,不存在的话存储传入的值,并返回出去,一般使用GetOrSetByHandler
func (this *Generic[K, V]) GetOrSet(key K, value V, expiration ...time.Duration) (V, bool) {
	val, has := this.Get(key)
	if !has {
		this.Set(key, value, expiration...)
		val = value
	}
	return val, has
}

// GetOrSetByHandler
// 尝试获取数据,存在则直接返回数据,
// 不存在的话调用函数,生成数据,储存并返回最新数据
// 执行函数时,增加了锁,避免并发,瞬时大量请求
// check-lock-check
func (this *Generic[K, V]) GetOrSetByHandler(key K, handler func() (V, error), expiration ...time.Duration) (V, error) {
	val, has := this.Get(key)
	if !has && handler != nil {
		muAny, _ := this.hmu.LoadOrStore(key, &sync.Mutex{})
		mu := muAny.(*sync.Mutex)
		mu.Lock()
		defer mu.Unlock()
		val, has = this.Get(key)
		if !has {
			value, err := handler()
			if err != nil {
				var zero V
				return zero, err
			}
			this.Set(key, value, expiration...)
			val = value
		}
	}
	return val, nil
}

// GetOrSetByHandler2 无错误
// 尝试获取数据,存在则直接返回数据,
// 不存在的话调用函数,生成数据,储存并返回最新数据
// 执行函数时,增加了锁,避免并发,瞬时大量请求
// check-lock-check
func (this *Generic[K, V]) GetOrSetByHandler2(key K, handler func() V, expiration ...time.Duration) V {
	val, has := this.Get(key)
	if !has && handler != nil {
		muAny, _ := this.hmu.LoadOrStore(key, &sync.Mutex{})
		mu := muAny.(*sync.Mutex)
		mu.Lock()
		defer mu.Unlock()
		val, has = this.Get(key)
		if !has {
			value := handler()
			this.Set(key, value, expiration...)
			val = value
		}
	}
	return val
}

// Range 遍历数据,返回false结束遍历
func (this *Generic[K, V]) Range(fn func(key K, value V) bool) {
	this.m.Range(func(key K, value *Value[V]) bool {
		v, _ := value.Val()
		return fn(key, v)
	})
}

// Map 复制数据到map[any]any
func (this *Generic[K, V]) Map() map[K]V {
	m := map[K]V{}
	this.Range(func(key K, value V) bool {
		m[key] = value
		return true
	})
	return m
}

// GMap 复制数据到map[string]any
func (this *Generic[K, V]) GMap() map[string]any {
	m := map[string]any{}
	this.Range(func(key K, value V) bool {
		m[conv.String(key)] = value
		return true
	})
	return m
}

// Clone 复制数据 todo
func (this *Generic[K, V]) Clone() *Generic[K, V] {
	m := NewGeneric[K, V]()
	this.m.Range(func(key K, value *Value[V]) bool {
		m.m.Set(key, value)
		return true
	})
	return m
}

// Clear 清除过期数据
func (this *Generic[K, V]) Clear() {
	this.m.Range(func(key K, value *Value[V]) bool {
		if !value.Valid() {
			this.m.Del(key)
		}
		return true
	})
}

// RunClear 定时清理过期数据
func (this *Generic[K, V]) RunClear(interval time.Duration) *Generic[K, V] {
	this.clearOnce.Do(func() {
		go func() {
			for {
				<-time.After(interval)
				this.Clear()
			}
		}()
	})
	return this
}

/*



 */

func (this *Generic[K, V]) OnGet(f func(k K, v V, exist bool)) {
	this.onGet = append(this.onGet, f)
}

func (this *Generic[K, V]) OnSet(f func(k K, v V)) {
	this.onSet = append(this.onSet, f)
}

func (this *Generic[K, V]) OnDel(f func(k K)) {
	this.onDel = append(this.onDel, f)
}

// Chan 订阅特定key的数据
func (this *Generic[K, V]) Chan(key K, cap ...int) *chans.Safe[V] {
	return this.Subscribe(key, cap...)
}

// Subscribe 订阅指定key的数据
func (this *Generic[K, V]) Subscribe(key K, cap ...int) *chans.Safe[V] {
	this.subscribeOnce.Do(func() {
		this.subscribe = chans.NewSubscribe[K, V]()
		this.OnSet(func(k K, v V) {
			this.subscribe.Publish(k, v, 0)
		})
	})
	return this.subscribe.Subscribe(key, cap...)
}
