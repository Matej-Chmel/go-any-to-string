package goanytostring_test

import (
	"runtime"
	"testing"

	ats "github.com/Matej-Chmel/go-any-to-string"
)

func check[T any](data T, expected string, t *testing.T, o ...ats.Options) {
	var actual string

	if len(o) > 0 {
		actual = ats.AnyToStringCustom(data, o[0])
	} else {
		actual = ats.AnyToString(data)
	}

	if actual == expected {
		return
	}

	_, _, line, ok := runtime.Caller(1)

	if !ok {
		t.Errorf("%s != %s", actual, expected)
	}

	t.Errorf("(line %d) %s != %s", line, actual, expected)
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
