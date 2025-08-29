package maps

import (
	"sync"
	"sync/atomic"
)

type Bit interface {
	Get(key uint64) bool
	Set(key uint64, value bool)
}

// NewBit 创建无锁 Bit map
func NewBit() Bit {
	return &bit{}
}

type bit struct {
	m sync.Map // key: group(uint64), value: *uint64
}

// Set 设置 key 的值
func (b *bit) Set(key uint64, value bool) {
	group := key / 64
	offset := key % 64

	ptrAny, _ := b.m.LoadOrStore(group, new(uint64))
	ptr := ptrAny.(*uint64)

	for {
		old := atomic.LoadUint64(ptr)
		var _new uint64
		if value {
			_new = old | (1 << offset)
		} else {
			_new = old & ^(uint64(1) << offset)
		}
		if atomic.CompareAndSwapUint64(ptr, old, _new) {
			break
		}
	}
}

// Get 获取 key 的值
func (b *bit) Get(key uint64) bool {
	group := key / 64
	offset := key % 64

	ptrAny, ok := b.m.Load(group)
	if !ok {
		return false
	}
	ptr := ptrAny.(*uint64)
	return atomic.LoadUint64(ptr)&(1<<offset) != 0
}

// Del 删除 key
func (b *bit) Del(key uint64) {
	offset := key % 64
	group := key / 64

	ptrAny, ok := b.m.Load(group)
	if !ok {
		return
	}
	ptr := ptrAny.(*uint64)

	for {
		old := atomic.LoadUint64(ptr)
		newVal := old & ^(uint64(1) << offset)
		if old == newVal {
			return
		}
		if atomic.CompareAndSwapUint64(ptr, old, newVal) {
			if newVal == 0 {
				// 所有位都清零了，直接从 map 删除
				b.m.Delete(group)
			}
			return
		}
	}
}
