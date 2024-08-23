package safe

import "sync/atomic"

func NewUse() *Use { return &Use{} }

type Use struct {
	used uint32
	//成功占用执行的事件,例初始化操作
	OnUse   func()
	OnUnUse func()
}

// Use 尝试占用,返回是否占用成功,失败标识被其他进程占用
func (this *Use) Use() bool {
	if !atomic.CompareAndSwapUint32(&this.used, 0, 1) {
		return false
	}
	if this.OnUse != nil {
		this.OnUse()
	}
	return true
}

// UnUse 取消占用
func (this *Use) UnUse() bool {
	if !atomic.CompareAndSwapUint32(&this.used, 1, 0) {
		return false
	}
	if this.OnUnUse != nil {
		this.OnUnUse()
	}
	return true
}

// Used 是否被占用
func (this *Use) Used() bool {
	return atomic.LoadUint32(&this.used) == 1
}
