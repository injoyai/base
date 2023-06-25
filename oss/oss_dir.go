package oss

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	DefaultName = "injoy"
)

// ExecName 当前执行的程序名称
func ExecName() string {
	fullName, _ := os.Executable()
	return fullName
}

// ExecDir 当前执行的程序路径
func ExecDir() string {
	return filepath.Dir(ExecName())
}

// FuncName 当前执行的函数名称
func FuncName() string {
	_, filename, _, _ := runtime.Caller(1)
	return filename
}

// FuncDir 当前执行的函数路径
func FuncDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}

// UserDir 系统用户路径
func UserDir() string {
	dir, _ := os.UserHomeDir()
	return dir
}

// UserLocalDir 系统用户数据路径
func UserLocalDir(join ...string) string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(append([]string{dir, "AppData/Local"}, join...)...)
}

// UserStartupDir  程序自启目录
func UserStartupDir(join ...string) string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(append([]string{dir, "AppData/Roaming/Microsoft/Windows/Start Menu/Programs/Startup"}, join...)...)
}

// UserDefaultDir 默认系统用户数据子路径(个人使用)
func UserDefaultDir() string {
	return UserLocalDir(DefaultName)
}
