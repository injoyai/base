package maps

import (
	"github.com/injoyai/conv"
	"time"
)

type Value[V any] struct {
	Value V     //值
	valid int64 //纳秒
}

func (this *Value[V]) Var() *conv.Var {
	v, _ := this.Val()
	return conv.New(v)
}

// Val 获取值,返回值和是否存在(是否有效)
func (this *Value[V]) Val() (V, bool) {
	if this.Valid() {
		return this.Value, true
	}
	var zero V
	return zero, false
}

// Valid 是否在有效期内,即数据是否有效
// 有效期小于等于0表示永久有效,不过期
func (this *Value[V]) Valid() bool {
	return this.valid <= 0 || this.valid > time.Now().UnixNano()
}

// SetExpiration 设置有效期
func (this *Value[V]) SetExpiration(expiration ...time.Duration) {
	if !this.Valid() {
		return
	}
	if len(expiration) > 0 && expiration[0] > 0 {
		this.valid = time.Now().Add(expiration[0]).UnixNano()
	}
}

func (this *Value[V]) Set(value V, expiration ...time.Duration) *Value[V] {
	this.Value = value
	if len(expiration) > 0 && expiration[0] > 0 {
		this.valid = time.Now().Add(expiration[0]).UnixNano()
	}
	return this
}

func NewValue[V any](value V, expiration ...time.Duration) *Value[V] {
	v := &Value[V]{Value: value}
	v.SetExpiration(expiration...)
	return v
}
