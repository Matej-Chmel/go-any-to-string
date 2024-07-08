package main

import (
	"fmt"
	"reflect"

	at "github.com/Matej-Chmel/go-any-to-string"
)

type Example struct {
	a int
	B int
}

func (e Example) String() string {
	return fmt.Sprintf("(%d, %d)", e.a, e.B)
}

func main() {
	example := Example{a: 1, B: 2}
	fmt.Println(at.AnyToString(example))

	// Use Options to ignore the default String() method
	opt := at.NewOptions()
	opt.IgnoreCustomMethod = true
	fmt.Println(at.AnyToStringCustom(example, opt))

	// Customize the output
	opt.ShowFieldNames = true
	opt.ShowType = true
	fmt.Println(at.AnyToStringCustom(example, opt))

	// Truncate floating-point values
	floatVal := 4.5678
	opt.FloatDecimalPlaces = 2
	opt.ShowType = false
	fmt.Println(at.AnyToStringCustom(floatVal, opt))

	// Change the start, seperator and end symbols of an array
	arr := [...]bool{false, true, false}
	opt.ArrayStart = "(("
	opt.ArraySep = " - "
	opt.ArrayEnd = "))"
	fmt.Println(at.AnyToStringCustom(arr, opt))

	// Convert a byte slice to a string
	bytes := []byte{'H', 'e', 'l', 'l', 'o'}
	opt.ByteAsString = true
	fmt.Println(at.AnyToStringCustom(bytes, opt))

	// Change the order of the keys in a map
	data := map[int]string{1: "hello", 2: "world"}
	opt.GetLessFunc = func(_ *reflect.Value) at.KeyLessType {
		return func(a, b *reflect.Value) bool {
			return a.Int() > b.Int()
		}
	}
	fmt.Println(at.AnyToStringCustom(data, opt))
}
