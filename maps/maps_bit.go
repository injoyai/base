package maps

type Bit interface {
	Get(key uint64) bool
	Set(key uint64, value bool)
}

func NewBit() Bit {
	return &bit{NewSafe[uint64, *uint64]()}
}

type bit struct {
	*Safe[uint64, *uint64]
}

func (this *bit) Set(key uint64, value bool) {
	offset := key % 64
	group := key / 64
	base := uint64(0)
	if value {
		base = 1 << offset
	}
	v, ok := this.Safe.Get(group)
	if !ok {
		this.Safe.Set(group, &base)
		return
	}
	if ((*v)>>offset)%2 == 0 {
		*v += base
	} else if !value {
		*v -= base
	}
}

func (this *bit) Get(key uint64) bool {
	offset := key % 64
	group := key / 64
	val, ok := this.Safe.Get(group)
	if ok {
		return (*(val)>>offset)%2 == 1
	}
	return false
}
