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
	key     string                      //名称
	c       chan interface{}            //通道
	handler func(int, int, interface{}) //数据处理
	err     error                       //错误信息
	num     int                         //释放携程数量,默认1
	mu      sync.Mutex
	ctx     context.Context
	cancel  context.CancelFunc
	mCancel map[int]context.CancelFunc
}

// NewEntity
// @cap,通道大小
// @num,释放携程数量,默认1
func NewEntity(num int, cap ...int) *Entity {
	return NewEntityWithContext(context.Background(), num, cap...)
}

// NewEntityWithContext
// @ctx,上下文
// @cap,通道大小
// @num,释放携程数量,默认1
func NewEntityWithContext(ctx context.Context, num int, cap ...int) *Entity {
	ctxParent, parentCancel := context.WithCancel(ctx)
	data := &Entity{
		key:     time.Now().Format("20060102150405"),
		c:       make(chan interface{}, conv.GetDefaultInt(1, cap...)),
		ctx:     ctxParent,
		cancel:  parentCancel,
		num:     conv.SelectInt(num < 1, 1, num),
		mCancel: make(map[int]context.CancelFunc),
	}
	for i := 0; i < data.num; i++ {
		childCtx, childCancel := context.WithCancel(ctxParent)
		data.mCancel[i] = childCancel
		go data.run(childCtx, i)
	}
	return data
}

func (this *Entity) SetNum(num int) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if num > this.num {
		for i := this.num; i < num; i++ {
			ctx, cancel := context.WithCancel(this.ctx)
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
func (this *Entity) Close() error {
	this.CloseWithErr(errors.New("主动关闭"))
	return nil
}

func (this *Entity) CloseWithErr(err error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if err != nil {
		this.err = err
		this.cancel()
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
		case <-this.ctx.Done():
			return this.Err()
		case this.c <- v:
		default:
			//尝试加入队列失败
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
		case <-this.ctx.Done():
			return this.Err()
		case this.c <- v:
		}
	}
	return nil
}

// GetKey 获取标识
func (this *Entity) GetKey() string {
	return this.key
}

// SetKey 设置标识,用于区别
// @key,标识
func (this *Entity) SetKey(key string) *Entity {
	this.key = key
	return this
}

// SetHandler 设置数据处理方法
// @no,协程序号
// @num 执行次数
// @data,数据
func (this *Entity) SetHandler(fun func(no, count int, data interface{})) *Entity {
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
