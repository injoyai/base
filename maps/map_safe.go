package maps

import (
	"github.com/injoyai/conv"
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
	m           sync.Map //存储
	c           sync.Map //监听通道
	cUsed       bool     //是否使用通道功能,减少查询次数
	conv.Extend          //接口
}

func (this *Safe) Has(key interface{}) bool {
	_, has := this.m.Load(key)
	return has
}

func (this *Safe) Get(key interface{}) (interface{}, bool) {
	val, has := this.m.Load(key)
	if has {
		return val.(*Value).Val(), has
	}
	return val, has
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
	if this.cUsed {
		list, ok := this.c.Load(key)
		if ok {
			for _, c := range list.([]*Chan) {
				c.add(value)
			}
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

func (this *Safe) GetOrSet(key, value interface{}, expiration ...time.Duration) interface{} {
	val, has := this.Get(key)
	if !has {
		this.Set(key, value, expiration...)
		val = value
	}
	return val
}

func (this *Safe) GetOrSetByHandler(key interface{}, handler func() (interface{}, error), expiration ...time.Duration) (interface{}, error) {
	val, has := this.Get(key)
	if !has {
		value, err := handler()
		if err != nil {
			return nil, err
		}
		this.Set(key, value, expiration...)
		val = value
	}
	return val, nil
}

func (this *Safe) Range(fn func(key, value interface{}) bool) {
	this.m.Range(func(key, value interface{}) bool {
		return fn(key, value.(*Value).Val())
	})
}

func (this *Safe) Map() map[interface{}]interface{} {
	m := map[interface{}]interface{}{}
	this.m.Range(func(key, value interface{}) bool {
		m[key] = value
		return true
	})
	return m
}

func (this *Safe) GMap() map[string]interface{} {
	m := map[string]interface{}{}
	this.m.Range(func(key, value interface{}) bool {
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

//========================================Chan========================================

func (this *Safe) Chan(key interface{}, cap ...uint) *Chan {
	return this.typeChan(chanTryInput, key, cap...)
}

func (this *Safe) TryChan(key interface{}, cap ...uint) *Chan {
	return this.typeChan(chanTryInput, key, cap...)
}

func (this *Safe) MustChan(key interface{}, cap ...uint) *Chan {
	return this.typeChan(chanMustInput, key, cap...)
}

func (this *Safe) GoMustChan(key interface{}, cap ...uint) *Chan {
	return this.typeChan(chanGoMustInput, key, cap...)
}

func (this *Safe) TimeoutChan(key interface{}, timeout time.Duration, cap ...uint) *Chan {
	c := this.typeChan(chanTimeoutInput, key, cap...)
	c.inputTimeout = timeout
	return c
}

func (this *Safe) GoTimeoutChan(key interface{}, timeout time.Duration, cap ...uint) *Chan {
	c := this.typeChan(chanGoTimeoutInput, key, cap...)
	c.inputTimeout = timeout
	return c
}

func (this *Safe) typeChan(inputType string, key interface{}, cap ...uint) *Chan {
	this.cUsed = true
	c := newChan(inputType, key, cap...)
	c.setCloseFunc(func() {
		if val, ok := this.c.Load(key); ok {
			list := val.([]*Chan)
			for i, v := range list {
				if v == c {
					val = append(list[:i], list[i+1:]...)
					break
				}
			}
			this.c.Store(key, val)
		}
	})
	if list, ok := this.c.LoadOrStore(key, []*Chan{c}); ok {
		list = append(list.([]*Chan), c)
		this.c.Store(key, list)
	}
	return c
}
