package maps

type Write func(p []byte) (int, error)

func (this Write) Write(p []byte) (int, error) {
	return this(p)
}
