package g

import (
	"errors"
	"time"
)

const (
	Null        = ""
	TimeDefault = "2006-01-02 15:04:05"
	TimeDate    = "2006-01-02"
	TimeTime    = "15:04:05"
)

var (
	StartTime = time.Now()
)

var (
	ErrContext = errors.New("上下文关闭")
	ErrTimeout = errors.New("超时")
)
