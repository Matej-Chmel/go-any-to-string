package goanytostring

import "reflect"

type item struct {
	flag int
	ix   int
	val  *reflect.Value
}

const (
	none = iota
	bytes
	keyNext
	runes
	structData
	valueNext
)

func newItem(f int, i int, v *reflect.Value) *item {
	return &item{flag: f, ix: i, val: v}
}
