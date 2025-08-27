package types

import "github.com/injoyai/conv"

type Memory = Reuse[byte]

func NewMemory(max int) *Memory {
	return NewReuse[byte](max)
}

func NewReuse[T any](max int) *Reuse[T] {
	var zero T
	return &Reuse[T]{
		zero:  zero,
		ls:    make([]T, max),
		index: 0,
	}
}

type Reuse[T any] struct {
	zero  T
	ls    []T
	index int
}

func (this *Reuse[T]) Reset() {
	for i := range this.ls {
		this.ls[i] = this.zero
	}
	this.index = len(this.ls)
}

func (this *Reuse[T]) Take(start int, end ...int) []T {

	last := this.index

	i1 := start
	if i1 < 0 {
		i1 = last + i1
	}
	if i1 < 0 || i1 > last {
		return nil
	}

	i2 := conv.Default(last, end...)
	if i2 < 0 {
		i2 = last + i2
	}
	if i2 > last {
		i2 = last
	}
	if i2 < 0 || i2 <= i1 {
		return nil
	}

	return this.ls[i1:i2]
}

func (this *Reuse[T]) getIdx(idx int) int {
	if idx >= 0 && idx < this.index {
		return idx
	}
	if idx < 0 && -idx <= this.index {
		return this.index + idx
	}
	return -1
}

func (this *Reuse[T]) Set(idx int, v T) bool {
	if idx = this.getIdx(idx); idx >= 0 {
		this.ls[idx] = v
		return true
	}
	return false
}

func (this *Reuse[T]) Cut(start int, end ...int) []T {
	return this.Take(start, end...)
}

func (this *Reuse[T]) CopyFrom(src []T, start int, end ...int) {
	copy(this.Take(start, end...), src)
}
