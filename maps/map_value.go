package maps

import (
	"github.com/injoyai/conv"
	"time"
)

type Value struct {
	Value interface{} //值
	valid int64       //纳秒
}

func (this *Value) Var() *conv.Var {
	v, _ := this.Val()
	return conv.New(v)
}

// Val 获取值,返回值和是否存在(是否有效)
func (this *Value) Val() (interface{}, bool) {
	if this.Valid() {
		return this.Value, true
	}
	return nil, false
}

// Valid 是否在有效期内,即数据是否有效
// 有效期小于等于0表示永久有效,不过期
func (this *Value) Valid() bool {
	return this.valid <= 0 || this.valid > time.Now().UnixNano()
}

// SetExpiration 设置有效期
func (this *Value) SetExpiration(expiration ...time.Duration) {
	if !this.Valid() {
		return
	}
	if len(expiration) > 0 && expiration[0] > 0 {
		this.valid = time.Now().Add(expiration[0]).UnixNano()
	}
}

func (this *Value) Set(value interface{}, expiration ...time.Duration) *Value {
	this.Value = value
	if len(expiration) > 0 && expiration[0] > 0 {
		this.valid = time.Now().Add(expiration[0]).UnixNano()
	}
	return this
}

func NewValue(value interface{}, expiration ...time.Duration) *Value {
	v := &Value{Value: value}
	v.SetExpiration(expiration...)
	return v
}
