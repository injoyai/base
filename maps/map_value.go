package maps

import (
	"github.com/injoyai/conv"
	"time"
)

type Value struct {
	Value      interface{}
	Valid      time.Time
	expiration time.Duration
}

func (this *Value) Expiration() time.Duration {
	return this.expiration
}

func (this *Value) Var() *conv.Var {
	return conv.New(this.Val())
}

func (this *Value) Val() interface{} {
	if this.Valid.Unix() == 0 || this.Valid.Sub(time.Now()) > 0 {
		return this.Value
	}
	return nil
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
		expiration: conv.GetDefaultDuration(0, expiration...),
	}
}
