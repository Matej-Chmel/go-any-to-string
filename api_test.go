package goanytostring_test

import (
	"runtime"
	"testing"

	ats "github.com/Matej-Chmel/go-any-to-string"
)

func check[T any](data T, expected string, t *testing.T, o ...ats.Options) {
	checkImpl(2, data, expected, t, o...)
}

func checkImpl[T any](skip int, data T, expected string, t *testing.T, o ...ats.Options) {
	var actual string

	if len(o) > 0 {
		actual = ats.AnyToStringCustom(data, o[0])
	} else {
		actual = ats.AnyToString(data)
	}

	if actual == expected {
		return
	}

	_, _, line, ok := runtime.Caller(skip)

	if !ok {
		t.Errorf("%s != %s", actual, expected)
	}

	t.Errorf("(line %d) %s != %s", line, actual, expected)
}

func checkPtr[T any](data T, expected string, t *testing.T, o ...ats.Options) {
	checkImpl(2, &data, expected, t, o...)
}

func TestArrays(t *testing.T) {
	check([...]bool{false, true}, "[false true]", t)
	check([...]byte{12, 34}, "[12 34]", t)
	check([...]int{1, 2, 3}, "[1 2 3]", t)
	check([]int{4, 5, 6}, "[4 5 6]", t)
	check([]rune{'A', 'B'}, "[65 66]", t)
	check([]string{"hello", "world"}, "[hello world]", t)

	o := ats.NewOptions()
	o.ByteAsString = true
	o.RuneAsString = true

	check([]byte{67, 68}, "CD", t, o)
	check([]rune{'A', 'B'}, "AB", t, o)
}

func TestArrayPointers(t *testing.T) {
	checkPtr([...]bool{false, true}, "&[false true]", t)
	checkPtr([...]byte{12, 34}, "&[12 34]", t)
	checkPtr([...]int{1, 2, 3}, "&[1 2 3]", t)
	checkPtr([]int{4, 5, 6}, "&[4 5 6]", t)
	checkPtr([]rune{'A', 'B'}, "&[65 66]", t)
	checkPtr([]string{"hello", "world"}, "&[hello world]", t)

	o := ats.NewOptions()
	o.ByteAsString = true
	o.RuneAsString = true

	checkPtr([]byte{67, 68}, "&CD", t, o)
	checkPtr([]rune{'A', 'B'}, "&AB", t, o)
}

func TestBasicTypes(t *testing.T) {
	check(false, "false", t)
	check(true, "true", t)
	check(make(chan int), "chan int", t)
	check(float32(12.34), "12.34", t)
	check(12.3456, "12.346", t)
	check(uint(12), "12", t)
	check(uint8(255), "255", t)
	check(uint16(65535), "65535", t)
	check(uint32(4294967295), "4294967295", t)
	check(uint64(18446744073709551615), "18446744073709551615", t)
	check(int(-12), "-12", t)
	check(int8(-128), "-128", t)
	check(int16(-32768), "-32768", t)
	check(int32(2147483647), "2147483647", t)
	check(int64(9223372036854775807), "9223372036854775807", t)
	check("hello world", "hello world", t)
	check(byte(65), "65", t)
	check('A', "65", t)

	o := ats.NewOptions()
	o.ByteAsString = true
	o.RuneAsString = true

	check(byte(65), "A", t, o)
	check('A', "A", t, o)
}

func TestComplex(t *testing.T) {
	check(1+1i, "1 + 1i", t)
	check(1.2+4.3i, "1.2 + 4.3i", t)
	check(1.2345+4.3456i, "1.234 + 4.346i", t)
	checkPtr(1+1i, "&1 + 1i", t)
	checkPtr(1.2+4.3i, "&1.2 + 4.3i", t)
	checkPtr(1.2345+4.3456i, "&1.234 + 4.346i", t)
}

func TestPointers(t *testing.T) {
	checkPtr(false, "&false", t)
	checkPtr(true, "&true", t)
	checkPtr(make(chan int), "&chan int", t)
	checkPtr(float32(12.34), "&12.34", t)
	checkPtr(12.3456, "&12.346", t)
	checkPtr(uint(12), "&12", t)
	checkPtr(uint8(255), "&255", t)
	checkPtr(uint16(65535), "&65535", t)
	checkPtr(uint32(4294967295), "&4294967295", t)
	checkPtr(uint64(18446744073709551615), "&18446744073709551615", t)
	checkPtr(int(-12), "&-12", t)
	checkPtr(int8(-128), "&-128", t)
	checkPtr(int16(-32768), "&-32768", t)
	checkPtr(int32(2147483647), "&2147483647", t)
	checkPtr(int64(9223372036854775807), "&9223372036854775807", t)
	checkPtr("hello world", "&hello world", t)
	checkPtr(byte(65), "&65", t)
	checkPtr('A', "&65", t)

	o := ats.NewOptions()
	o.ByteAsString = true
	o.RuneAsString = true

	checkPtr(byte(65), "&A", t, o)
	checkPtr('A', "&A", t, o)
}

type Example struct {
	a int
	B string
	c rune
}

type NestedExample struct {
	Example
	b string
	C rune
}

type SliceExample struct {
	bytes []byte
	ints  []int
}

func TestStruct(t *testing.T) {
	a := Example{12, "hello", '*'}
	b := NestedExample{Example{34, "world", '%'}, "super", 'X'}
	c := SliceExample{[]byte("hello world"), []int{1, 2, 3}}

	o := ats.NewOptions()
	o.ByteAsString = true
	o.RuneAsString = true

	check(a, "{12 hello *}", t, o)
	check(b, "{{34 world %} super X}", t, o)
	check(c, "{hello world [1 2 3]}", t, o)
}