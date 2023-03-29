package chans

import "context"

// NewTrunk 消息总线,发布和订阅
func NewTrunk(num int, cap ...int) *Trunk {
	t := &Trunk{
		Entity:    NewEntity(num, cap...),
		subscribe: nil,
	}
	t.SetHandler(func(ctx context.Context, no, num int, data interface{}) {
		for _, fn := range t.subscribe {
			if fn != nil {
				fn(ctx, data)
			}
		}
	})
	return t
}

// Trunk 消息总线,发布和订阅
type Trunk struct {
	*Entity
	subscribe []func(ctx context.Context, data interface{})
}

// Publish 发布接口输入
func (this *Trunk) Publish(data ...interface{}) error {
	return this.Entity.Do(data...)
}

// Subscribe 订阅消息总线
func (this *Trunk) Subscribe(handler ...func(ctx context.Context, data interface{})) *Trunk {
	this.subscribe = append(this.subscribe, handler...)
	return this
}
