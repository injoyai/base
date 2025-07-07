package types

import (
	"testing"
)

type A struct {
	Name string
	Age  int
}

func (this *A) Default() {
	this.Name = "张三"
	this.Age = 18
}

func WithName(name string) Option[A] {
	return func(a *A) {
		a.Name = name
	}
}

func WithAge(age int) Option[A] {
	return func(a *A) {
		a.Age = age
	}
}

func TestNewOption(t *testing.T) {
	x := NewOption[A](
		WithName("李四"),
		WithAge(15),
	)
	t.Log(x)
}
