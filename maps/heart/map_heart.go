package heart

import (
	"github.com/injoyai/conv"
	"sync"
	"time"
)

// Heart 心跳
type Heart struct {
	mOnline    map[string]*Info
	mOffline   map[string]int64
	mu         sync.RWMutex
	interval   time.Duration
	fnOnline   func([]*Info)
	fnOffline  func([]string)
	onlineList []*Info
}

type Info struct {
	No   string      //唯一标识符
	Date int64       //心跳时间
	Data interface{} //心跳内容
}

func NewHeart(interval time.Duration, fnOnline func([]*Info), fnOffline func([]string)) *Heart {
	data := &Heart{
		mOnline:   make(map[string]*Info),
		mOffline:  make(map[string]int64),
		interval:  interval,
		fnOnline:  fnOnline,
		fnOffline: fnOffline,
	}
	go data.run()
	return data
}

func (this *Heart) Keep(no string, v ...interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()
	data := &Info{
		No:   no,
		Date: time.Now().Unix(),
		Data: conv.GetDefault(v...).Interface(nil),
	}
	if _, ok := this.mOnline[no]; !ok {
		this.onlineList = append(this.onlineList, data)
	}
	this.mOnline[no] = data
	delete(this.mOffline, no)
}

func (this *Heart) run() {
	go func() {
		for {
			time.Sleep(this.interval / 2)
			if len(this.onlineList) > 0 {
				list := this.onlineList
				this.onlineList = []*Info{}
				this.fnOnline(list)
			}
		}
	}()
	for {
		time.Sleep(this.interval / 2)
		t := time.Now().Unix()
		list := []string{}
		this.mu.Lock()
		for i, v := range this.mOnline {
			if v.Date+int64(this.interval/time.Second) < t {
				this.mOffline[i] = t
				delete(this.mOnline, i)
				list = append(list, i)
			}
		}
		this.mu.Unlock()
		this.fnOffline(list)
	}
}
