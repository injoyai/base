package sort

import (
	"errors"
	"fmt"
	"github.com/injoyai/conv"
	json "github.com/json-iterator/go"
	"sort"
)

// List 排序任意列表
func List(fn func(i, j interface{}) bool) *list {
	return &list{fn: fn}
}

type list struct {
	list []interface{}
	fn   func(i, j interface{}) bool
}

// Add 需要[]interface类型,或者任意类型的子元素
func (this *list) Add(list ...interface{}) *list {
	this.list = append(this.list, list...)
	return this
}

func (this *list) Set(list []interface{}) *list {
	this.list = list
	return this
}

func (this *list) Sort() (_ []interface{}, err error) {
	if e := recover(); e != nil {
		err = errors.New(fmt.Sprintln(e))
	}
	sort.Sort(this)
	return this.list, nil
}

func (this *list) Bind(pointer interface{}) error {
	this.Add(conv.Interfaces(pointer)...)
	data, err := this.Sort()
	if err != nil {
		return err
	}
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, pointer)
}

//------------------------

// Len 实现自带排序接口
func (this *list) Len() int {
	return len(this.list)
}

// Less 实现自带排序接口
func (this *list) Less(i, j int) bool {
	return this.fn(this.list[i], this.list[j])
}

// Swap 实现自带排序接口
func (this *list) Swap(i, j int) {
	this.list[i], this.list[j] = this.list[j], this.list[i]
}
