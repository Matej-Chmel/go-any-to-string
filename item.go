package goanytostring

import "reflect"

type item struct {
	flag int
	ix   int
	keys []reflect.Value
	val  *reflect.Value
}

const (
	none = iota
	bytes
	keyNext
	runes
	structData
	top
	valueNext
)

func newItem(f int, i int, v *reflect.Value) *item {
	return &item{flag: f, ix: i, keys: nil, val: v}
}
