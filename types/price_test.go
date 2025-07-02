package types

import "testing"

func TestPrice_FloatUnit(t *testing.T) {
	p := Price(16)
	t.Log(p.FloatUnit())
	p = Price(230)
	t.Log(p.FloatUnit())
	p = Price(1000)
	t.Log(p.FloatUnit())
	p = Price(1_0000_000)
	t.Log(p.FloatUnit())
	p = Price(500_0000_000)
	t.Log(p.FloatUnit())
}
