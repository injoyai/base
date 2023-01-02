package heart

import (
	"fmt"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"time"
)

const (
	BySet     = "bySet"
	ByTimeout = "byTimeout"
	ByHeart   = "byHeart"
)

// Heart 心跳
type Heart struct {
	mapOnline  *maps.Safe //在线缓存
	mapOffline *maps.Safe //离线缓存

	timeout     time.Duration //超时时间
	funcOnline  func([]*Info) //处理在线函数
	funcOffline func([]*Info) //处理离线函数
	lastOnline  []*Info       //最新在线数据
	lastOffline []*Info       //最新离线数据
}

type Info struct {
	Key    string      //唯一标识符
	Date   time.Time   //数据时间,无心跳则为0
	Data   interface{} //心跳内容
	Reason string      //在线离线原因
}

func (this *Info) String() string {
	return fmt.Sprintf("标识:%s 时间:%s", this.Key, this.Date.Format("2006-01-02 15:04:05"))
}

func (this *Info) IsTimeout(timeout time.Duration) bool {
	return this.Date.Unix() > 0 && time.Now().Sub(this.Date) > timeout
}

func NewInfoWithNoDate(key string, v ...interface{}) *Info {
	return &Info{
		Key:    key,
		Date:   time.Time{},
		Data:   conv.GetDefaultInterface(nil, v...),
		Reason: BySet,
	}
}

func NewInfoWithDate(key, reason string, v ...interface{}) *Info {
	return &Info{
		Key:    key,
		Date:   time.Now(),
		Data:   conv.GetDefaultInterface(nil, v...),
		Reason: reason,
	}
}

func New(timeout time.Duration, fnOnline func([]*Info), fnOffline func([]*Info)) *Heart {
	data := &Heart{
		mapOnline:   maps.NewSafe(),
		mapOffline:  maps.NewSafe(),
		timeout:     timeout,
		funcOnline:  fnOnline,
		funcOffline: fnOffline,
	}
	go data.run()
	return data
}

// SetOnline 设置在线
func (this *Heart) SetOnline(key string, info *Info) {
	if !this.mapOnline.Has(key) {
		this.lastOnline = append(this.lastOnline, info)
		this.mapOnline.Set(key, info)
	}
	this.mapOffline.Del(key)
}

// SetOffline 设置离线
func (this *Heart) SetOffline(key string) {
	info := NewInfoWithDate(key, BySet)
	if !this.mapOffline.Has(key) {
		this.lastOffline = append(this.lastOffline, info)
		this.mapOffline.Set(key, info)
	}
	this.mapOnline.Del(key)
}

// Keep 保持心跳
func (this *Heart) Keep(key string, v ...interface{}) {
	info := NewInfoWithDate(key, ByHeart, v...)
	this.SetOnline(key, info)
}

func (this *Heart) run() {
	for {
		time.Sleep(this.timeout / 2)

		lastOffline := this.lastOffline
		this.mapOnline.Range(func(key, value interface{}) bool {
			if value.(*Info).IsTimeout(this.timeout) {
				info := NewInfoWithDate(key.(string), ByTimeout)
				this.mapOffline.Set(key, info)
				this.mapOnline.Del(key)
				lastOffline = append(lastOffline, info)
			}
			return true
		})

		if len(lastOffline) > 0 {
			this.funcOffline(lastOffline)
		}

		if len(this.lastOnline) > 0 {
			this.funcOnline(this.lastOnline)
			this.lastOnline = []*Info{}
		}
	}
}
