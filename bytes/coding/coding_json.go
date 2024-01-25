package coding

import jsoniter "github.com/json-iterator/go"

type Json struct{}

func (this *Json) Marshal(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func (this *Json) Unmarshal(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}

var (
	JsonMarshal   = jsoniter.Marshal
	JsonUnmarshal = jsoniter.Unmarshal
)
