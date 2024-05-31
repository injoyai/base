package list

import (
	"fmt"
	"github.com/injoyai/conv"
	"sort"
	"sync"
	"time"
)

func New(v ...interface{}) *List {
	return new(List).append(v...)
}

func NewSafe(v ...interface{}) *List {
	return new(List).Safe().append(v...)
}

// Entity 兼容之前的名称
type Entity = List

// List 列表
type List struct {
	list     []interface{}
	mu       sync.Mutex
	lessFunc func(i, j int) bool
	safe     bool
}

//==============属性==============

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

// Safe 并发安全模式,执行操作(增加,修改,删除)前会加锁
func (this *List) Safe(b ...bool) *List {
	this.safe = len(b) == 0 || b[0]
	return this
}

// Less 实现排序接口,比较元素
func (this *List) Less(i, j int) bool {
	return this.lessFunc == nil || this.lessFunc(i, j)
}

// Swap 实现排序接口,交换元素
func (this *List) Swap(i, j int) {
	this.list[i], this.list[j] = this.list[j], this.list[i]
}

// String 对应%s
func (this *List) String() string {
	return fmt.Sprintf("%v", this.list)
}

// Copy 复制
func (this *List) Copy() *List {
	return New(this.list...)
}

//==============查询==============

// Len 元素长度
func (this *List) Len() int {
	return len(this.list)
}

// Count 参照SQL的命名,获取总数量
func (this *List) Count() int64 {
	return int64(this.Len())
}

// Cap 总长
func (this *List) Cap() int {
	return cap(this.list)
}

// List 全部列表,引用类型
func (this *List) List() []interface{} {
	return this.list
}

// Find 参照SQL的名称,解析数据到ptr
func (this *List) Find(ptr interface{}) error {
	return conv.Unmarshal(this.list, ptr)
}

// Scan 参照其他的名称,解析数据到ptr
func (this *List) Scan(ptr interface{}) error {
	return conv.Unmarshal(this.list, ptr)
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

// First 获取第一个元素
func (this *List) First() interface{} {
	return this.MustGet(0)
}

// Last 获取最后一个元素
func (this *List) Last() interface{} {
	return this.MustGet(-1)
}

// GetVar 获取数据转成*conv.Var
func (this *List) GetVar(idx int) *conv.Var {
	return conv.New(this.MustGet(idx))
}

// GetAndDel 获取元素并删除
func (this *List) GetAndDel(idx int) (interface{}, bool) {
	if this.safe {
		this.mu.Lock()
		defer this.mu.Unlock()
	}

	if idx = this.getIdx(idx); idx >= 0 {
		val := this.list[idx]
		this.list = append(this.list[:idx], this.list[idx+1:]...)
		return val, true
	}
	return nil, false
}

//==============增加==============

func (this *List) append(v ...interface{}) *List {
	switch {
	case len(v) == 0:
	case len(v) == 1:
		switch x := v[0].(type) {
		case *List:
			this.list = append(this.list, x.List()...)
			this.list = append(this.list, x.List()...)
		default:
			this.list = append(this.list, conv.Interfaces(x)...)
		}
	default:
		for _, k := range v {
			switch x := k.(type) {
			case *List:
				this.list = append(this.list, x.List()...)
			default:
				this.list = append(this.list, x)
			}
		}
	}
	return this
}

// Append 从最后添加元素
func (this *List) Append(v ...interface{}) *List {
	if this.safe {
		this.mu.Lock()
		defer this.mu.Unlock()
	}
	return this.append(v...)
}

// Insert 中间插入元素
func (this *List) Insert(idx int, v ...interface{}) {
	if this.safe {
		this.mu.Lock()
		defer this.mu.Unlock()
	}

	if idx = this.getIdx(idx); idx >= 0 {
		after := this.list[idx:]
		cache := make([]interface{}, len(after))
		copy(cache, after)
		this.list = this.list[:idx]
		this.append(v...)
		this.list = append(this.list, cache...)
	}
}

//==============修改==============

// Set 替换元素,替换已有的元素
func (this *List) Set(idx int, v interface{}) *List {
	if idx = this.getIdx(idx); idx >= 0 {
		this.list[idx] = v
	}
	return this
}

// Replace 替换元素,替换已有的元素
func (this *List) Replace(idx int, v interface{}) *List {
	return this.Set(idx, v)
}

// MoveToLast 获取数据并移到队列最后
func (this *List) MoveToLast(idx int) interface{} {
	if this.safe {
		this.mu.Lock()
		defer this.mu.Unlock()
	}

	if idx = this.getIdx(idx); idx >= 0 {
		v := this.list[idx]
		this.list = append(this.list[:idx], this.list[idx+1:]...)
		this.list = append(this.list, v)
		return v
	}
	return nil
}

// Reverse 倒序
func (this *List) Reverse() *List {
	if this.safe {
		this.mu.Lock()
		defer this.mu.Unlock()
	}

	for i := 0; i < this.Len()/2; i++ {
		this.Swap(i, this.Len()-1-i)
	}
	return this
}

// Sort 排序
func (this *List) Sort(fn func(a, b interface{}) bool) *List {
	if this.safe {
		this.mu.Lock()
		defer this.mu.Unlock()
	}

	this.lessFunc = func(i, j int) bool {
		return fn(this.list[i], this.list[j])
	}
	sort.Sort(this)
	return this
}

// OrderBy 仿照SQL命名的排序
func (this *List) OrderBy(fn func(v interface{}) interface{}) *List {
	return this.Sort(func(a, b interface{}) bool {
		flied1 := fn(a)
		flied2 := fn(b)
		switch f1 := conv.GetNature(flied1).(type) {
		case int, int8, int16, int32,
			uint, uint8, uint16,
			uint32, uint64:
			return conv.Uint64(flied1) < conv.Uint64(flied2)
		case int64:
			return f1 < conv.Int64(flied2)
		case float32:
			return f1 < conv.Float32(flied2)
		case float64:
			return f1 < conv.Float64(flied2)
		case string:
			return f1 < conv.String(flied2)
		case bool:
			return f1
		case time.Time:
			return f1.UnixNano()-conv.Int64(flied2) > 0
		default:
			return conv.String(flied1) < conv.String(flied2)
		}
	})
}

//==============删除==============

// Where 筛选数据,参照SQL的命名
func (this *List) Where(fn func(i int, v interface{}) bool) *List {
	if this.Len() > 0 {
		if this.safe {
			this.mu.Lock()
			defer this.mu.Unlock()
		}

		cache := make([]interface{}, 0, this.Len())
		for i, v := range this.list {
			if fn(i, v) {
				cache = append(cache, v)
			}
		}
		this.list = cache
	}
	return this
}

func (this *List) WhereValue(fn func(v interface{}) bool) *List {
	return this.Where(func(i int, v interface{}) bool {
		return fn(v)
	})
}

// Del 移除元素
func (this *List) Del(idx ...int) *List {
	m := make(map[int]bool)
	for _, v := range idx {
		if i := this.getIdx(v); i > 0 {
			m[i] = true
		}
	}
	return this.Where(func(i int, v interface{}) bool {
		return !m[i]
	})
}

// Delete 移除元素
func (this *List) Delete(idx ...int) *List {
	return this.Del(idx...)
}

// Remove 移除元素
func (this *List) Remove(idx ...int) *List {
	return this.Del(idx...)
}

// RemoveNil 移除nil的元素
func (this *List) RemoveNil() *List {
	return this.Where(func(i int, v interface{}) bool { return v != nil })
}

// Cut 剪切,新值 , 同 list[start:end]
func (this *List) Cut(start, end int) *List {
	if this.Len() > 0 {
		if this.safe {
			this.mu.Lock()
			defer this.mu.Unlock()
		}

		start = this.getIdx(start)
		_end := this.getIdx(end)

		if end > this.Len() {
			_end = this.Len()
		}
		if start < 0 {
			start = 0
		}
		if start > _end || start < 0 {
			this.list = this.list[:0]
		} else {
			this.list = this.list[start:_end]
		}

	}
	return this
}

func (this *List) Limit(size int, offset ...int) *List {
	start := conv.DefaultInt(0, offset...)
	end := start + size
	return this.Cut(start, end)
}

// Clear 清除元素
func (this *List) Clear() *List {
	if this.safe {
		this.mu.Lock()
		defer this.mu.Unlock()
	}

	this.list = []interface{}{}
	return this
}
