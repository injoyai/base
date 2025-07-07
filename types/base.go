package types

// Err 错误,好处是能定义在const
type Err string

func (this Err) Error() string { return string(this) }

func (this Err) Is(target error) bool {
	return target != nil && this.Error() == target.Error()
}

type Debugger bool

func (this *Debugger) Debug(b ...bool) {
	*this = Debugger(len(b) == 0 || b[0])
}

func (this *Debugger) Debugged() bool {
	return bool(*this)
}
