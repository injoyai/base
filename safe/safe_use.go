package safe

import "sync/atomic"

func NewUse(f ...func()) *Use {
	return &Use{
		OnUser: func() {
			for _, v := range f {
				v()
			}
		},
	}
}

type Use struct {
	used uint32
	//成功占用执行的事件,例初始化操作
	OnUser func()
}

// Use 尝试占用,返回是否占用成功,失败标识被其他进程占用
func (this *Use) Use() bool {
	if !atomic.CompareAndSwapUint32(&this.used, 0, 1) {
		return false
	}
	if this.OnUser != nil {
		this.OnUser()
	}
	return true
}

// UnUse 取消占用
func (this *Use) UnUse() bool {
	return atomic.CompareAndSwapUint32(&this.used, 1, 0)
}

// Used 是否被占用
func (this *Use) Used() bool {
	return atomic.LoadUint32(&this.used) == 1
}
