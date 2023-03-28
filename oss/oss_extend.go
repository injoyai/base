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
	OpenFile  = os.OpenFile
	Rename    = os.Rename
	Remove    = os.Remove
	RemoveAll = os.RemoveAll

	ReadDir = ioutil.ReadDir
)
