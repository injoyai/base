package str

import (
	"strings"
	"testing"
)

func TestMustSplitN(t *testing.T) {
	testStr, testSep := "a,b,c", ","
	t.Log(strings.SplitN(testStr, testSep, 10))
	t.Log(MustSplitN(testStr, testSep, 3))
	t.Log(MustSplitN(testStr, testSep, 2))
	t.Log(len(MustSplitN(testStr, testSep, 2)))
	t.Log(MustSplitN(testStr, testSep, 1))
	t.Log(MustSplitN(testStr, testSep, 0))
	t.Log(MustSplitN(testStr, testSep, -1))
}
