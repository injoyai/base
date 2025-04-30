package types

import (
	"encoding/json"
	"github.com/injoyai/conv"
)

type Map[K comparable, V any] map[K]V

// GetVar 实现conv.Extend接口
func (this Map[K, V]) GetVar(key K) *conv.Var {
	return conv.New(this[key])
}

// Extend 获取扩展
func (this Map[K, V]) Extend() conv.Extend[K] {
	return conv.NewExtend[K](this)
}

// Unmarshal 解析到ptr中
func (this Map[K, V]) Unmarshal(ptr any, p ...conv.UnmarshalParam) error {
	return conv.Unmarshal(this, ptr, p...)
}

// Bytes map转字节
func (this Map[K, V]) Bytes() Bytes {
	bs, _ := json.Marshal(this)
	return bs
}

// Merge 合并多个map
func (this Map[K, V]) Merge(m ...Map[K, V]) Map[K, V] {
	for _, v := range m {
		for key, val := range v {
			this[key] = val
		}
	}
	return this
}

type Maps[K comparable, V any] List[Map[K, V]]

func (this Maps[K, V]) Merge() Map[K, V] {
	if len(this) == 0 {
		return nil
	}
	return this[0].Merge(this[1:]...)
}
