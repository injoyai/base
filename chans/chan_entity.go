package chans

import (
	"context"
	"errors"
	"github.com/injoyai/conv"
	"sync"
	"time"
)

// Entity 实例
type Entity[T any] struct {
	key     string                                           //名称
	c       chan T                                           //通道
	handler func(ctx context.Context, no, count int, data T) //数据处理
	err     error                                            //错误信息
	num     int                                              //释放携程数量,默认1
	mu      sync.Mutex
	ctx     context.Context
	cancel  context.CancelFunc
	mCancel map[int]context.CancelFunc
}

// NewEntity
// @cap,通道大小
// @num,释放携程数量,默认1
func NewEntity[T any](num int, cap ...int) *Entity[T] {
	return NewEntityWithContext[T](context.Background(), num, cap...)
}

// NewEntityWithContext
// @ctx,上下文
// @cap,通道大小
// @num,释放携程数量,默认1
func NewEntityWithContext[T any](ctx context.Context, num int, cap ...int) *Entity[T] {
	ctxParent, parentCancel := context.WithCancel(ctx)
	data := &Entity[T]{
		key:     time.Now().Format("20060102150405"),
		c:       make(chan T, conv.Default[int](1, cap...)),
		ctx:     ctxParent,
		cancel:  parentCancel,
		num:     conv.Select[int](num < 1, 1, num),
		mCancel: make(map[int]context.CancelFunc),
	}
	for i := 0; i < data.num; i++ {
		childCtx, childCancel := context.WithCancel(ctxParent)
		data.mCancel[i] = childCancel
		go data.run(childCtx, i)
	}
	return data
}

// SetNum 重置协程数量
func (this *Entity[T]) SetNum(num int) {
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

func (this *Entity[T]) Cap() int {
	return cap(this.c)
}

func (this *Entity[T]) Len() int {
	return len(this.c)
}

// Close 关闭通道
func (this *Entity[T]) Close() error {
	this.CloseWithErr(errors.New("主动关闭"))
	return nil
}

func (this *Entity[T]) CloseWithErr(err error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if err != nil {
		this.err = err
		this.cancel()
	}
}

// Err 错误信息
func (this *Entity[T]) Err() error {
	return this.err
}

// Try 尝试加入队列(如果满了则忽略)
func (this *Entity[T]) Try(data ...T) (succ bool, err error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	for _, v := range data {
		select {
		case <-this.ctx.Done():
			return false, this.Err()
		case this.c <- v:
			succ = true
		default:
			//尝试加入队列失败
		}
	}
	return
}

func (this *Entity[T]) Do(data ...T) error {
	return this.Must(data...)
}

// Must 添加数据,通道关闭,返回错误信息
// @data,数据任意类型
func (this *Entity[T]) Must(data ...T) error {
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

// Timeout 尝试加入,在超时后返回错误
func (this *Entity[T]) Timeout(timeout time.Duration, data ...T) error {
	this.mu.Lock()
	defer this.mu.Unlock()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for _, v := range data {
		timer.Reset(timeout)
		select {
		case <-this.ctx.Done():
			return this.Err()
		case this.c <- v:
		case <-timer.C:
			return errors.New("超时")
		}
		return nil
	}
	return nil
}

// GetKey 获取标识
func (this *Entity[T]) GetKey() string {
	return this.key
}

// SetKey 设置标识,用于区别
// @key,标识
func (this *Entity[T]) SetKey(key string) *Entity[T] {
	this.key = key
	return this
}

// SetHandler 设置数据处理方法
// @ctx,上下文
// @no,协程序号
// @num 执行次数
// @data,数据
func (this *Entity[T]) SetHandler(fun func(ctx context.Context, no, count int, data T)) *Entity[T] {
	this.handler = fun
	return this
}

// run 多个select等待数据,按携程释放顺序,出来数据
// 例如携程顺序为,携程2,携程3,携程1,select的顺序就是2,3,1
func (this *Entity[T]) run(ctx context.Context, n int) {
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			return
		case v := <-this.c:
			if this.handler != nil {
				this.handler(ctx, n, i, v)
			}
		}
	}
}
