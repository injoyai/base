package sort

import (
	"errors"
	"fmt"
	"github.com/injoyai/conv"
	json "github.com/json-iterator/go"
	"sort"
)

type Entity struct {
	list []interface{}
	fn   func(i, j interface{}) bool
	err  error
}

// Add 需要[]interface类型,或者任意类型的子元素
func (this *Entity) Add(list ...interface{}) *Entity {
	this.list = append(this.list, list...)
	return this
}

func (this *Entity) Set(list []interface{}) *Entity {
	this.list = list
	return this
}

func (this *Entity) Sort() ([]interface{}, error) {
	sort.Sort(this)
	return this.list, this.err
}

func (this *Entity) Bind(pointer interface{}) error {
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
func (this *Entity) Len() int {
	return len(this.list)
}

// Less 实现自带排序接口
func (this *Entity) Less(i, j int) bool {
	defer this.recover()
	return this.fn(this.list[i], this.list[j])
}

// Swap 实现自带排序接口
func (this *Entity) Swap(i, j int) {
	this.list[i], this.list[j] = this.list[j], this.list[i]
}

// recover 捕捉错误(主要类型强转错误)
func (this *Entity) recover() {
	if err := recover(); err != nil {
		this.err = errors.New(fmt.Sprintln(err))
	}
}
