package shell

import (
	"github.com/injoyai/base/oss"
	"runtime"
)

func OpenBrowser(uri string) (err error) {
	switch runtime.GOOS {
	case oss.OS_windows:
		_, err = Exec("start", uri)
	case oss.OS_darwin:
		_, err = Exec("open", uri)
	case oss.OS_linux:
		_, err = Exec("xdg-open", uri)
	}
	return
}
