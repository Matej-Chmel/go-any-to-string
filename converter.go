package goanytostring

import (
	"fmt"
	"io"
	r "reflect"
	"runtime"
	"strconv"
	"strings"
	"unsafe"

	gs "github.com/Matej-Chmel/go-generic-stack"
)

type converter struct {
	builder strings.Builder
	options Options
	stack   gs.Stack[*item]
	writer  io.Writer
}

func newConverter(o Options, val *r.Value, writer io.Writer) converter {
	c := converter{options: o, stack: gs.Stack[*item]{}, writer: writer}
	c.push(none, 0, val)
	return c
}

func (c *converter) run() error {
	for c.stack.HasItems() {
		top, _ := c.stack.Top()
		err := c.processItem(top)

		if err != nil {
			return err
		}
	}

	_, err := c.writer.Write([]byte(c.builder.String()))
	return err
}

func (c *converter) processItem(it *item) error {
	if it.flag == structData {
		return c.processStruct(it)
	}

	kind := it.val.Kind()

	switch kind {
	case r.Array, r.Slice:
		if it.flag > none {
			return c.processBytes(it)
		}

		if it.ix == 0 {
			elemKind := it.val.Type().Elem().Kind()

			if c.options.ByteAsString && elemKind == r.Uint8 {
				it.flag = bytes
				return c.processBytes(it)
			}

			if c.options.RuneAsString && elemKind == r.Int32 {
				it.flag = runes
				return c.processBytes(it)
			}
		}

		return c.processArray(it)
	case r.Map:
		return c.processMap(it)
	case r.Pointer:
		return c.processPointer(it)
	case r.Struct:
		return c.processStruct(it)
	}

	c.stack.Pop()

	if it.flag > 0 {
		if it.flag == bytes {
			return c.processByte(it.val)
		} else if it.flag == runes {
			return c.processRune(it.val)
		}
	}

	switch kind {

	case r.Bool:
		return c.processBool(it.val)
	case r.Chan:
		return c.processChan(it.val)
	case r.Complex64, r.Complex128:
		return c.processComplex(it.val)
	case r.Float32:
		return c.processFloat(it.val, 32)
	case r.Float64:
		return c.processFloat(it.val, 64)
	case r.Func:
		return c.processFunc(it.val)

	case r.Int32:
		if c.options.RuneAsString {
			return c.processRune(it.val)
		}

		return c.processInt(it.val)

	case r.Int, r.Int8, r.Int16, r.Int64:
		return c.processInt(it.val)

	case r.Interface, r.Invalid:
		return c.processInterface(it.val)

	case r.Uint, r.Uint16, r.Uint32, r.Uint64:
		return c.processUint(it.val)

	case r.Uint8:
		if c.options.ByteAsString {
			return c.processByte(it.val)
		}

		return c.processUint(it.val)

	case r.String:
		return c.processString(it.val)

	case r.Uintptr:
		return c.processUintptr(it.val)
	case r.UnsafePointer:
		return c.processUnsafe(it.val)

	default:
		return c.write(kind.String())
	}
}

func (c *converter) push(f int, i int, v *r.Value) {
	c.stack.Push(newItem(f, i, v))
}

func (c *converter) processArray(it *item) error {
	l := it.val.Len()

	if it.ix == 0 {
		err := c.writeRune('[')

		if err != nil {
			return err
		}
	} else if it.ix < l {
		err := c.writeRune(' ')

		if err != nil {
			return err
		}
	} else if it.ix == l {
		err := c.writeRune(']')

		if err != nil {
			return err
		}

		c.stack.Pop()
		return nil
	}

	elem := it.val.Index(it.ix)
	c.push(none, 0, &elem)

	it.ix++
	return nil
}

func (c *converter) processBytes(it *item) error {
	l := it.val.Len()

	if it.ix == l {
		c.stack.Pop()
		return nil
	}

	elem := it.val.Index(it.ix)
	c.push(it.flag, 0, &elem)

	it.ix++
	return nil
}

func (c *converter) processMap(it *item) error {
	if it.flag == none && it.ix == 0 {
		it.keys = it.val.MapKeys()

		err := c.writeRune('{')

		if err != nil {
			return err
		}

		it.flag = keyNext
	} else if it.flag == keyNext && it.ix < it.val.Len() {
		err := c.writeRune(' ')

		if err != nil {
			return err
		}
	} else if it.ix == it.val.Len() {
		err := c.writeRune('}')

		if err != nil {
			return err
		}

		c.stack.Pop()
		return nil
	}

	key := it.keys[it.ix]

	if it.flag == keyNext {
		c.push(none, 0, &key)
		it.flag = valueNext
	} else if it.flag == valueNext {
		c.writeRune(':')
		val := it.val.MapIndex(key)
		c.push(none, 0, &val)
		it.flag = keyNext
		it.ix++
	}

	return nil
}

func (c *converter) processPointer(it *item) error {
	err := c.write("&")

	if err != nil {
		return err
	}

	if elem := it.val.Elem(); elem.Kind() == r.Struct {
		it.flag = structData
	} else {
		it.val = &elem
	}

	return nil
}

func (c *converter) processStruct(it *item) error {
	if it.flag != structData {
		tmp := r.New(it.val.Type())
		tmp.Elem().Set(*it.val)
		elem := tmp.Elem()
		it.val = &elem
	}

	if it.ix == 0 {
		err := c.writeRune('{')

		if err != nil {
			return err
		}
	} else if it.ix < it.val.NumField() {
		err := c.writeRune(' ')

		if err != nil {
			return err
		}
	} else if it.ix == it.val.NumField() {
		err := c.writeRune('}')

		if err != nil {
			return err
		}

		c.stack.Pop()
		return nil
	}

	field := it.val.Field(it.ix)

	if field.CanInterface() {
		c.push(none, 0, &field)
	} else {
		addr := unsafe.Pointer(field.UnsafeAddr())
		data := r.NewAt(field.Type(), addr).Elem()
		c.push(none, 0, &data)
	}

	it.ix++
	return nil
}

func (c *converter) processBool(val *r.Value) error {
	return c.write(strconv.FormatBool(val.Bool()))
}

func (c *converter) processByte(val *r.Value) error {
	return c.writeByte(byte(val.Uint()))
}

func (c *converter) processChan(val *r.Value) error {
	err := c.write("chan ")

	if err != nil {
		return err
	}

	return c.write(val.Type().Elem().String())
}

func floatToString(f float64, bitSize int) string {
	s := strconv.FormatFloat(f, 'f', 3, bitSize)
	s = strings.TrimRight(s, "0")
	return strings.TrimRight(s, ".")
}

func (c *converter) processComplex(val *r.Value) error {
	realPart := floatToString(real(val.Complex()), 64)
	imagPart := floatToString(imag(val.Complex()), 64)
	return c.write(fmt.Sprintf("%s + %si", realPart, imagPart))
}

func (c *converter) processFloat(val *r.Value, bitSize int) error {
	return c.write(floatToString(val.Float(), bitSize))
}

func (c *converter) processFunc(val *r.Value) error {
	typ := val.Type()
	in, out := typ.NumIn(), typ.NumOut()

	name := runtime.FuncForPC(val.Pointer()).Name()
	parts := strings.Split(name, ".")
	lastIx := len(parts) - 1
	// pkg := strings.Join(parts[:lastIx])
	funcName := parts[lastIx]

	err := c.write(funcName)
	if err != nil {
		return err
	}

	err = c.writeRune('(')
	if err != nil {
		return err
	}

	for i := 0; i < in; i++ {
		if i > 0 {
			err = c.write(", ")
			if err != nil {
				return err
			}
		}

		err = c.write(typ.In(i).String())
		if err != nil {
			return err
		}
	}

	err = c.write(") ")
	if err != nil {
		return err
	}

	if out > 1 {
		err = c.writeRune('(')
		if err != nil {
			return err
		}
	}

	for i := 0; i < out; i++ {
		if i > 0 {
			err = c.write(", ")
			if err != nil {
				return err
			}
		}

		err = c.write(typ.Out(i).String())
		if err != nil {
			return err
		}
	}

	if out > 1 {
		err = c.writeRune(')')
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *converter) processInt(val *r.Value) error {
	return c.write(strconv.FormatInt(val.Int(), 10))
}

func (c *converter) processInterface(_ *r.Value) error {
	return c.write("interface{}")
}

func (c *converter) processRune(val *r.Value) error {
	return c.writeRune(rune(val.Int()))
}

func (c *converter) processString(val *r.Value) error {
	return c.write(val.String())
}

func (c *converter) processUint(val *r.Value) error {
	return c.write(strconv.FormatUint(val.Uint(), 10))
}

func (c *converter) processUintptr(val *r.Value) error {
	return c.write(fmt.Sprintf("0x%X", val.Uint()))
}

func (c *converter) processUnsafe(val *r.Value) error {
	return c.write(fmt.Sprintf("Ux%X", val.Pointer()))
}

func (c *converter) write(s string) error {
	fmt.Println(s)
	_, err := c.builder.WriteString(s)
	return err
}

func (c *converter) writeByte(b byte) error {
	return c.builder.WriteByte(b)
}

func (c *converter) writeRune(r rune) error {
	_, err := c.builder.WriteRune(r)
	return err
}
