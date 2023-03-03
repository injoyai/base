package crypt

const Padding = byte(58)

func DealLength(bs []byte, length int) []byte {
	beforeLength := len(bs)
	if beforeLength > length {
		bs = bs[:length]
	} else if beforeLength < length {
		for i := 0; i < length-beforeLength; i++ {
			bs = append(bs, Padding)
		}
	}
	return bs
}
