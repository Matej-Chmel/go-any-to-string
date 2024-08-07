package internal

import (
	r "reflect"
	"strings"

	gs "github.com/Matej-Chmel/go-generic-stack"
)

// Counts the number of dimensions of an array or a slice
func countDimensions(val *r.Value) (d uint32) {
	t := val.Type()

	for {
		kind := t.Kind()

		if kind == r.Array || kind == r.Slice {
			d++
			t = t.Elem()
		} else {
			break
		}
	}

	return
}

// Returns the name of the type of the Value that represents
// a variable of a basic type
func FormatBasicType(val *r.Value) string {
	return val.Type().String()
}

// Returns the name of the type of the Value that represents
// a variable of a composite type
func formatCompositeType(val *r.Value) string {
	var builder strings.Builder
	stack := gs.Stack[typeInfo]{}
	stack.Push(newTypeInfo(val.Type(), false))

	for stack.HasItems() {
		top := stack.TopPointer()
		stack.Pop()
		aType := top.aType

		if top.endMap {
			builder.WriteRune(']')
		}

		switch aType.Kind() {
		case r.Array, r.Slice:
			builder.WriteString("[]")
			stack.Push(newTypeInfo(aType.Elem(), false))
		case r.Map:
			builder.WriteString("map[")
			stack.Push(newTypeInfo(aType.Elem(), true))
			stack.Push(newTypeInfo(aType.Key(), false))
		case r.Struct:
			builder.WriteString(aType.Name())
		default:
			builder.WriteString(aType.String())
		}
	}

	return builder.String()
}

// Returns the name of the type of the given Value
func FormatType(val *r.Value) string {
	if IsCompositeType(val) {
		return formatCompositeType(val)
	}

	return FormatBasicType(val)
}

// Returns true if given Value represents a composite type
func IsCompositeType(val *r.Value) bool {
	switch val.Kind() {
	case r.Array, r.Map, r.Pointer, r.Slice, r.Struct:
		return true
	}

	return false
}

// Returns true if given Value represents a nil pointer
func IsNil(val *r.Value) bool {
	switch val.Kind() {
	case r.Chan, r.Func, r.Map, r.Pointer, r.Slice:
		return val.IsNil()
	case r.Interface, r.Invalid:
		return !val.IsValid() || val.IsZero()
	}

	return false
}

// Internal struct for a Type in a stack
type typeInfo struct {
	aType  r.Type
	endMap bool
}

// Constructs new typeInfo
func newTypeInfo(aType r.Type, endMap bool) typeInfo {
	return typeInfo{aType: aType, endMap: endMap}
}
