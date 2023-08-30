package maps

import (
	"github.com/injoyai/conv"
	"time"
)

type Value struct {
	Value interface{}
	Valid time.Time
}

func (this *Value) Var() *conv.Var {
	v, _ := this.Val()
	return conv.New(v)
}

func (this *Value) Val() (interface{}, bool) {
	if this.Valid.Unix() == 0 || this.Valid.Sub(time.Now()) > 0 {
		return this.Value, true
	}
	return nil, false
}

func newValue(value interface{}, expiration ...time.Duration) *Value {
	return &Value{
		Value: value,
		Valid: func() time.Time {
			if len(expiration) > 0 && expiration[0] > 0 {
				return time.Now().Add(expiration[0])
			}
			return time.Unix(0, 0)
		}(),
	}
}
