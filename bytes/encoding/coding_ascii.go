package coding

type ASCII struct{}

func (this ASCII) Decode(s string) ([]byte, error) {
	return []byte(s), nil
}

func (this ASCII) Encode(bs []byte) string {
	return string(bs)
}

// DecodeASCII []byte()
func DecodeASCII(s string) ([]byte, error) {
	return []byte(s), nil
}

// EncodeASCII string()
func EncodeASCII(bs []byte) string {
	return string(bs)
}
