package shell

import (
	"errors"
	"fmt"
	"github.com/injoyai/base/str"
	"os/exec"
	"runtime"
	"strings"
)

func Execf(format string, args ...interface{}) (string, error) {
	return Exec(fmt.Sprintf(format, args...))
}

func Exec(args ...string) (string, error) {
	list := append([]string{"/c"}, args...)
	switch runtime.GOOS {
	case "windows":
		result, err := exec.Command("cmd", list...).CombinedOutput()
		if err != nil {
			return "", err
		}
		result, err = str.GbkToUtf8(result)
		return string(result), err
	case "linux":
		list[0] = "-c"
		result, err := exec.Command("bash", list...).CombinedOutput()
		if err != nil {
			return "", err
		}
		return string(result), nil
	}
	return "", errors.New("未知操作系统:" + runtime.GOOS)
}

// Stop 结束程序 "taskkill.exe", "/f", "/im", "edge.exe"
func Stop(name string) error {
	switch runtime.GOOS {
	case "windows":
		result, err := Exec("taskkill.exe", "/f", "/im", name)
		if err != nil && !strings.Contains(err.Error(), "exit status") {
			return err
		} else if err == nil && !strings.Contains(result, "成功") {
			return errors.New(result)
		}
	}
	return nil
}

// Start 启动程序
// windows "cmd", "/c", "start ./xxx.exe"
func Start(path string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("cmd", "/c", "start "+path).Start()
	}
	return nil
}
