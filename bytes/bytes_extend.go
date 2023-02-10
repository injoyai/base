package bytes

import "bytes"

type (
	Buffer = bytes.Buffer
	Reader = bytes.Reader
)

var (
	NewBuffer       = bytes.NewBuffer
	NewBufferString = bytes.NewBufferString
	NewReader       = bytes.NewReader
	Replace         = bytes.Replace
	ReplaceAll      = bytes.ReplaceAll
	Fields          = bytes.Fields
	Count           = bytes.Count
	Map             = bytes.Map
	Join            = bytes.Join
	Compare         = bytes.Compare
	Split           = bytes.Split
	Contains        = bytes.Contains
	Index           = bytes.Index
	ToLower         = bytes.ToLower
	ToUpper         = bytes.ToUpper
)
