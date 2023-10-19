package maps

import (
	"github.com/injoyai/conv"
	"sync"
)

// map合集,不开放,防止数据类型强转错误
var (
	defaultMaps *Safe
	defaultOnce sync.Once
)

// Take 名字待定
func Take(keys ...string) *Safe {
	key := conv.GetDefaultString("_default", keys...)
	defaultOnce.Do(func() { defaultMaps = NewSafe() })
	val := defaultMaps.GetVar(key)
	if !val.IsNil() {
		value, ok := val.Val().(*Safe)
		if ok {
			return value
		}
	}
	newMap := NewSafe()
	defaultMaps.Set(key, newMap)
	return newMap
}
