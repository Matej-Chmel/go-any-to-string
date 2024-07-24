package goanytostring_test

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"
	"unsafe"

	ats "github.com/Matej-Chmel/go-any-to-string"
)

type Example struct {
	a int
	B string
	c rune
}

type ExampleCustom struct {
	a rune
	b rune
	c rune
}

func (e ExampleCustom) String() string {
	return fmt.Sprintf("%c -> %c -> %c", e.a, e.b, e.c)
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

// Wrapper around the test state
type tester struct {
	failed bool
	*testing.T
}

func newTester(t *testing.T) *tester {
	return &tester{failed: false, T: t}
}

func (t *tester) fail(skip int, format string, data ...any) {
	_, _, line, ok := runtime.Caller(skip)

	if ok {
		format = fmt.Sprintf("(line %d) %s", line, format)
	}

	t.Errorf(format, data...)
	t.failed = true
}

func check[T any](data T, expected string, t *tester, o ...*ats.Options) {
	checkImpl(2, data, expected, t, o...)
}

func checkImpl[T any](skip int, data T, expected string, t *tester, o ...*ats.Options) {
	if t.failed {
		return
	}

	var actual string

	if len(o) > 0 {
		actual = ats.AnyToStringCustom(data, o[0])
	} else {
		actual = ats.AnyToString(data)
	}

	if actual == expected {
		return
	}

	t.fail(skip, "\n\n%s\n\n!=\n\n%s", actual, expected)
}

func checkPtr[T any](data T, expected string, t *tester, o ...*ats.Options) {
	checkImpl(2, &data, expected, t, o...)
}

func hello(a int) string {
	return fmt.Sprintf("%d", a)
}

func readFile(path string) string {
	file, err := os.Open(path)

	if err != nil {
		return ""
	}

	data, err := io.ReadAll(file)

	if err != nil {
		return ""
	}

	res := strings.TrimSpace(string(data))
	return strings.ReplaceAll(res, "\r", "")
}

func tuple(a, b int, c string) (int, int, string) {
	return a, b, c
}

func Test2D(ot *testing.T) {
	t := newTester(ot)
	data := [][]int32{
		{1, 2, 3},
		{4, -5, 6},
		{7, 8, -9, 0, 1},
	}
	exp := readFile("test_data/2D.txt")
	check(data, exp, t)
}

func Test3D(ot *testing.T) {
	t := newTester(ot)
	data := [][][]int32{
		{
			{-1, -2, -3},
			{4, 5, 6},
			{7, 8, 9},
		},
		{
			{0, 0, 6},
			{0, 0},
		},
		{
			{1, 1, -1, -1, 0},
			{0},
			{1, 1},
		},
	}
	exp := readFile("test_data/3D.txt")
	check(data, exp, t)
}

func Test4D(ot *testing.T) {
	t := newTester(ot)
	data := [][][][]int32{
		{
			{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
			{
				{-4, -5, -6},
				{-7, -8, -9},
				{-1, -2, -3},
			},
		},
		{
			{
				{0, 0},
				{0, 0},
				{0},
			},
			{
				{100, 10},
				{1000, 10},
				{100_000_000},
			},
		},
		{
			{
				{-1, 1},
				{1, -1},
			},
		},
	}
	exp := readFile("test_data/4D.txt")
	check(data, exp, t)
}

func TestArrays(ot *testing.T) {
	t := newTester(ot)
	check([...]bool{false, true}, "[false true]", t)
	check([...]byte{12, 34}, "[12 34]", t)
	check([...]int{1, 2, 3}, "[1 2 3]", t)
	check([]int{4, 5, 6}, "[4 5 6]", t)
	check([]rune{'A', 'B'}, "[65 66]", t)
	check([]string{"hello", "world"}, "[hello world]", t)

	checkPtr([...]bool{false, true}, "&[false true]", t)
	checkPtr([...]byte{12, 34}, "&[12 34]", t)
	checkPtr([...]int{1, 2, 3}, "&[1 2 3]", t)
	checkPtr([]int{4, 5, 6}, "&[4 5 6]", t)
	checkPtr([]rune{'A', 'B'}, "&[65 66]", t)
	checkPtr([]string{"hello", "world"}, "&[hello world]", t)

	o := ats.NewOptions()
	o.ByteAsString = true
	o.RuneAsString = true

	check([]byte{67, 68}, "CD", t, o)
	check([]rune{'A', 'B'}, "AB", t, o)

	checkPtr([]byte{67, 68}, "&CD", t, o)
	checkPtr([]rune{'A', 'B'}, "&AB", t, o)
}

func TestBasicTypes(ot *testing.T) {
	t := newTester(ot)
	check(true, "true", t)
	check(make(chan int), "chan int", t)
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

func TestComplex(ot *testing.T) {
	t := newTester(ot)
	check(1+1i, "(1+1i)", t)
	check(1.2+4.3i, "(1.2+4.3i)", t)
	check(1.2345+4.3456i, "(1.234+4.346i)", t)
}

func TestFloat(ot *testing.T) {
	t := newTester(ot)
	check(0.0, "0.0", t)
	check(1.0, "1.0", t)
	check(1.020, "1.02", t)
	check(1.0209, "1.021", t)
	check(1.0211, "1.021", t)
	check(127.1239, "127.124", t)
}

func TestFormat(ot *testing.T) {
	t := newTester(ot)
	a := [...]int{4, 5, 6}
	m := map[int]string{12: "hello", 34: "world"}

	o := ats.NewOptions()
	o.ArrayStart = "<< "
	o.ArrayEnd = " >>"
	o.ArraySep = ", "
	o.MapStart = ":: "
	o.MapEnd = "!"
	o.MapSepVal = " - "

	check(a, "<< 4, 5, 6 >>", t, o)
	check(m, ":: 12:hello - 34:world!", t, o)

	o.ShowType = true
	check(a, "[]int << 4, 5, 6 >>", t, o)
	check(m, "map[int]string :: 12:hello - 34:world!", t, o)
}

func TestFunc(ot *testing.T) {
	t := newTester(ot)
	check(hello, "hello(int) string", t)
	checkPtr(hello, "&hello(int) string", t)

	check(tuple, "tuple(int, int, string) (int, int, string)", t)
	checkPtr(tuple, "&tuple(int, int, string) (int, int, string)", t)

	actual := ats.AnyToString(func(i int) int {
		return i + 1
	})
	check(actual, "func1(int) int", t)
}

func TestInterface(ot *testing.T) {
	t := newTester(ot)
	var i interface{}
	check(i, "nil", t)
	checkPtr(i, "&nil", t)
}

func TestMap(ot *testing.T) {
	t := newTester(ot)
	fa, tr := false, true
	i := map[int]string{12: "hello", 34: "world"}
	s := map[string]*bool{"F": &fa, "T": &tr}

	actual := ats.AnyToString(i)
	check(actual, "{12:hello 34:world}", t)

	actual = ats.AnyToString(s)
	check(actual, "{F:&false T:&true}", t)
}

func TestMemory(ot *testing.T) {
	t := newTester(ot)
	check(uintptr(0x12345678), "0x12345678", t)
	checkPtr(uintptr(0x12345678), "&0x12345678", t)

	check(unsafe.Pointer(uintptr(0x34125678)), "Ux34125678", t)
	checkPtr(unsafe.Pointer(uintptr(0x34125678)), "&Ux34125678", t)
}

func TestPointers(ot *testing.T) {
	t := newTester(ot)
	checkPtr(false, "&false", t)
	checkPtr(true, "&true", t)
	checkPtr(make(chan int), "&chan int", t)
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

func TestStruct(ot *testing.T) {
	t := newTester(ot)
	a := Example{12, "hello", '*'}
	b := NestedExample{Example{34, "world", '%'}, "super", 'X'}
	c := SliceExample{[]byte("hello world"), []int{1, 2, 3}}
	d := ExampleCustom{'A', 'b', 'C'}

	o := ats.NewOptions()
	o.ByteAsString = true
	o.RuneAsString = true

	check(a, "{12 hello *}", t, o)
	check(b, "{{34 world %} super X}", t, o)
	check(c, "{hello world [1 2 3]}", t, o)
	check(d, "A -> b -> C", t, o)

	checkPtr(a, "&{12 hello *}", t, o)
	checkPtr(b, "&{{34 world %} super X}", t, o)
	checkPtr(c, "&{hello world [1 2 3]}", t, o)
	checkPtr(d, "A -> b -> C", t, o)

	o.IgnoreCustomMethod = true
	o.ShowFieldNames = true

	check(a, "{a:12 B:hello c:*}", t, o)
	check(b, "{Example:{a:34 B:world c:%} b:super C:X}", t, o)
	check(c, "{bytes:hello world ints:[1 2 3]}", t, o)
	check(d, "{a:A b:b c:C}", t, o)

	checkPtr(a, "&{a:12 B:hello c:*}", t, o)
	checkPtr(b, "&{Example:{a:34 B:world c:%} b:super C:X}", t, o)
	checkPtr(c, "&{bytes:hello world ints:[1 2 3]}", t, o)
	checkPtr(d, "&{a:A b:b c:C}", t, o)
}

func TestZero(ot *testing.T) {
	t := newTester(ot)
	check[interface{}](nil, "nil", t)
	check[*int](nil, "nil", t)
	check(false, "false", t)
	check(float32(0.0), "0.0", t)
	check(0.0, "0.0", t)
	check(uint(0), "0", t)
	check(uint8(0), "0", t)
	check(uint16(0), "0", t)
	check(uint32(0), "0", t)
	check(uint64(0), "0", t)
	check(int(0), "0", t)
	check(int8(0), "0", t)
	check(int16(0), "0", t)
	check(int32(0), "0", t)
	check(int64(0), "0", t)
	check("", "", t)
	check(byte(0), "0", t)
	check('\000', "0", t)
}
