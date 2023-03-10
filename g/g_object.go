package g

import (
	"context"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/conv"
	"github.com/injoyai/conv/cfg"
	"math/rand"
	"time"
)

//========================================Rand========================================

var (
	r  *rand.Rand
	rs string
)

// Rand 随机数
func Rand() *rand.Rand {
	if r == nil {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
		rs = "abcdefghigklmnopqrstuvwxyzABCDEFGHIGKLMNOPQRSTUVWXYZ0123456789"
	}
	return r
}

// RandString 随机字符串
func RandString(length int, str ...string) string {
	r := Rand()
	xs := conv.GetDefaultString(rs, str...)
	var s []byte
	for i := 0; i < length; i++ {
		n := r.Intn(len(xs))
		s = append(s, xs[n])
	}
	return string(s)
}

//========================================Context========================================

// Ctx context.Background
func Ctx() context.Context { return context.Background() }

// WithCancel context.WithCancel
func WithCancel(ctx ...context.Context) (context.Context, context.CancelFunc) {
	var c context.Context
	if len(ctx) > 0 && ctx[0] != nil {
		c = ctx[0]
	} else {
		c = context.Background()
	}
	return context.WithCancel(c)
}

// WithTimeout context.WithTimeout
func WithTimeout(timeout time.Duration, ctx ...context.Context) (context.Context, context.CancelFunc) {
	var c context.Context
	if len(ctx) > 0 && ctx[0] != nil {
		c = ctx[0]
	} else {
		c = context.Background()
	}
	return context.WithTimeout(c, timeout)
}

//========================================Time========================================

// Now 当前时间
func Now() time.Time { return time.Now() }

// Date 当前日期
func Date() (year int, month time.Month, day int) { return time.Now().Date() }

// Unix 当前时间戳
func Unix() int64 { return time.Now().Unix() }

// UnixNano 当前纳秒
func UnixNano() int64 { return time.Now().UnixNano() }

// Year 年[1970-]
func Year() int { return time.Now().Year() }

// Month 月[1-12]
func Month() int { return int(time.Now().Month()) }

// Day 日[1-31]
func Day() int { return time.Now().Day() }

// Hour 时[0-23]
func Hour() int { return time.Now().Hour() }

// Minute 分[0-59]
func Minute() int { return time.Now().Minute() }

// Second 秒[0-59]
func Second() int { return time.Now().Second() }

//========================================Other========================================

// Cfg 读取配置文件
func Cfg(paths ...string) *cfg.Entity { return cfg.New(paths...) }

// Chan chans.NewEntity
func Chan(num int, cap ...int) *chans.Entity { return chans.NewEntity(num, cap...) }
