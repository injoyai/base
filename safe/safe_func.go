package safe

import (
	"fmt"
	"github.com/injoyai/conv"
	"runtime/debug"
)

// Recover 错误捕捉
func Recover(err *error, stack ...bool) {
	if er := recover(); er != nil {
		if err != nil {
			if len(stack) > 0 && stack[0] {
				*err = fmt.Errorf("%v\n%v", er, string(debug.Stack()))
			} else {
				*err = fmt.Errorf("%v", er)
			}
		}
	}
}

// RecoverFunc 捕捉错误并执行函数
func RecoverFunc(fn func(err error, stack string)) {
	if er := recover(); er != nil {
		if fn != nil {
			fn(fmt.Errorf("%v", er), string(debug.Stack()))
		}
	}
}

// Try 尝试运行,捕捉错误
func Try(fn func() error, catch ...func(err error)) (err error) {
	defer RecoverFunc(func(er error, stack string) {
		if er != nil {
			err = er
			for _, v := range catch {
				v(err)
			}
		}
	})
	return fn()
}

func Retry(fn func() error, nums ...int) (err error) {
	num := conv.GetDefaultInt(3, nums...)
	for i := 0; i < num; i++ {
		if err = Try(fn); err == nil {
			return
		}
	}
	return
}
