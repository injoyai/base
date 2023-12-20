package chans

import (
	"errors"
	"time"
)

// MustWriter chan []byte 实现 io.Writer,必须等到写入成功为止
type MustWriter chan []byte

func (this MustWriter) Write(p []byte) (int, error) {
	this <- p
	return len(p), nil
}

func (this MustWriter) WriteTimeout(p []byte, timeout time.Duration) (int, error) {
	select {
	case this <- p:
		return len(p), nil
	case <-time.After(timeout):
		return 0, errors.New("写入超时")
	}
}

// TryWriter chan []byte 实现 io.Writer,尝试写入,不管是否成功
type TryWriter chan []byte

func (this TryWriter) Write(p []byte) (int, error) {
	select {
	case this <- p:
		return len(p), nil
	default:
		return 0, nil
	}
}
