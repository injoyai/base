package chans

import "C"
import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func NewDistribute() *Distribute {
	return &Distribute{}
}

/*
Distribute
分发
可选监听,均衡
*/
type Distribute struct {
	all []*SubsetChan
	one []*SubsetChan
	mu  sync.RWMutex
}

func (this *Distribute) Add(i interface{}, timeout ...time.Duration) {

	this.mu.RLock()
	if len(this.one) > 0 {
		idx := rand.Intn(len(this.one))
		c := this.one[idx].C
		defer c.Add(i, timeout...)
	}
	this.mu.RUnlock()

	for _, v := range this.all {
		v.C.Add(i, timeout...)
	}
}

func (this *Distribute) ChanOne(cap int) *SubsetChan {
	s := &SubsetChan{C: make(chan interface{}, cap)}
	s.key = fmt.Sprintf("%p", s)
	s.close = func() {
		for i, v := range this.one {
			if v.key == s.key {
				this.mu.Lock()
				this.one = append(this.one[:i], this.one[i+1:]...)
				this.mu.Unlock()
				break
			}
		}
	}
	this.mu.Lock()
	this.one = append(this.one, s)
	this.mu.Unlock()
	return s
}

func (this *Distribute) ChanAll(cap int) *SubsetChan {
	s := &SubsetChan{C: make(chan interface{}, cap)}
	s.key = fmt.Sprintf("%p", s)
	s.close = func() {
		for i, v := range this.all {
			if v.key == s.key {
				this.mu.Lock()
				this.all = append(this.all[:i], this.all[i+1:]...)
				this.mu.Unlock()
				break
			}
		}
	}
	this.mu.Lock()
	this.all = append(this.all, s)
	this.mu.Unlock()
	return s
}

type SubsetChan struct {
	C     Chan
	key   string
	close func()
}

func (this *SubsetChan) Close() {
	if this.close != nil {
		this.close()
	}
}
