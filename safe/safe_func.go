package safe

import (
	"fmt"
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

//========================================TryCatch========================================

type TryErr struct {
	err error
}

func newErr(err error) *TryErr {
	return &TryErr{err: err}
}

func (this *TryErr) Error() string {
	if this.err != nil {
		return this.err.Error()
	}
	return ""
}

func (this *TryErr) Catch(fn ...func(err error)) *TryErr {
	if this.err != nil {
		for _, v := range fn {
			v(this.err)
		}
	}
	return this
}

func (this *TryErr) Finally(fn func(err error)) {
	if fn != nil {
		fn(this.err)
	}
}

// Try 尝试运行,捕捉错误
func Try(fn func() error) (err *TryErr) {
	defer RecoverFunc(func(er error, stack string) {
		err = newErr(er)
	})
	return newErr(fn())
}

// TryCatch 其他语言的try catch
func TryCatch(fn func() error, catch ...func(err error)) *TryErr {
	return Try(fn).Catch(catch...)
}
