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
	runes
)

func newItem(f int, i int, v *reflect.Value) *item {
	return &item{flag: f, ix: i, val: v}
}
