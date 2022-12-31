package chans

import (
	"context"
	"errors"
	"github.com/injoyai/conv"
	"sync"
	"time"
)

// Entity 性能,10万数据耗时35秒
type Entity struct {
	name         string                      //名称
	c            chan interface{}            //通道
	handler      func(int, int, interface{}) //数据处理
	err          error                       //错误信息
	num          int                         //释放携程数量,默认1
	mu           sync.Mutex
	parenCtx     context.Context
	parentCancel context.CancelFunc
	mCancel      map[int]context.CancelFunc
}

// NewEntity
// @cap,通道大小
// @num,释放携程数量,默认1
func NewEntity(num int, cap ...uint) *Entity {
	ctxParent, parentCancel := context.WithCancel(context.Background())
	data := &Entity{
		name:         time.Now().Format("20060102150405"),
		c:            make(chan interface{}, conv.GetDefaultUint(1, cap...)),
		parenCtx:     ctxParent,
		parentCancel: parentCancel,
		num:          conv.SelectInt(num < 1, 1, num),
		mCancel:      make(map[int]context.CancelFunc),
	}
	for i := 0; i < data.num; i++ {
		ctx, cancel := context.WithCancel(ctxParent)
		data.mCancel[i] = cancel
		go data.run(ctx, i)
	}
	return data
}

func (this *Entity) SetNum(num int) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if num > this.num {
		for i := this.num; i < num; i++ {
			ctx, cancel := context.WithCancel(this.parenCtx)
			this.mCancel[i] = cancel
			go this.run(ctx, i)
		}
	} else {
		for i := num; i < this.num; i++ {
			if cancel := this.mCancel[i]; cancel != nil {
				cancel()
			}
			delete(this.mCancel, i)
		}
	}
	this.num = num
}

func (this *QueueFunc) Cap() int {
	return cap(this.c)
}

func (this *QueueFunc) Len() int {
	return len(this.c)
}

//Close 关闭通道
func (this *Entity) Close(err ...error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.err == nil {
		this.err = conv.GetDefaultErr(errors.New("通道关闭"), err...)
		this.parentCancel()
		//close(this.c)
	}
}

// Err 错误信息
func (this *Entity) Err() error {
	return this.err
}

// Try 尝试加入队列(如果满了则忽略)
func (this *Entity) Try(data ...interface{}) error {
	this.mu.Lock()
	defer this.mu.Unlock()
	for _, v := range data {
		select {
		case <-this.parenCtx.Done():
			return this.err
		case this.c <- v:
		default: //尝试加入队列失败
		}
	}
	return nil
}

// Do 添加数据,通道关闭,返回错误信息
// @data,数据任意类型
func (this *Entity) Do(data ...interface{}) error {
	this.mu.Lock()
	defer this.mu.Unlock()
	for _, v := range data {
		select {
		case <-this.parenCtx.Done():
			return this.err
		case this.c <- v:
		}
	}
	return nil
}

// GetName 获取名字
func (this *Entity) GetName() string {
	return this.name
}

// SetName 设置名称,用于区别
// @name,实例名称
func (this *Entity) SetName(name string) *Entity {
	this.name = name
	return this
}

// SetHandler 设置数据处理方法
// @no,协程序号
// @num 执行次数
// @data,数据
func (this *Entity) SetHandler(fun func(no, num int, data interface{})) *Entity {
	this.handler = fun
	return this
}

// run 多个select等待数据,按携程释放顺序,出来数据
// 例如携程顺序为,携程2,携程3,携程1,select的顺序就是2,3,1
func (this *Entity) run(ctx context.Context, n int) {
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			return
		case v := <-this.c:
			if this.handler != nil {
				this.handler(n, i, v)
			}
		}
	}
}
