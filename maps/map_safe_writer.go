package maps

type WriteFunc func(p []byte) (int, error)

func (this WriteFunc) Write(p []byte) (int, error) {
	return this(p)
}
