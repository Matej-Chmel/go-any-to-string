package goanytostring

import (
	"io"
	"reflect"
	"strings"
)

func AnyToString(a any) string {
	return AnyToStringCustom(a, NewOptions())
}

func AnyToStringCustom(a any, o Options) string {
	var builder strings.Builder
	val := reflect.ValueOf(a)
	c := newConverter(o, &val, &builder)
	err := c.run()

	if err != nil {
		return err.Error()
	}

	return builder.String()
}

func AnyToWriter(a any, w io.Writer) error {
	return AnyToWriterCustom(a, NewOptions(), w)
}

func AnyToWriterCustom(a any, o Options, w io.Writer) error {
	val := reflect.ValueOf(a)
	c := newConverter(o, &val, w)
	return c.run()
}
