package encoding

type Decoder interface {
	Decode([]byte, interface{})
}

type Encoder interface {
	Encode(interface{}) ([]byte, error)
}
