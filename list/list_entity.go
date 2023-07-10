package list

import (
	"fmt"
	"github.com/injoyai/conv"
	"sort"
	"sync"
)

func New() *Entity {
	return &Entity{}
}

// Entity 列表
type Entity struct {
	list     []interface{}
	mu       sync.Mutex
	lessFunc func(i, j int) bool
}

// getIdx 处理下标,支持负数-1表示最后1个,同python
func (this *Entity) getIdx(idx int) int {
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
func (this *Entity) Len() int {
	return len(this.list)
}

// Cap 总长
func (this *Entity) Cap() int {
	return cap(this.list)
}

// Less 实现排序接口,比较元素
func (this *Entity) Less(i, j int) bool {
	return this.lessFunc == nil || this.lessFunc(i, j)
}

// Swap 实现排序接口,交换元素
func (this *Entity) Swap(i, j int) {
	this.list[i], this.list[j] = this.list[j], this.list[i]
}

// Cut 剪切,新值 , 同 list[start:end]
func (this *Entity) Cut(start, end int) *Entity {
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
	list := &Entity{}
	if end2 > start2 {
		list.Append(this.list[start2:end2]...)
	}
	return list
}

// GoString 对应%#v
func (this *Entity) GoString() string {
	return fmt.Sprintf("%#v", this.list)
}

// String 对应%s
func (this *Entity) String() string {
	return fmt.Sprintf("%v", this.list)
}

// Copy 复制
func (this *Entity) Copy() *Entity {
	x := &this
	return *x
}

// Exist 元素是否存在
func (this *Entity) Exist(idx int) bool {
	return this.getIdx(idx) >= 0
}

// Get 获取元素
func (this *Entity) Get(idx int) (interface{}, bool) {
	if idx = this.getIdx(idx); idx >= 0 {
		return this.list[idx], true
	}
	return nil, false
}

// MustGet 获取元素,不存在返回nil
func (this *Entity) MustGet(idx int) interface{} {
	if idx = this.getIdx(idx); idx >= 0 {
		return this.list[idx]
	}
	return nil
}

// GetVar 获取数据转成*conv.Var
func (this *Entity) GetVar(idx int) *conv.Var {
	return conv.New(this.MustGet(idx))
}

// GetAndDel 获取元素并删除
func (this *Entity) GetAndDel(idx int) (interface{}, bool) {
	if idx = this.getIdx(idx); idx >= 0 {
		val := this.list[idx]
		this.list = append(this.list[:idx], this.list[idx+1:]...)
		return val, true
	}
	return nil, false
}

// MustGetAndDel 获取元素(不存在返回nil)并删除
func (this *Entity) MustGetAndDel(idx int) interface{} {
	v, _ := this.GetAndDel(idx)
	return v
}

// MoveToLast 获取数据并移到队列最后
func (this *Entity) MoveToLast(idx int) interface{} {
	if idx = this.getIdx(idx); idx >= 0 {
		v := this.MustGet(idx)
		this.Remove(idx)
		this.Append(v)
		return v
	}
	return nil
}

// Insert 中间插入元素
func (this *Entity) Insert(idx int, v ...interface{}) {
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
func (this *Entity) Replace(idx int, v interface{}) {
	if idx = this.getIdx(idx); idx > 0 {
		this.list[idx] = v
	}
}

// Join 拼接*List的元素
func (this *Entity) Join(list ...*Entity) {
	for _, v := range list {
		this.Append(v.List()...)
	}
}

// Append 从最后添加元素
func (this *Entity) Append(v ...interface{}) {
	for _, k := range v {
		if list, ok := k.(*Entity); ok {
			this.list = append(this.list, list.List()...)
		} else {
			this.list = append(this.list, k)
		}
	}
}

// Del 移除元素
func (this *Entity) Del(idx ...int) {
	this.Remove(idx...)
}

// Delete 移除元素
func (this *Entity) Delete(idx ...int) {
	this.Remove(idx...)
}

// Remove 移除元素
func (this *Entity) Remove(idx ...int) {
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
func (this *Entity) RemoveNil() {
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

// Reverse 倒序
func (this *Entity) Reverse() *Entity {
	for i := 0; i < this.Len()/2; i++ {
		this.Swap(i, this.Len()-1-i)
	}
	return this
}

// Clear 清除元素
func (this *Entity) Clear() {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.list = []interface{}{}
}

// List 全部列表,引用类型
func (this *Entity) List() []interface{} {
	return this.list
}

// Sort 排序
func (this *Entity) Sort(fn func(a, b interface{}) bool) (err error) {
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
