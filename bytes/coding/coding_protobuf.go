package coding

import "google.golang.org/protobuf/proto"

type Proto struct{}

func (this *Proto) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (this *Proto) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}

var (
	ProtoMarshal   = proto.Marshal
	ProtoUnmarshal = proto.Unmarshal
)
