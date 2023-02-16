package json

import (
	"github.com/json-iterator/go"
)

type Entity struct{}

func (this *Entity) Marshal(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func (this *Entity) Unmarshal(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}

var (
	Marshal         = jsoniter.Marshal
	MarshalIndent   = jsoniter.MarshalIndent
	MarshalToString = jsoniter.MarshalToString
	Unmarshal       = jsoniter.Unmarshal
	NewDecoder      = jsoniter.NewDecoder
	NewEncoder      = jsoniter.NewEncoder
	Valid           = jsoniter.Valid
)
