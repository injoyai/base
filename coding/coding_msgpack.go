package coding

import (
	"bytes"
	"github.com/ugorji/go/codec"
)

type Msgpack struct{}

func (this *Msgpack) Marshal(v any) ([]byte, error) {
	return MsgpackMarshal(v)
}

func (this *Msgpack) Unmarshal(data []byte, v any) error {
	return MsgpackUnmarshal(data, v)
}

func MsgpackMarshal(v any) ([]byte, error) {
	var mh codec.MsgpackHandle
	w := bytes.NewBuffer(nil)
	err := codec.NewEncoder(w, &mh).Encode(v)
	return w.Bytes(), err
}

func MsgpackUnmarshal(data []byte, v any) error {
	var mh codec.MsgpackHandle
	r := bytes.NewBuffer(data)
	err := codec.NewDecoder(r, &mh).Decode(v)
	return err
}
