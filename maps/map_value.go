package maps

import (
	"github.com/injoyai/conv"
	"time"
)

type Value struct {
	Value      interface{}
	valid      time.Time
	expiration time.Duration
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
func (this *Value) Valid() bool {
	return this.valid.Unix() <= 0 || this.valid.Sub(time.Now()) > 0
}

// SetExpiration 设置有效期
func (this *Value) SetExpiration(expiration ...time.Duration) {
	if !this.Valid() {
		return
	}
	this.expiration = conv.GetDefaultDuration(0, expiration...)
	this.valid = func() time.Time {
		if len(expiration) > 0 && expiration[0] > 0 {
			return time.Now().Add(expiration[0])
		}
		return time.Unix(0, 0)
	}()
}

// ResetExpiration 重置有效期
func (this *Value) ResetExpiration() {
	this.SetExpiration(this.expiration)
}

func newValue(value interface{}, expiration ...time.Duration) *Value {
	v := &Value{Value: value}
	v.SetExpiration(expiration...)
	return v
}
