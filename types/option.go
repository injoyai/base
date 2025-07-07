package types

type Option[T any] func(*T)

type Defaulter[T any] interface {
	*T //表示限制这种类型
	Default()
}

func NewOption[T any, P Defaulter[T]](ops ...Option[T]) *T {
	t := new(T)
	// 初始化 等同于 P(t).Default()
	P.Default(t)
	for _, op := range ops {
		op(t)
	}
	return t
}
