package maps

type Chan struct {
	key       interface{}
	C         chan interface{}
	closeFunc func()
	closed    bool
}

func (this *Chan) Close() {
	this.closed = true
	defer func() { recover() }()
	if this.closeFunc != nil {
		this.closeFunc()
	} else {
		close(this.C)
	}
}

func newChan(key interface{}, cap ...uint) *Chan {
	c := &Chan{key: key}
	if len(cap) > 0 && cap[0] > 0 {
		c.C = make(chan interface{}, cap[0])
	} else {
		c.C = make(chan interface{})
	}
	return c
}
