package chans

type LimitGo interface {
	Go(f func())
	Wait()
}

func NewLimitGo(limit int) LimitGo {
	return &limitGo{limit: NewWaitLimit(limit)}
}

type limitGo struct {
	limit WaitLimit
}

func (this *limitGo) Go(f func()) {

	this.limit.Add()
	go func() {
		defer this.limit.Done()
		f()
	}()
}

func (this *limitGo) Wait() {
	this.limit.Wait()
}
