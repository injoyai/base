package types

import (
	"sort"

	"github.com/injoyai/conv"
)

type List[T any] []T

// Len 元素长度
func (this List[T]) Len() int {
	return len(this)
}

// Cap 总长
func (this List[T]) Cap() int {
	return cap(this)
}

// Swap 实现排序接口,交换元素
func (this List[T]) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

// Copy 复制
func (this List[T]) Copy() List[T] {
	ls := make(List[T], len(this))
	copy(ls, this)
	return ls
}

// getIdx 处理下标,支持负数-1表示最后1个,同python
func (this List[T]) getIdx(idx int) int {
	length := this.Len()
	if idx < length && idx >= 0 {
		return idx
	}
	if idx < 0 && -idx <= length {
		return length + idx
	}
	return -1
}

// Exist 元素是否存在
func (this List[T]) Exist(idx int) bool {
	return this.getIdx(idx) >= 0
}

// Get 获取元素
func (this List[T]) Get(idx int) (any, bool) {
	if idx = this.getIdx(idx); idx >= 0 {
		return this[idx], true
	}
	return nil, false
}

// MustGet 获取元素,不存在返回nil
func (this List[T]) MustGet(idx int) T {
	if idx = this.getIdx(idx); idx >= 0 {
		return this[idx]
	}
	var zero T
	return zero
}

// GetVar 获取数据转成*conv.Var
func (this List[T]) GetVar(idx int) *conv.Var {
	return conv.New(this.MustGet(idx))
}

// Set 替换元素,替换已有的元素
func (this List[T]) Set(idx int, v T) List[T] {
	if idx = this.getIdx(idx); idx >= 0 {
		this[idx] = v
	}
	return this
}

// Reverse 倒序
func (this List[T]) Reverse() List[T] {
	for i := 0; i < this.Len()/2; i++ {
		this.Swap(i, this.Len()-1-i)
	}
	return this
}

// Sort 排序
func (this List[T]) Sort(fn func(a, b T) bool) List[T] {
	Sort(this, func(i, j int) bool {
		return fn(this[i], this[j])
	})
	return this
}

// Where 筛选数据,参照SQL的命名
func (this List[T]) Where(fn func(i int, v T) bool) List[T] {
	cache := make([]T, 0, this.Len())
	for i, v := range this {
		if fn(i, v) {
			cache = append(cache, v)
		}
	}
	return cache
}

// Del 移除元素
func (this List[T]) Del(idx ...int) List[T] {
	m := make(map[int]bool)
	for _, v := range idx {
		if i := this.getIdx(v); i > 0 {
			m[i] = true
		}
	}
	return this.Where(func(i int, v T) bool {
		return !m[i]
	})
}

// Cut 剪切,新值 , 安全的 list[start:end]
func (this List[T]) Cut(start int, end ...int) List[T] {
	_end := conv.Default(this.Len(), end...)
	cache := List[T]{}
	if this.Len() > 0 {
		start = this.getIdx(start)
		_end = this.getIdx(_end)
		if _end < 0 {
			_end = len(this)
		}
		if start < 0 {
			start = 0
		}
		if start > _end {
			cache = this[:0]
		} else {
			cache = this[start:_end]
		}
	}
	return cache
}

// Limit 取一定数量的数据,可选偏移值
func (this List[T]) Limit(size int, offset ...int) List[T] {
	start := conv.Default(0, offset...)
	end := start + size
	start = conv.Select(start < 0, 0, start)
	end = conv.Select(end > this.Len(), this.Len(), end)
	end = conv.Select(end < 0, 0, end)
	if start > end {
		return this[:0]
	}
	return this[start:end]
}

// Split 分割
func (this List[T]) Split(size int) []List[T] {
	cache := make([]List[T], 0, this.Len()/size)
	for i := 0; i < this.Len(); i += size {
		cache = append(cache, this.Limit(size, i))
	}
	return cache
}

// Unmarshal 解析到ptr中
func (this List[T]) Unmarshal(ptr any, p ...conv.UnmarshalParam) error {
	return conv.Unmarshal(this, ptr, p...)
}

/******/

// MergeAlternate 交替合并
func (this List[T]) MergeAlternate(ls []T) List[T] {
	res := make([]T, 0, this.Len()+len(ls))
	for i := 0; i < this.Len() || i < len(ls); i++ {
		if i < this.Len() {
			res = append(res, this[i])
		}
		if i < len(ls) {
			res = append(res, ls[i])
		}
	}
	return res
}

// IsBand 是否是波段
func (this List[T]) IsBand(f func(a, b T) bool) bool {
	if this.Len() < 2 {
		return true
	}
	for i := 0; i+1 < this.Len(); i += 2 {
		if !f(this[i], this[i+1]) {
			return false
		}
	}
	for i := 1; i+1 < this.Len(); i += 2 {
		if f(this[i], this[i+1]) {
			return false
		}
	}
	return true
}

// IsSort 是否排序
func (this List[T]) IsSort(f func(a, b T) bool) bool {
	if this.Len() < 2 {
		return true
	}
	for i := 0; i+1 < this.Len(); i++ {
		if !f(this[i], this[i+1]) {
			return false
		}
	}
	return true
}

/*



 */

type _sort[T any] struct {
	lessFunc func(i, j int) bool
	Sorter
}

func (this _sort[T]) Less(i, j int) bool {
	return this.lessFunc(i, j)
}

func Sort(s Sorter, fn func(i, j int) bool) {
	sort.Sort(_sort[any]{
		Sorter:   s,
		lessFunc: fn,
	})
}
