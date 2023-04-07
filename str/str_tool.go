package str

import (
	"strings"
	"unsafe"
)

// Pointer 获取字符串指针,一般用于节约内存
func Pointer(s string) *string {
	return &s
}

// Bytes 字符串转字节,不同类型共用同一地址
func Bytes(s *string) *[]byte {
	return (*[]byte)(unsafe.Pointer(s))
}

// Reverse 把字符串顺序逆转过来
func Reverse(str string) string {
	str1 := []rune(str)
	str2 := []rune(str)
	n := len(str1)
	for i := 0; i < n; i++ {
		str1[i] = str2[n-i-1]
	}
	return string(str1)
}

// HasRepeat 是否有重复
func HasRepeat(str string) bool {
	m := make(map[int32]bool)
	for _, v := range str {
		if m[v] {
			return true
		}
		m[v] = true
	}
	return false
}

// CropFirst 裁剪,剪短
// 例: "0123456789", "2" >>> "23456789"
// @str,被裁剪的字符串
// @crop,裁剪的字符串
// @retain,是否保留裁剪字符串,默认保留
func CropFirst(str, crop string, retain ...bool) string {
	n := strings.Index(str, crop)
	if n >= 0 {
		if len(retain) > 0 && !retain[0] {
			return str[n+len(crop):]
		}
		return str[n:]
	}
	return str
}

// CropLast 裁剪,剪短
// 例: "0123456789", "2" >>> "012"
// @str,被裁剪的字符串
// @crop,裁剪的字符串
// @retain,是否保留裁剪字符串,默认保留 "0123456789", "2" >>> "012" "0123456789", "2" >>> "01"
func CropLast(str, crop string, retain ...bool) string {
	n := strings.LastIndex(str, crop)
	if n+len(crop) <= len(str) && n >= 0 {
		if len(retain) > 0 && !retain[0] {
			return str[:n]
		}
		return str[:n+len(crop)]
	}
	return str
}

// GetLine 按行分割 获取第几行数据
func GetLine(s string, idx int) string {
	return GetSplit(s, "\n", idx)
}

// GetSplit 分割字符串,并取对应下标的字符串,默认""
func GetSplit(s, sep string, idx int) string {
	list := strings.SplitN(s, sep, idx)
	if len(list) > idx {
		return list[idx]
	}
	return ""
}

// FindCommon 获取多个列表的公共部分
// stringSlice1 = []string{"aaa", "bbb", "ccc"}
// stringSlice2 = []string{"bbb", "ccc", "ddd"}
// stringSlice3 = []string{"aaa", "bbb", "ccc", "eee"}
// FindCommon(stringSlice1, stringSlice2, stringSlice3) shall return a string slice with {"bbb", "ccc"}
func FindCommon(stringsArray ...[]string) (commonStrings []string) {
	if len(stringsArray) == 0 {
		return nil
	} else if len(stringsArray) == 1 {
		return stringsArray[0]
	}
	commonStringsSet := make([]string, 0)
	hash := make(map[string]bool)
	for _, s := range stringsArray[0] {
		hash[s] = true
	}
	for _, s := range stringsArray[1] {
		if _, ok := hash[s]; ok {
			commonStringsSet = append(commonStringsSet, s)
		}
	}
	stringsArray = append([][]string{commonStringsSet}, stringsArray[2:]...)
	return FindCommon(stringsArray...)
}

// MustSplitN 分割字符串,并获取指定下标字符,不存在返回""
func MustSplitN(s, sep string, idx int) string {
	list := strings.SplitN(s, sep, idx+2) //至少需要分割2次
	if len(list) > idx && idx >= 0 {
		return list[idx]
	}
	return ""
}
