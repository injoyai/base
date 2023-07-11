package maps

import (
	"github.com/injoyai/base/chans"
	"github.com/injoyai/conv"
	"io"
	"sync"
	"time"
)

// NewSafe 新建
func NewSafe() *Safe {
	e := &Safe{}
	e.Extend = conv.NewExtend(e)
	return e
}

// Safe 读写分离,适合读多写少
type Safe struct {
	m        sync.Map //存储
	hmu      sync.Map //函数锁
	listened bool     //是否数据监听
	listen   sync.Map //数据监听

	conv.Extend //接口
}

func (this *Safe) Exist(key interface{}) bool {
	_, has := this.Get(key)
	return has
}

func (this *Safe) Has(key interface{}) bool {
	_, has := this.Get(key)
	return has
}

func (this *Safe) Get(key interface{}) (interface{}, bool) {
	val, has := this.m.Load(key)
	if has {
		return val.(*Value).Val()
	}
	return nil, false
}

func (this *Safe) MustGet(key interface{}) interface{} {
	value, _ := this.Get(key)
	return value
}

func (this *Safe) GetVar(key string) *conv.Var {
	val, _ := this.Get(key)
	return conv.New(val)
}

func (this *Safe) Set(key, value interface{}, expiration ...time.Duration) {
	this.m.Store(key, newValue(value, expiration...))
	if this.listened {
		//list, ok := this.c.Load(key)
		//if ok {
		//	for _, c := range list.([]*Chan) {
		//		c.add(value)
		//	}
		//}
		listen, ok := this.listen.Load(key)
		if ok {
			listen.(chans.Listen).Publish(value)
		}
	}
}

func (this *Safe) Del(key interface{}) {
	this.m.Delete(key)
}

func (this *Safe) GetAndDel(key interface{}) interface{} {
	value, has := this.Get(key)
	if !has {
		return nil
	}
	this.Del(key)
	return value
}

func (this *Safe) GetAndSet(key, value interface{}, expiration ...time.Duration) interface{} {
	val, _ := this.Get(key)
	this.Set(key, value, expiration...)
	return val
}

// GetOrSet 尝试获取数据,存在则返回数据,不存在的话存储传入的值,并返回出去,一般使用GetOrSetByHandler
func (this *Safe) GetOrSet(key, value interface{}, expiration ...time.Duration) interface{} {
	val, has := this.Get(key)
	if !has {
		this.Set(key, value, expiration...)
		val = value
	}
	return val
}

// GetOrSetByHandler
// 尝试获取数据,存在则直接返回数据,
// 不存在的话调用函数,生成数据,储存并返回最新数据
// 执行函数时,增加了锁,避免并发,瞬时大量请求
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

func (this *Safe) Range(fn func(key, value interface{}) bool) {
	this.m.Range(func(key, value interface{}) bool {
		v, _ := value.(*Value).Val()
		return fn(key, v)
	})
}

func (this *Safe) Map() map[interface{}]interface{} {
	m := map[interface{}]interface{}{}
	this.Range(func(key, value interface{}) bool {
		m[key] = value
		return true
	})
	return m
}

func (this *Safe) GMap() map[string]interface{} {
	m := map[string]interface{}{}
	this.Range(func(key, value interface{}) bool {
		m[conv.String(key)] = value
		return true
	})
	return m
}

func (this *Safe) Clone() *Safe {
	m := NewSafe()
	this.m.Range(func(key, value interface{}) bool {
		m.m.Store(key, value)
		return true
	})
	return m
}

func (this *Safe) Writer(key interface{}) io.Writer {
	return Write(func(p []byte) (int, error) {
		this.Set(key, p)
		return len(p), nil
	})
}

func (this *Safe) Chan(key interface{}, cap ...uint) *chans.Subscribe {
	this.listened = true
	l, ok := this.listen.Load(key)
	if !ok {
		l = chans.NewListen()
		this.listen.Store(key, l)
	}
	return l.(chans.Listen).Subscribe(cap...)
}
