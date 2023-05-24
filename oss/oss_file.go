package oss

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/injoyai/conv"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	defaultPerm = 0777
)

// Read 读取文件
func Read(filename string) ([]byte, error) {
	return ReadBytes(filename)
}

// ReadBytes 读取文件内容
func ReadBytes(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// ReadString 读取文件
func ReadString(filename string) (string, error) {
	bs, err := ReadBytes(filename)
	return string(bs), err
}

// ReadHEX 读取文件
func ReadHEX(filename string) (string, error) {
	bs, err := ReadBytes(filename)
	return hex.EncodeToString(bs), err
}

// ReadBase64 读取文件
func ReadBase64(filename string) (string, error) {
	bs, err := ReadBytes(filename)
	return base64.StdEncoding.EncodeToString(bs), err
}

// NewDir 新建文件夹
// @path,路径
func NewDir(path string) error {
	return os.MkdirAll(path, defaultPerm)
}

// NewFile 新建文件
func NewFile(filename string) (io.ReadWriteCloser, error) {
	return os.Create(filename)
}

// OpenFile 打开文件
func OpenFile(filename string) (io.ReadWriteCloser, error) {
	return os.Open(filename)
}

// New 新建文件,会覆盖
func New(filename string, v ...interface{}) error {
	dir := filepath.Dir(filename)
	name := filepath.Base(filename)
	if err := os.MkdirAll(dir, defaultPerm); err != nil {
		return err
	}
	if len(name) == 0 {
		return nil
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if len(v) == 0 {
		return nil
	}
	data := []byte(nil)
	for _, k := range v {
		data = append(data, conv.Bytes(k)...)
	}
	_, err = f.Write(data)
	return err
}

// ReadDirFunc 遍历目录
func ReadDirFunc(dir string, fn func(info os.FileInfo) error) error {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, info := range fileInfos {
		if err = fn(info); err != nil {
			return err
		}
	}
	return nil
}

// OpenFunc 打开文件,并执行函数
func OpenFunc(filename string, fn func(f *os.File) error) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return fn(f)
}

// OpenWithWriteTo 打开文件,并写入到io.Writer
func OpenWithWriteTo(filename string, writer io.Writer) error {
	return OpenFunc(filename, func(f *os.File) error {
		_, err := io.Copy(writer, f)
		return err
	})
}

// OpenWithWriteFrom 打开文件,并从io.Reader写入
func OpenWithWriteFrom(filename string, reader io.Reader) error {
	return OpenFunc(filename, func(f *os.File) error {
		_, err := io.Copy(f, reader)
		return err
	})
}
