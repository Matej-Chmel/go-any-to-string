package goanytostring

import (
	r "reflect"

	ite "github.com/Matej-Chmel/go-any-to-string/internal"
)

// Type of a getter function of a sorting order
type GetLessType = ite.GetLessType

// Decorator for a less function and two values of a pointer type
func GetLessPointer(less KeyLessType) KeyLessType {
	return ite.GetLessPointer(less)
}

// Function type that returns true if first value should be sorted
// before the second one
type KeyLessType = ite.KeyLessType

// KeyLessType for bool values
func LessBool(a, b *r.Value) bool {
	return ite.LessBool(a, b)
}

// KeyLessType for complex numbers
func LessComplex(a, b *r.Value) bool {
	return ite.LessComplex(a, b)
}

// KeyLessType for floating-point numbers
func LessFloat(a, b *r.Value) bool {
	return ite.LessFloat(a, b)
}

// KeyLessType for signed integers
func LessInt(a, b *r.Value) bool {
	return ite.LessInt(a, b)
}

// KeyLessType for strings
func LessString(a, b *r.Value) bool {
	return ite.LessString(a, b)
}

// KeyLessType for unsigned integers
func LessUint(a, b *r.Value) bool {
	return ite.LessUint(a, b)
}
