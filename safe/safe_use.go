package safe

import "sync/atomic"

type Use struct {
	used uint32
}

// Use 尝试占用,返回是否占用成功,失败标识被其他进程占用
func (this *Use) Use() bool {
	return atomic.CompareAndSwapUint32(&this.used, 0, 1)
}

// UnUse 取消占用
func (this *Use) UnUse() bool {
	return atomic.CompareAndSwapUint32(&this.used, 1, 0)
}

// Used 是否被占用
func (this *Use) Used() bool {
	return atomic.LoadUint32(&this.used) == 1
}
