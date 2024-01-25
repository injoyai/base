package chans

import "sync"

func NewOrder() *Order {
	return &Order{}
}

type Order struct {
	node []*OrderNode
	mu   sync.Mutex
}

func (this *Order) New() *OrderNode {
	this.mu.Lock()
	defer this.mu.Unlock()
	node := &OrderNode{
		sign: make(chan struct{}, 1),
	}
	if len(this.node) > 0 {
		//把左后一个节点指向新节点
		this.node[len(this.node)-1].next = node
	} else {
		//增加初始信号
		node.sign <- struct{}{}
	}
	//新节点添加到列表
	this.node = append(this.node, node)
	//新节点指向第一个节点,形成循环
	node.next = this.node[0]
	return node
}

type OrderNode struct {
	sign chan struct{}
	next *OrderNode
}

func (this *OrderNode) Do(fn func()) {
	<-this.sign
	fn()
	this.next.sign <- struct{}{}
}
