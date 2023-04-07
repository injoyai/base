package str

import "testing"

func TestMustSplitN(t *testing.T) {
	testStr, testSep := "a,b,c", ","
	t.Log(MustSplitN(testStr, testSep, 3))
	t.Log(MustSplitN(testStr, testSep, 2))
	t.Log(len(MustSplitN(testStr, testSep, 2)))
	t.Log(MustSplitN(testStr, testSep, 1))
	t.Log(MustSplitN(testStr, testSep, 0))
	t.Log(MustSplitN(testStr, testSep, -1))
}
