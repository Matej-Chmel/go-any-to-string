package internal

import "reflect"

type Item struct {
	flag int
	ix   int
	keys []reflect.Value
	val  *reflect.Value
}

const (
	None = iota
	Bytes
	Dim2
	InnerDim
	KeyNext
	OtherDim
	Runes
	StructData
	Top
	ValueNext
)

func NewItem(f int, i int, v *reflect.Value) *Item {
	return &Item{flag: f, ix: i, keys: nil, val: v}
}
