package json

import (
	"github.com/json-iterator/go"
)

type Encoder struct{}

func (this *Encoder) Encode(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func (this *Encoder) Decode(data []byte, v interface{}) error {
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
