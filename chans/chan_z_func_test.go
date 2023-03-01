package chans

import (
	"testing"
)

func TestRange(t *testing.T) {
	{
		list := []int(nil)
		for v := range Range(5) {
			list = append(list, v)
		}
		if len(list) == 5 &&
			list[0] == 0 &&
			list[1] == 1 &&
			list[2] == 2 &&
			list[3] == 3 &&
			list[4] == 4 {
			t.Log("Range函数,1个参数,测试通过")
		} else {
			t.Error("Range函数,1个参数,测试失败")
		}
	}
	{
		list := []int(nil)
		for v := range Range(1, 5) {
			list = append(list, v)
		}
		if len(list) == 4 &&
			list[0] == 1 &&
			list[1] == 2 &&
			list[2] == 3 &&
			list[3] == 4 {
			t.Log("Range函数,2个参数,测试通过")
		} else {
			t.Error("Range函数,2个参数,测试失败")
		}
	}
	{
		list := []int(nil)
		for v := range Range(0, 5, 2) {
			list = append(list, v)
		}
		if len(list) == 3 &&
			list[0] == 0 &&
			list[1] == 2 &&
			list[2] == 4 {
			t.Log("Range函数,3个参数,测试通过")
		} else {
			t.Error("Range函数,3个参数,测试失败")
		}
	}
}
