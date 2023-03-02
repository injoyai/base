package json

import (
	"encoding/json"

	"testing"
)

var TestStr = `{"a":1,"b":"2","c":true,"longLongLong":"longLongLong"}`

type Test struct {
	A            int    `json:"a"`
	B            string `json:"b"`
	C            bool   `json:"c"`
	LongLongLong string `json:"longLongLong"`
}

func TestEncoder_Encode(t *testing.T) {
	m := Test{} //make(map[string]interface{})
	for i := 0; i < 100000; i++ {
		Unmarshal([]byte(TestStr), &m)
	}
}

func TestEncoder_Encode2(t *testing.T) {
	m := Test{} //make(map[string]interface{})
	for i := 0; i < 100000; i++ {
		json.Unmarshal([]byte(TestStr), &m)
	}
}
