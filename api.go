package goanytostring

import (
	"io"
	"reflect"

	ite "github.com/Matej-Chmel/go-any-to-string/internal"
)

// Convert any variable to a string
func AnyToString(a any) string {
	return AnyToStringCustom(a, NewOptions())
}

// Convert any variable to a string according to specified Options
func AnyToStringCustom(a any, o *Options) string {
	val := reflect.ValueOf(a)
	return ValueToStringCustom(&val, o)
}

// Write any variable to a Writer
func AnyToWriter(a any, w io.Writer) error {
	return AnyToWriterCustom(a, NewOptions(), w)
}

// Write any variable to a Writer according to specified Options
func AnyToWriterCustom(a any, o *Options, w io.Writer) error {
	val := reflect.ValueOf(a)
	return ValueToWriterCustom(&val, o, w)
}

// Convert Value to string
func ValueToString(val *reflect.Value) string {
	return ValueToStringCustom(val, NewOptions())
}

// Convert Value to string according to specified Options
func ValueToStringCustom(val *reflect.Value, o *Options) (res string) {
	if ite.IsNil(val) {
		res = "nil"

		if o.ShowType {
			res = res + " " + ite.FormatType(val)
		}
	} else if ite.IsCompositeType(val) {
		c := ite.NewCompositeConverter(o, val)
		return c.ConvertStackToString()
	} else {
		c := ite.NewLeafConverter(o)
		res = c.ConvertToString(val)

		if o.ShowType {
			res = res + " " + ite.FormatBasicType(val)
		}
	}

	return
}

// Write a Value to a Writer
func ValueToWriter(val *reflect.Value, w io.Writer) error {
	return ValueToWriterCustom(val, NewOptions(), w)
}

// Write a Value to a Writer according to specified Options
func ValueToWriterCustom(val *reflect.Value, o *Options, w io.Writer) error {
	return write(ValueToStringCustom(val, o), w)
}

// Internal function that writes a string to a Writer
func write(data string, w io.Writer) error {
	_, err := w.Write([]byte(data))
	return err
}
