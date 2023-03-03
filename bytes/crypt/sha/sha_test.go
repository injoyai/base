package sha

import (
	"testing"
)

func TestEncrypt256(t *testing.T) {
	t.Log(Encrypt256("test"))
	t.Log(Encrypt1("test"))
	t.Log(Encrypt512("test"))
}
