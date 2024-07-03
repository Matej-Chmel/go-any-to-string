package internal

import "reflect"

// Item in a stack
type Item struct {
	// Two-part bit array, higher 16 bits represent number of dimensions
	// of the original array or slice, lower 16 bits number of dimensions
	// of the current layer. It's 0, if this Item isn't an array or slice.
	dim uint32
	// Flag indicating the current stage of processing
	flag uint
	// If Item is an array or slice, ix is an index into that data
	ix int
	// If Item is a map, the order of keys is saved here
	keys []reflect.Value
	// The reflection data of this Item
	val *reflect.Value
}

const (
	// The item has no special flag
	None uint = iota
	// Item is a part of a byte array of slice
	Bytes
	// Item is a part of multidimensional array of slice
	InnerDim
	// Item is a map and a key should be processed in the next stage
	KeyNext
	// Item is a part of a rune array of slice
	Runes
	// Item is a struct and a pointer to this struct has been created
	StructData
	// Item is a map and a value should be processed in the next stage
	ValueNext
)

// Constructs a new Item
func NewItem(flag uint, index int, val *reflect.Value) *Item {
	return &Item{
		dim:  0,
		flag: flag,
		ix:   index,
		keys: nil,
		val:  val,
	}
}

// Returns the number of dimensions of the current layer
func (i *Item) GetCurrentDim() uint32 {
	return i.dim & 0xFFFF
}

// Returns the number of dimensions of the outer most layer
func (i *Item) GetOriginalDim() uint32 {
	return (i.dim >> 16) & 0xFFFF
}

// Sets the number of dimensions of the current layer
func (i *Item) SetCurrentDim(currentDim uint32) {
	i.dim = (i.dim & uint32(0xFFFF0000)) | (currentDim & 0xFFFF)
}

// Sets the number of dimensions of the outer most layer
func (i *Item) SetOriginalDim(originalDim uint32) {
	i.dim = (i.dim & 0x0000FFFF) | ((originalDim & 0xFFFF) << 16)
}
