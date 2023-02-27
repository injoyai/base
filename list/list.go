package list

import (
	"fmt"
	"github.com/injoyai/conv"
	"sort"
	"sync"
)

func New() *List {
	return &List{}
}

// List 列表
type List struct {
	list     []interface{}
	mu       sync.Mutex
	next     int
	lessFunc func(i, j int) bool
}

// getIdx 处理下标,支持负数-1表示最后1个,同python
func (this *List) getIdx(idx int) int {
	length := this.Len()
	if idx < length && idx >= 0 {
		return idx
	}
	if idx < 0 && -idx <= length {
		return length + idx
	}
	return -1
}

// Len 元素长度
func (this *List) Len() int {
	return len(this.list)
}

// Cap 总长
func (this *List) Cap() int {
	return cap(this.list)
}

// Less 实现排序接口,比较元素
func (this *List) Less(i, j int) bool {
	return this.lessFunc == nil || this.lessFunc(i, j)
}

// Swap 实现排序接口,交换元素
func (this *List) Swap(i, j int) {
	this.list[i], this.list[j] = this.list[j], this.list[i]
}

// Cut 剪切,新值 , 同 list[start:end]
func (this *List) Cut(start, end int) *List {
	if this.Len() == 0 {
		return New()
	}
	start2 := this.getIdx(start)
	end2 := this.getIdx(end)
	if end > 0 && end2 < 0 {
		end2 = this.Len()
	}
	if start2 < 0 {
		start2 = 0
	}
	list := &List{}
	if end2 > start2 {
		list.Append(this.list[start2:end2]...)
	}
	return list
}

// NextIdx 下一个元素下标
func (this *List) NextIdx() int {
	return this.next
}

// Next 获取下一个数据
func (this *List) Next() interface{} {
	if this.Len() == 0 {
		return nil
	}
	if this.next < this.Len() {
		v := this.MustGet(this.next)
		this.next++
		if this.next >= this.Len() {
			this.next = 0
		}
		return v
	}
	this.next = 0
	return this.Next()
}

// GoString 对应%#v
func (this *List) GoString() string {
	return fmt.Sprintf("%#v", this.list)
}

// String 对应%s
func (this *List) String() string {
	return fmt.Sprintf("%v", this.list)
}

// Copy 复制
func (this *List) Copy() *List {
	x := &this
	return *x
}

// Exist 元素是否存在
func (this *List) Exist(idx int) bool {
	return this.getIdx(idx) >= 0
}

// Get 获取元素
func (this *List) Get(idx int) (interface{}, bool) {
	if idx = this.getIdx(idx); idx >= 0 {
		return this.list[idx], true
	}
	return nil, false
}

// MustGet 获取元素,不存在返回nil
func (this *List) MustGet(idx int) interface{} {
	if idx = this.getIdx(idx); idx >= 0 {
		return this.list[idx]
	}
	return nil
}

// GetVar 获取数据转成*conv.Var
func (this *List) GetVar(idx int) *conv.Var {
	return conv.New(this.MustGet(idx))
}

// MoveToLast 获取数据并移到队列最后
func (this *List) MoveToLast(idx int) interface{} {
	if idx = this.getIdx(idx); idx >= 0 {
		v := this.MustGet(idx)
		this.Remove(idx)
		this.Append(v)
		return v
	}
	return nil
}

// Insert 中间插入元素
func (this *List) Insert(idx int, v ...interface{}) {
	if idx = this.getIdx(idx); idx >= 0 {
		this.mu.Lock()
		defer this.mu.Unlock()
		length := len(v)
		list := make([]interface{}, this.Len()+length)
		for i, k := range this.list {
			if i < idx {
				list[i] = k
			} else {
				list[i+length] = k
			}
		}
		for i, k := range v {
			list[idx+i] = k
		}
		this.list = list
	}
}

// Replace 替换元素,替换已有的元素
func (this *List) Replace(idx int, v interface{}) {
	if idx = this.getIdx(idx); idx > 0 {
		this.list[idx] = v
	}
}

// Join 拼接*List的元素
func (this *List) Join(list ...*List) {
	for _, v := range list {
		this.Append(v.List()...)
	}
}

// Append 从最后添加元素
func (this *List) Append(v ...interface{}) {
	for _, k := range v {
		if list, ok := k.(*List); ok {
			this.list = append(this.list, list.List()...)
		} else {
			this.list = append(this.list, k)
		}
	}
}

// Delete 移除元素
func (this *List) Delete(idx ...int) {
	this.Remove(idx...)
}

// Remove 移除元素
func (this *List) Remove(idx ...int) {
	if this.Len() > 0 {
		this.mu.Lock()
		defer this.mu.Unlock()
		list := []interface{}(nil)
		m := map[int]bool{}
		for _, i := range idx {
			m[i] = true
		}
		for i, v := range this.list {
			if !m[i] {
				list = append(list, v)
			}
		}
		this.list = list
	}
}

// RemoveNil 移除nil的元素
func (this *List) RemoveNil() {
	if this.Len() > 0 {
		this.mu.Lock()
		defer this.mu.Unlock()
		list := []interface{}(nil)
		for _, v := range this.list {
			if v != nil {
				list = append(list, v)
			}
		}
		this.list = list
	}
}

// Clear 清除元素
func (this *List) Clear() {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.list = []interface{}{}
	this.next = 0
}

// List 全部列表,引用类型
func (this *List) List() []interface{} {
	return this.list
}

// Sort 排序
func (this *List) Sort(fn func(a, b interface{}) bool) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	this.lessFunc = func(i, j int) bool {
		return fn(this.list[i], this.list[j])
	}
	sort.Sort(this)
	return
}
