package oss

import (
	"os"
)

// Exists 是否存在
func Exists(name string) bool {
	stat, err := os.Stat(name)
	return stat != nil && !os.IsNotExist(err)
}

// IsDir 是否是文件夹
func IsDir(name string) bool {
	s, err := os.Stat(name)
	return err == nil && s.IsDir()
}

// IsFile 是否是文件
func IsFile(name string) bool {
	s, err := os.Stat(name)
	return err == nil && !s.IsDir()
}
