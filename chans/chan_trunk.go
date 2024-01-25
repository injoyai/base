package chans

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// NewTrunk 消息总线,发布和订阅
func NewTrunk(num int, cap ...int) *Trunk {
	t := &Trunk{
		Entity:    NewEntity(num, cap...),
		subscribe: nil,
	}
	t.SetHandler(func(ctx context.Context, no, num int, data interface{}) {
		for _, sub := range t.subscribe {
			if sub != nil && sub.fn != nil {
				sub.fn(ctx, data)
			}
		}
	})
	return t
}

// Trunk 消息总线,发布和订阅
type Trunk struct {
	*Entity
	subscribe []*trunkSubscribe
	sync.Mutex
}

// Publish 发布接口输入
func (this *Trunk) Publish(data ...interface{}) error {
	return this.Entity.Do(data...)
}

// Subscribe 订阅消息总线
func (this *Trunk) Subscribe(handler func(ctx context.Context, data interface{})) string {
	key := fmt.Sprintf("%p%d", handler, time.Now().UnixNano())
	this.subscribe = append(this.subscribe, &trunkSubscribe{
		key: key,
		fn:  handler,
	})
	return key
}

// Unsubscribe 取消订阅
func (this *Trunk) Unsubscribe(key string) bool {
	if len(key) == 0 {
		return false
	}
	this.Lock()
	defer this.Unlock()
	for i, v := range this.subscribe {
		if v.key == key {
			this.subscribe = append(this.subscribe[:i], this.subscribe[i+1:]...)
			return true
		}
	}
	return false
}

type trunkSubscribe struct {
	key string
	fn  func(ctx context.Context, data interface{})
}
