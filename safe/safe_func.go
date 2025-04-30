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
func RecoverFunc(fn func(err any, stack []byte)) {
	if fn != nil {
		if e := recover(); e != nil {
			fn(e, debug.Stack())
		}
	}
}

// Try 尝试运行,捕捉错误
func Try(fn func() error, catch ...func(err error)) (err error) {
	defer RecoverFunc(func(e any, stack []byte) {
		err = fmt.Errorf("%v", e)
		for _, v := range catch {
			v(err)
		}
	})
	return fn()
}

func Retry(fn func() error, retry ...int) (err error) {
	num := conv.Default[int](3, retry...)
	for i := 0; i < num; i++ {
		if err = Try(fn); err == nil {
			return
		}
	}
	return
}
