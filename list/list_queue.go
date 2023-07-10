package list

type Queue interface {
	Len() int                       //队列长度
	Cap() int                       //cap
	Get() (interface{}, bool)       //获取
	Ack()                           //确认
	GetAndAck() (interface{}, bool) //获取并确认
}

func NewQueue() Queue {
	return &queue{New()}
}

var _ Queue = new(queue)

type queue struct {
	*Entity
}

func (this *queue) Get() (interface{}, bool) {
	return this.Entity.Get(0)
}

func (this *queue) Ack() {
	this.Entity.Del(0)
}

func (this *queue) GetAndAck() (val interface{}, ok bool) {
	val, ok = this.Get()
	this.Ack()
	return
}
