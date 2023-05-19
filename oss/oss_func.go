package oss

import (
	"fmt"
	"github.com/injoyai/conv"
	"os"
	"os/signal"
	"os/user"
	"strings"
	"sync"
	"syscall"
)

var (
	exitFunc []func()
	exitOnce sync.Once
)

// ListenExit 监听退出信号
func ListenExit(handler ...func()) {
	exitOnce.Do(func() {
		exitChan := make(chan os.Signal)
		signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGTERM)
		go func() {
			<-exitChan
			for _, v := range exitFunc {
				v()
			}
			os.Exit(-127)
		}()
	})
	exitFunc = append(exitFunc, handler...)
}

// Wait 一直等待
func Wait() { select {} }

// UserHome 获取用户账号路径 例 C:/Users/injoy/
func UserHome() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(u.HomeDir, `\`, `/`) + "/", nil
}

// Input 监听用户输入
// reader := bufio.NewReader(os.Stdin)
// msg, _ := reader.ReadString('\n')
func Input(hint ...interface{}) (s string) {
	if len(hint) > 0 {
		fmt.Println(hint...)
	}
	fmt.Scanln(&s)
	return
}

// InputVar 监听用户输入
func InputVar(hint ...interface{}) *conv.Var { return conv.New(Input(hint...)) }
