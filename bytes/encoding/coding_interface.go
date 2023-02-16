package coding

type Coding interface {
	Decoder
	Encoder
}

type Decoder interface {
	Decode([]byte, interface{})
}

type Encoder interface {
	Encode(interface{}) ([]byte, error)
}

type IMarshal interface {
	Marshal
	Unmarshal
}

type Marshal interface {
	Marshal(interface{}) ([]byte, error)
}

type Unmarshal interface {
	Unmarshal([]byte, interface{}) error
}
