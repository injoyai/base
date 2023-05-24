package oss

import (
	"io/ioutil"
	"os"
)

var (
	Exit      = os.Exit
	Chmod     = os.Chmod
	Create    = os.Create
	Open      = os.Open
	Rename    = os.Rename
	Remove    = os.Remove
	RemoveAll = os.RemoveAll

	ReadDir = ioutil.ReadDir
)
