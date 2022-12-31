package chans

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"
)

func UID() string {
	c := make(chan string, 1)
	defer close(c)
	GetQueueFunc("uid").Do(func(no, num int) {
		x := md5.New()
		x.Write([]byte(fmt.Sprintf("%d#%d#%d", time.Now().UnixNano(), no, num)))
		bs := x.Sum(nil)
		c <- strings.ToUpper(string(bs))
	})
	return <-c
}
