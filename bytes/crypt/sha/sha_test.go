package sha

import (
	"testing"
)

func TestEncrypt256(t *testing.T) {
	t.Log(Encrypt256([]byte("test")))
	t.Log(Encrypt1([]byte("test")))
	t.Log(Encrypt512([]byte("test")))
}
