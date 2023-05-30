package g

import (
	"errors"
	"fmt"
	"github.com/injoyai/base/bytes/crypt/md5"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/base/maps/wait"
	"github.com/injoyai/base/oss"
	"github.com/injoyai/base/safe"
	"github.com/injoyai/conv"
	uuid "github.com/satori/go.uuid"
	"runtime/debug"
	"time"
)

//========================================Chan========================================

// Range 仿python的range 参数1-3个
// 当1个参数, 例 Range(5) 输出 0,1,2,3,4
// 当2个参数, 例 Range(1,5) 输出 1,2,3,4
// 当3个参数, 例 Range(0,5,2) 输出 0,2,4
func Range(n int, m ...int) <-chan int {
	return chans.Range(n, m...)
}

// Count 遍历
// @num,数量,-1为死循环
// @interval,间隔
func Count(num int, interval ...time.Duration) <-chan int {
	return chans.Count(num, interval...)
}

// Interval 间隔触发
func Interval(interval time.Duration, nums ...int) <-chan int {
	num := conv.GetDefaultInt(-1, nums...)
	return chans.Count(num, interval)
}

//========================================Crypt========================================

// MD5 加密,返回hex的32位小写
func MD5(s string) string { return md5.Encrypt(s) }

// HmacMD5 加密,返回hex的32位小写
func HmacMD5(s, key string) string { return md5.Hmac(s, key) }

//========================================Runtime========================================

// Recover 错误捕捉
func Recover(err *error, stack ...bool) {
	if er := recover(); er != nil {
		if err != nil {
			if len(stack) > 0 && stack[0] {
				*err = errors.New(fmt.Sprintln(er) + string(debug.Stack()))
			} else {
				*err = errors.New(fmt.Sprintln(er))
			}
		}
	}
}

// Try 尝试运行,捕捉错误 其他语言的try catch
func Try(fn func() error, catch ...func(err error)) *safe.TryErr {
	return safe.Try(fn).Catch(catch...)
}

// Retry 重试,默认3次
func Retry(fn func() error, nums ...int) (err error) {
	num := conv.GetDefaultInt(3, nums...)
	for i := 0; i < num; i++ {
		if err = Try(fn); err == nil {
			return
		}
	}
	return
}

// PanicErr 如果是错误则panic
func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

//========================================Wait========================================

// Wait 等待
func Wait(key string) (interface{}, error) { return wait.Wait(key) }

// Done 结束等待
func Done(key string, value interface{}) { wait.Done(key, value) }

//========================================OS========================================

// Input 监听用户输入
func Input(hint ...interface{}) string { return oss.Input(hint...) }

// ListenExit 监听退出信号
func ListenExit(fn ...func()) { oss.ListenExit(fn...) }

// ExecName 当前执行的程序名称
func ExecName() string { return oss.ExecName() }

// ExecDir 当前执行的程序路径
func ExecDir() string { return oss.ExecDir() }

// FuncName 当前执行的函数名称
func FuncName() string { return oss.FuncName() }

// FuncDir 当前执行的函数路径
func FuncDir() string { return oss.FuncDir() }

// UserDir 系统用户路径
func UserDir() string { return oss.UserDir() }

// UserDataDir 系统用户数据路径
func UserDataDir() string { return oss.UserDataDir() }

// UserDefaultDir 默认系统用户数据子路径(个人使用)
func UserDefaultDir() string { return oss.UserDefaultDir() }

//========================================Third========================================

// UUID uuid
func UUID() string { return uuid.NewV4().String() }
