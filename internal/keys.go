package internal

import (
	r "reflect"
	"sort"
)

// Default getter function of a sorting order
func DefaultGetLess(key *r.Value) KeyLessType {
	switch kind := key.Kind(); kind {
	case r.Bool:
		return LessBool
	case r.Complex64, r.Complex128:
		return LessComplex
	case r.Float32, r.Float64:
		return LessFloat
	case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
		return LessInt
	case r.Pointer:
		elem := key.Elem()
		return GetLessPointer(DefaultGetLess(&elem))
	case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64:
		return LessInt
	}

	return LessString
}

// Type of a getter function of a sorting order
type GetLessType = func(*r.Value) KeyLessType

// Decorator for a less function and two values of a pointer type
func GetLessPointer(less KeyLessType) KeyLessType {
	return func(a, b *r.Value) bool {
		aElem := a.Elem()
		bElem := b.Elem()
		return less(&aElem, &bElem)
	}
}

// Function type that returns true if first value should be sorted
// before the second one
type KeyLessType = func(*r.Value, *r.Value) bool

// Implementation of sort.Interface
type Keys struct {
	data []r.Value
	less KeyLessType
}

// Returns the number of keys
func (k *Keys) Len() int {
	return len(k.data)
}

// Returns true if key at index i should be sorted before key at index j
func (k *Keys) Less(i, j int) bool {
	return k.less(&k.data[i], &k.data[j])
}

// Swaps keys at indices i and j
func (k *Keys) Swap(i, j int) {
	k.data[i], k.data[j] = k.data[j], k.data[i]
}

// KeyLessType for bool values
func LessBool(a, b *r.Value) bool {
	aVal := a.Bool()
	return aVal == false || aVal == b.Bool()
}

// KeyLessType for complex numbers
func LessComplex(a, b *r.Value) bool {
	return magnitude(a.Complex()) <= magnitude(b.Complex())
}

// KeyLessType for floating-point numbers
func LessFloat(a, b *r.Value) bool {
	return a.Float() <= b.Float()
}

// KeyLessType for signed integers
func LessInt(a, b *r.Value) bool {
	return a.Int() <= b.Int()
}

// KeyLessType for strings
func LessString(a, b *r.Value) bool {
	return a.String() <= b.String()
}

// KeyLessType for unsigned integers
func LessUint(a, b *r.Value) bool {
	return a.Uint() <= b.Uint()
}

// Returns the magnitude of a complex number
func magnitude(val complex128) float64 {
	img := imag(val)
	rel := real(val)
	return img*img + rel*rel
}

// Sorts keys in the slice data according to sort order
// that is returned when the function getLess is invoked
// with the first key
func SortKeys(data []r.Value, getLess GetLessType) {
	if len(data) <= 0 {
		return
	}

	keys := &Keys{data: data, less: getLess(&data[0])}
	sort.Sort(keys)
}
