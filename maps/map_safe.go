package maps

import (
	"github.com/injoyai/base/chans"
	"github.com/injoyai/conv"
	"sync"
	"time"
)

type (
	SafeAny = Safe[any, any]
	SafeSA  = Safe[string, any]
)

func NewSafeAny() *SafeAny {
	return NewSafe[any, any]()
}

func NewSafeDefault() *SafeSA {
	return NewSafe[string, any]()
}

// NewSafe 新建
func NewSafe[K comparable, V any]() *Safe[K, V] {
	e := &Safe[K, V]{
		m: WithMutex[K, *Value[V]](),
	}
	//e.Extend = conv.NewExtend(e)
	return e
}

// Safe 读写分离,适合读多写少
// 千万次写速度 8.3s,
// 千万次读速度 2.3s
type Safe[K comparable, V any] struct {
	m              Mapper[K, *Value[V]] //接口
	hmu            sync.Map             //函数锁
	listened       bool                 //是否数据监听
	listen         sync.Map             //数据监听
	clearOnce      sync.Once            //清理过期数据,单次执行
	conv.Extend[K]                      //接口
}

// Exist 是否存在
func (this *Safe[K, V]) Exist(key K) bool {
	_, has := this.Get(key)
	return has
}

// Has 是否存在
func (this *Safe[K, V]) Has(key K) bool {
	_, has := this.Get(key)
	return has
}

// Get 获取数据
func (this *Safe[K, V]) Get(key K) (V, bool) {
	val, has := this.m.Get(key)
	if has {
		return val.Val()
	}
	var zero V
	return zero, false
}

// MustGet 获取数据,不管是否存在
func (this *Safe[K, V]) MustGet(key K) V {
	value, _ := this.Get(key)
	return value
}

// GetVar 实现conv.Extend的接口
func (this *Safe[K, V]) GetVar(key K) *conv.Var {
	val, _ := this.Get(key)
	return conv.New(val)
}

// Set 设置数据,可选有效期
func (this *Safe[K, V]) Set(key K, value V, expiration ...time.Duration) {
	this.m.Set(key, NewValue[V](value, expiration...))
	if this.listened {
		listen, ok := this.listen.Load(key)
		if ok {
			listen.(*chans.Subscribe[V]).Publish(value)
		}
	}
}

// SetExpiration 设置有效期
func (this *Safe[K, V]) SetExpiration(key K, expiration time.Duration) bool {
	val, has := this.m.Get(key)
	if has {
		val.SetExpiration(expiration)
	}
	return has
}

// Del 删除键
func (this *Safe[K, V]) Del(key K) {
	this.m.Del(key)
	this.hmu.Delete(key)
}

// GetAndDel 获取数据,并删除数据
func (this *Safe[K, V]) GetAndDel(key K) (V, bool) {
	value, has := this.Get(key)
	if !has {
		return value, has
	}
	this.Del(key)
	return value, has
}

// GetAndSet 获取老数据,并设置新数据
func (this *Safe[K, V]) GetAndSet(key K, value V, expiration ...time.Duration) (V, bool) {
	val, has := this.Get(key)
	this.Set(key, value, expiration...)
	return val, has
}

// GetOrSet 尝试获取数据,存在则返回数据,不存在的话存储传入的值,并返回出去,一般使用GetOrSetByHandler
func (this *Safe[K, V]) GetOrSet(key K, value V, expiration ...time.Duration) (V, bool) {
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
func (this *Safe[K, V]) GetOrSetByHandler(key K, handler func() (V, error), expiration ...time.Duration) (V, error) {
	val, has := this.Get(key)
	if !has && handler != nil {
		muAny, _ := this.hmu.LoadOrStore(key, &sync.Mutex{})
		mu := muAny.(*sync.Mutex)
		mu.Lock()
		defer mu.Unlock()
		val, has = this.Get(key)
		if !has && handler != nil {
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

// Range 遍历数据,返回false结束遍历
func (this *Safe[K, V]) Range(fn func(key K, value V) bool) {
	this.m.Range(func(key K, value *Value[V]) bool {
		v, _ := value.Val()
		return fn(key, v)
	})
}

// Map 复制数据到map[any]any
func (this *Safe[K, V]) Map() map[K]V {
	m := map[K]V{}
	this.Range(func(key K, value V) bool {
		m[key] = value
		return true
	})
	return m
}

// GMap 复制数据到map[string]any
func (this *Safe[K, V]) GMap() map[string]any {
	m := map[string]any{}
	this.Range(func(key K, value V) bool {
		m[conv.String(key)] = value
		return true
	})
	return m
}

// Clone 复制数据 todo
func (this *Safe[K, V]) Clone() *Safe[K, V] {
	m := NewSafe[K, V]()
	this.m.Range(func(key K, value *Value[V]) bool {
		m.m.Set(key, value)
		return true
	})
	return m
}

//// Writer 将特定的key写入实现成io.Writer
//func (this *Safe[K, V]) Writer(key K) io.Writer {
//	return WriteFunc(func(p []byte) (int, error) {
//		this.Set(key, p)
//		return len(p), nil
//	})
//}

// Chan 订阅特定key的数据
func (this *Safe[K, V]) Chan(key any, cap ...uint) *chans.Safe[V] {
	this.listened = true
	l, ok := this.listen.Load(key)
	if !ok {
		l = chans.NewSubscribe[V]()
		this.listen.Store(key, l)
	}
	return l.(*chans.Subscribe[V]).Subscribe(cap...)
}

// Clear 清除过期数据
func (this *Safe[K, V]) Clear() {
	this.m.Range(func(key K, value *Value[V]) bool {
		if !value.Valid() {
			this.m.Del(key)
		}
		return true
	})
}

// RunClear 定时清理过期数据
func (this *Safe[K, V]) RunClear(interval time.Duration) *Safe[K, V] {
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
