package g

import (
	"errors"
	"github.com/injoyai/base/bytes"
	"github.com/injoyai/conv"
	json "github.com/json-iterator/go"
)

type (
	Var   = conv.Var
	DMap  = conv.Map
	Any   = interface{}
	List  []interface{}
	Bytes = bytes.Entity
)

type KV struct {
	K string      `json:"key"`
	V interface{} `json:"value"`
	L string      `json:"label"`
}

//========================================Type========================================

const (
	String Type = "string"
	Bool   Type = "bool"
	Int    Type = "int"
	Float  Type = "float"
	Array  Type = "array"
	Object Type = "object"
)

// Type 数据类型
type Type string

func (this Type) String() string {
	switch this {
	case Float:
		return "浮点"
	case Int:
		return "整数"
	case String:
		return "字符"
	case Bool:
		return "布尔"
	case Array:
		return "数组"
	case Object:
		return "对象"
	}
	return "未知"
}

func (this Type) Value(v interface{}) interface{} {
	switch this {
	case String:
		return conv.String(v)
	case Bool:
		return conv.Bool(v)
	case Int:
		return conv.Int(v)
	case Float:
		return conv.Float64(v)
	}
	return v
}

func (this *Type) Check() error {
	switch *this {
	case String, Bool, Int, Float, Array, Object:
	case "":
		*this = String
	default:
		return errors.New("未知数据类型")
	}
	return nil
}

//========================================Map========================================

type Map map[string]interface{}

// Struct json Marshal
func (this Map) Struct(ptr interface{}) error {
	bs, err := json.Marshal(this)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, ptr)
}

// Json map转json
func (this Map) Json() string {
	return conv.String(this)
}

// Bytes map转字节
func (this Map) Bytes() []byte {
	return conv.Bytes(this)
}

// GetVar 实现conv.Extend接口
func (this Map) GetVar(key string) *conv.Var {
	return conv.New(this[key])
}

// Merge 合并2个map
func (this Map) Merge(m Map) Map {
	for key, val := range m {
		this[key] = val
	}
	return this
}
