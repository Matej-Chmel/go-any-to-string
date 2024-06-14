package goanytostring

import (
	"io"
	"reflect"
	"strings"

	ite "github.com/Matej-Chmel/go-any-to-string/internal"
)

func AnyToString(a any) string {
	return AnyToStringCustom(a, NewOptions())
}

func AnyToStringCustom(a any, o Options) string {
	val := reflect.ValueOf(a)
	return ValueToStringCustom(&val, o)
}

func AnyToWriter(a any, w io.Writer) error {
	return AnyToWriterCustom(a, NewOptions(), w)
}

func AnyToWriterCustom(a any, o Options, w io.Writer) error {
	val := reflect.ValueOf(a)
	return ValueToWriterCustom(&val, o, w)
}

func ValueToString(val *reflect.Value) string {
	return ValueToStringCustom(val, NewOptions())
}

func ValueToStringCustom(val *reflect.Value, o Options) string {
	var builder strings.Builder
	c := ite.NewConverter(o, val, &builder)
	err := c.Run()

	if err != nil {
		return err.Error()
	}

	return builder.String()
}

func ValueToWriter(val *reflect.Value, w io.Writer) error {
	return ValueToWriterCustom(val, NewOptions(), w)
}

func ValueToWriterCustom(val *reflect.Value, o Options, w io.Writer) error {
	c := ite.NewConverter(o, val, w)
	return c.Run()
}
