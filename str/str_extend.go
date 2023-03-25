package str

import "strings"

var (
	HasPrefix  = strings.HasPrefix
	HasSuffix  = strings.HasSuffix
	Contains   = strings.Contains
	Split      = strings.Split
	ReplaceAll = strings.ReplaceAll
	Count      = strings.Count
	Join       = strings.Join
	NewReader  = strings.NewReader
	TrimPrefix = strings.TrimPrefix
	TrimSuffix = strings.TrimSuffix
	TrimSpace  = strings.TrimSpace
	Index      = strings.Index
	Title      = strings.Title
	ToUpper    = strings.ToUpper
	ToLower    = strings.ToLower

	// IsBegin 是否开始
	// 例: "0123456" , "012" , "0" 满足条件
	IsBegin = strings.HasPrefix

	// IsEnd 是否结尾.
	// 例: "0123456" , "456" , "3456" 满足条件
	IsEnd = strings.HasSuffix
)
