package maps

import (
	"github.com/injoyai/base/chans"
	"github.com/injoyai/conv"
	"io"
	"sync"
	"time"
)

// NewSafe 新建
func NewSafe(m ...Map) *Safe {
	e := &Safe{}
	if len(m) > 0 && m[0] != nil {
		e.m = m[0]
	} else {
		e.m = WithSync()
	}
	e.Extend = conv.NewExtend(e)
	return e
}

// Safe 读写分离,适合读多写少
// 千万次写速度 8.3s,
// 千万次读速度 2.3s
type Safe struct {
	m           Map       //接口
	hmu         sync.Map  //函数锁
	listened    bool      //是否数据监听
	listen      sync.Map  //数据监听
	clearOnce   sync.Once //清理过期数据,单次执行
	conv.Extend           //接口
}

// Exist 是否存在
func (this *Safe) Exist(key interface{}) bool {
	_, has := this.Get(key)
	return has
}

// Has 是否存在
func (this *Safe) Has(key interface{}) bool {
	_, has := this.Get(key)
	return has
}

// Get 获取数据
func (this *Safe) Get(key interface{}) (interface{}, bool) {
	val, has := this.m.Get(key)
	if has {
		return val.(*Value).Val()
	}
	return nil, false
}

// MustGet 获取数据,不管是否存在
func (this *Safe) MustGet(key interface{}) interface{} {
	value, _ := this.Get(key)
	return value
}

// GetVar 实现conv.Extend的接口
func (this *Safe) GetVar(key string) *conv.Var {
	val, _ := this.Get(key)
	return conv.New(val)
}

// Set 设置数据,可选有效期
func (this *Safe) Set(key, value interface{}, expiration ...time.Duration) {
	this.m.Set(key, NewValue(value, expiration...))
	if this.listened {
		listen, ok := this.listen.Load(key)
		if ok {
			listen.(*chans.Listen).Publish(value)
		}
	}
}

// SetExpiration 设置有效期
func (this *Safe) SetExpiration(key string, expiration time.Duration) bool {
	val, has := this.m.Get(key)
	if has {
		val.(*Value).SetExpiration(expiration)
	}
	return has
}

// Del 删除键
func (this *Safe) Del(key interface{}) {
	this.m.Del(key)
	this.hmu.Delete(key)
}

// GetAndDel 获取数据,并删除数据
func (this *Safe) GetAndDel(key interface{}) (interface{}, bool) {
	value, has := this.Get(key)
	if !has {
		return nil, has
	}
	this.Del(key)
	return value, has
}

// GetAndSet 获取老数据,并设置新数据
func (this *Safe) GetAndSet(key, value interface{}, expiration ...time.Duration) (interface{}, bool) {
	val, has := this.Get(key)
	this.Set(key, value, expiration...)
	return val, has
}

// GetOrSet 尝试获取数据,存在则返回数据,不存在的话存储传入的值,并返回出去,一般使用GetOrSetByHandler
func (this *Safe) GetOrSet(key, value interface{}, expiration ...time.Duration) (interface{}, bool) {
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
func (this *Safe) GetOrSetByHandler(key interface{}, handler func() (interface{}, error), expiration ...time.Duration) (interface{}, error) {
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
				return nil, err
			}
			this.Set(key, value, expiration...)
			val = value
		}
	}
	return val, nil
}

// Range 遍历数据,返回false结束遍历
func (this *Safe) Range(fn func(key, value interface{}) bool) {
	this.m.Range(func(key, value interface{}) bool {
		v, _ := value.(*Value).Val()
		return fn(key, v)
	})
}

// Map 复制数据到map[interface{}]interface{}
func (this *Safe) Map() map[interface{}]interface{} {
	m := map[interface{}]interface{}{}
	this.Range(func(key, value interface{}) bool {
		m[key] = value
		return true
	})
	return m
}

// GMap 复制数据到map[string]interface{}
func (this *Safe) GMap() map[string]interface{} {
	m := map[string]interface{}{}
	this.Range(func(key, value interface{}) bool {
		m[conv.String(key)] = value
		return true
	})
	return m
}

// Clone 复制数据 todo
func (this *Safe) Clone() *Safe {
	m := NewSafe()
	this.m.Range(func(key, value interface{}) bool {
		m.m.Set(key, value)
		return true
	})
	return m
}

// Writer 将特定的key写入实现成io.Writer
func (this *Safe) Writer(key interface{}) io.Writer {
	return WriteFunc(func(p []byte) (int, error) {
		this.Set(key, p)
		return len(p), nil
	})
}

// Chan 订阅特定key的数据
func (this *Safe) Chan(key interface{}, cap ...uint) *chans.Subscribe {
	this.listened = true
	l, ok := this.listen.Load(key)
	if !ok {
		l = chans.NewListen()
		this.listen.Store(key, l)
	}
	return l.(*chans.Listen).Subscribe(cap...)
}

// Clear 清除过期数据
func (this *Safe) Clear() {
	this.m.Range(func(key, value interface{}) bool {
		if !value.(*Value).Valid() {
			this.m.Del(key)
		}
		return true
	})
}

// RunClear 定时清理过期数据
func (this *Safe) RunClear(interval time.Duration) *Safe {
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
