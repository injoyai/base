package ios

import "io"

type (
	// Closed 是否已关闭
	Closed interface {
		Closed() bool
	}

	Closer interface {
		io.Closer
		Closed
	}

	AReader interface {
		ReadAck() (Acker, error)
	}

	AReadCloser interface {
		AReader
		io.Closer
	}

	AReadWriter interface {
		AReader
		io.Writer
	}

	AReadWriteCloser interface {
		AReader
		io.Writer
		io.Closer
	}

	MReader interface {
		ReadMessage() ([]byte, error)
	}

	MReadWriter interface {
		MReader
		io.Writer
	}

	MReadCloser interface {
		MReader
		io.Closer
	}

	MReadWriteCloser interface {
		MReader
		io.Writer
		io.Closer
	}
)

// Acker 兼容MQ等需要确认的场景
type Acker interface {
	Payload() []byte
	Ack() error
}

//=================================Func=================================

// ReadFunc 读取函数
type ReadFunc func(p []byte) (int, error)

func (this ReadFunc) Read(p []byte) (int, error) { return this(p) }

type AReadFunc func() (Acker, error)

func (this AReadFunc) ReadAck() (Acker, error) { return this() }

type MReadFunc func() ([]byte, error)

func (this MReadFunc) ReadAck() (Acker, error) {
	bs, err := this()
	return Ack(bs), err
}

func (this MReadFunc) ReadMessage() ([]byte, error) { return this() }

// WriteFunc 写入函数
type WriteFunc func(p []byte) (int, error)

func (this WriteFunc) Write(p []byte) (int, error) { return this(p) }

// CloseFunc 关闭函数
type CloseFunc func() error

func (this CloseFunc) Close() error { return this() }

type Ack []byte

func (this Ack) Ack() error { return nil }

func (this Ack) Payload() []byte { return this }
