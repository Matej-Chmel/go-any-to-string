package main

import (
	"fmt"

	at "github.com/Matej-Chmel/go-any-to-string"
)

type Example struct {
	a int
	B int
}

func main() {
	example := Example{a: 1, B: 2}
	fmt.Println(at.AnyToString(example))

	floatVal := 4.5678
	options := at.NewOptions()
	options.FloatDecimalPlaces = 2
	fmt.Println(at.AnyToStringCustom(floatVal, options))

	arr := [...]bool{false, true, false}
	options.ArraySep = " - "
	fmt.Println(at.AnyToStringCustom(arr, options))

	bytes := []byte{'H', 'e', 'l', 'l', 'o'}
	options.ByteAsString = true
	fmt.Println(at.AnyToStringCustom(bytes, options))
}
