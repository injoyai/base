package keep

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	k := 1
	x := New[int](time.Second * 1)
	t.Log(x.Keeping(k))
	x.Keep(k)
	t.Log(x.Keeping(k))

	<-time.After(time.Second * 2)
	x.Keep(k)
	t.Log(x.Keeping(k))
	t.Log(x.Keeping(k, 2))
	x.Keep(k)
	t.Log(x.Keeping(k, 2))

}
