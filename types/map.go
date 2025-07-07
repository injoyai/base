package types

import (
	"encoding/json"
	"github.com/injoyai/conv"
	"sort"
)

type Map[K comparable, V any] map[K]V

// GetVar 实现conv.Extend接口
func (this Map[K, V]) GetVar(key K) *conv.Var {
	return conv.New(this[key])
}

// Extend 获取扩展
func (this Map[K, V]) Extend() conv.ExtendGeneric[K] {
	return conv.NewExtendGeneric[K](this)
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

/**/

type SortMap[K Comparable, V any] map[K]V

func (this SortMap[K, V]) Sort(desc ...bool) []V {
	items := make([]sortMapItem[K, V], 0, len(this))
	for k, v := range this {
		items = append(items, sortMapItem[K, V]{
			K: k,
			V: v,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		b := items[i].K < items[j].K
		if len(desc) > 0 && desc[0] {
			b = !b
		}
		return b
	})
	ret := make([]V, 0, len(items))
	for _, item := range items {
		ret = append(ret, item.V)
	}
	return ret
}

type sortMapItem[K comparable, V any] struct {
	K K
	V V
}

/**/

type Maps[K comparable, V any] List[Map[K, V]]

func (this Maps[K, V]) Merge() Map[K, V] {
	if len(this) == 0 {
		return nil
	}
	return this[0].Merge(this[1:]...)
}
