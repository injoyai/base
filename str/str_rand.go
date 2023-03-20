package str

import (
	"math/rand"
	"strings"
	"time"
)

var (
	r              *rand.Rand
	defaultRandStr = "abcdefghigklmnopqrstuvwxyzABCDEFGHIGKLMNOPQRSTUVWXYZ0123456789"
)

// Rand 随机字符串
func Rand(length int, str ...string) string {
	if r == nil {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	rs := defaultRandStr
	if len(str) > 0 {
		rs = strings.Join(str, "")
	}
	var s []byte
	for i := 0; i < length; i++ {
		n := r.Intn(len(rs))
		s = append(s, rs[n])
	}
	return string(s)
}
