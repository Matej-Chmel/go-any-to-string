package internal

import "reflect"

type Item struct {
	dim  uint32
	flag int
	ix   int
	keys []reflect.Value
	val  *reflect.Value
}

const (
	None = iota
	Bytes
	InnerDim
	KeyNext
	Runes
	StructData
	Top
	ValueNext
)

func NewItem(flag int, index int, val *reflect.Value) *Item {
	return &Item{
		dim:  0,
		flag: flag,
		ix:   index,
		keys: nil,
		val:  val,
	}
}

func (i *Item) GetCurrentDim() uint32 {
	return i.dim & 0xFFFF
}

func (i *Item) GetOriginalDim() uint32 {
	return (i.dim >> 16) & 0xFFFF
}

func (i *Item) SetCurrentDim(currentDim uint32) {
	i.dim = (i.dim & uint32(0xFFFF0000)) | (currentDim & 0xFFFF)
}

func (i *Item) SetOriginalDim(originalDim uint32) {
	i.dim = (i.dim & 0x0000FFFF) | ((originalDim & 0xFFFF) << 16)
}
