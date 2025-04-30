package coding

type Coding interface {
	Decoder
	Encoder
}

type Decoder interface {
	Decode([]byte) ([]byte, error)
}

type Encoder interface {
	Encode([]byte) string
}

type IMarshal interface {
	Marshal
	Unmarshal
}

type Marshal interface {
	Marshal(any) ([]byte, error)
}

type Unmarshal interface {
	Unmarshal([]byte, any) error
}
