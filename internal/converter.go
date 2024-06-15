package internal

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

type Converter struct {
	builder strings.Builder
	options Options
	stack   gs.Stack[*Item]
	writer  io.Writer
}

func NewConverter(o Options, val *r.Value, writer io.Writer) Converter {
	c := Converter{options: o, stack: gs.Stack[*Item]{}, writer: writer}
	c.push(Top, 0, val)
	return c
}

func (c *Converter) Run() error {
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

func countDimensions(val *r.Value) (d int) {
	t := val.Type()

	for {
		kind := t.Kind()

		if kind == r.Array || kind == r.Slice {
			d++
			t = t.Elem()
		} else {
			break
		}
	}

	return
}

func formatType(t r.Type, top bool) (s string) {
	switch t.Kind() {
	case r.Array, r.Slice:
		s = fmt.Sprintf("[]%s", formatType(t.Elem(), false))
	case r.Map:
		s = fmt.Sprintf("map[%s]%s", formatType(t.Key(), false), formatType(t.Elem(), false))
	case r.Struct:
		s = t.Name()
	default:
		s = t.String()
	}

	if top {
		s += " "
	}

	return
}

func (c *Converter) processItem(it *Item) error {
	if it.flag == Top {
		it.flag = None

		if c.options.ShowType {
			err := c.write(formatType(it.val.Type(), true))

			if err != nil {
				return err
			}
		}
	}

	if it.flag == StructData {
		return c.processStruct(it)
	}

	kind := it.val.Kind()

	switch kind {
	case r.Array, r.Slice:
		if it.flag == Bytes || it.flag == Runes {
			return c.processBytes(it)
		}

		if it.ix == 0 {
			elemKind := it.val.Type().Elem().Kind()

			if c.options.ByteAsString && elemKind == r.Uint8 {
				it.flag = Bytes
				return c.processBytes(it)
			}

			if c.options.RuneAsString && elemKind == r.Int32 {
				it.flag = Runes
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
		if it.flag == Bytes {
			return c.processByte(it.val)
		} else if it.flag == Runes {
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

func (c *Converter) push(f int, i int, v *r.Value) {
	c.stack.Push(NewItem(f, i, v))
}

func (c *Converter) processArray(it *Item) error {
	if it.flag == None {
		dim := countDimensions(it.val)

		if dim == 2 {
			it.flag = Dim2
		} else {
			it.flag = OtherDim
		}
	}

	if it.flag == Dim2 {
		return c.processArray2D(it)
	}

	if it.ix == 0 && it.flag != InnerDim {
		if err := c.write(c.options.ArrayStart); err != nil {
			return err
		}
	} else if l := it.val.Len(); it.ix > 0 && it.ix < l {
		if err := c.write(c.options.ArraySep); err != nil {
			return err
		}
	} else if it.ix == l {
		if it.flag != InnerDim {
			if err := c.write(c.options.ArrayEnd); err != nil {
				return err
			}
		}

		c.stack.Pop()
		return nil
	}

	elem := it.val.Index(it.ix)
	c.push(None, 0, &elem)

	it.ix++
	return nil
}

func (c *Converter) processArray2D(it *Item) error {
	if l := it.val.Len(); it.ix > 0 && it.ix < l {
		if err := c.write(c.options.ArraySep2D); err != nil {
			return err
		}
	} else if it.ix == l {
		c.stack.Pop()
		return nil
	}

	elem := it.val.Index(it.ix)
	c.push(InnerDim, 0, &elem)

	it.ix++
	return nil
}

func (c *Converter) processBytes(it *Item) error {
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

func (c *Converter) processMap(it *Item) error {
	if it.flag == None && it.ix == 0 {
		it.keys = it.val.MapKeys()

		err := c.write(c.options.MapStart)

		if err != nil {
			return err
		}

		it.flag = KeyNext
	} else if it.flag == KeyNext && it.ix < it.val.Len() {
		err := c.write(c.options.MapSep)

		if err != nil {
			return err
		}
	} else if it.ix == it.val.Len() {
		err := c.write(c.options.MapEnd)

		if err != nil {
			return err
		}

		c.stack.Pop()
		return nil
	}

	key := it.keys[it.ix]

	if it.flag == KeyNext {
		c.push(None, 0, &key)
		it.flag = ValueNext
	} else if it.flag == ValueNext {
		c.writeRune(':')
		val := it.val.MapIndex(key)
		c.push(None, 0, &val)
		it.flag = KeyNext
		it.ix++
	}

	return nil
}

func (c *Converter) processPointer(it *Item) error {
	err := c.write("&")

	if err != nil {
		return err
	}

	elem := it.val.Elem()

	if elem.Kind() == r.Struct {
		it.flag = StructData
	}

	it.val = &elem
	return nil
}

func (c *Converter) processStruct(it *Item) error {
	if it.flag != StructData {
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
		c.push(None, 0, &field)
	} else {
		addr := unsafe.Pointer(field.UnsafeAddr())
		data := r.NewAt(field.Type(), addr).Elem()
		c.push(None, 0, &data)
	}

	it.ix++
	return nil
}

func (c *Converter) processBool(val *r.Value) error {
	return c.write(strconv.FormatBool(val.Bool()))
}

func (c *Converter) processByte(val *r.Value) error {
	return c.writeByte(byte(val.Uint()))
}

func (c *Converter) processChan(val *r.Value) error {
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

func (c *Converter) processComplex(val *r.Value) error {
	realPart := floatToString(real(val.Complex()), 64)
	imagPart := floatToString(imag(val.Complex()), 64)
	return c.write(fmt.Sprintf("%s + %si", realPart, imagPart))
}

func (c *Converter) processFloat(val *r.Value, bitSize int) error {
	return c.write(floatToString(val.Float(), bitSize))
}

func (c *Converter) processFunc(val *r.Value) error {
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

func (c *Converter) processInt(val *r.Value) error {
	return c.write(strconv.FormatInt(val.Int(), 10))
}

func (c *Converter) processInterface(_ *r.Value) error {
	return c.write("interface{}")
}

func (c *Converter) processRune(val *r.Value) error {
	return c.writeRune(rune(val.Int()))
}

func (c *Converter) processString(val *r.Value) error {
	return c.write(val.String())
}

func (c *Converter) processUint(val *r.Value) error {
	return c.write(strconv.FormatUint(val.Uint(), 10))
}

func (c *Converter) processUintptr(val *r.Value) error {
	return c.write(fmt.Sprintf("0x%X", val.Uint()))
}

func (c *Converter) processUnsafe(val *r.Value) error {
	return c.write(fmt.Sprintf("Ux%X", val.Pointer()))
}

func (c *Converter) write(s string) error {
	_, err := c.builder.WriteString(s)
	return err
}

func (c *Converter) writeByte(b byte) error {
	return c.builder.WriteByte(b)
}

func (c *Converter) writeRune(r rune) error {
	_, err := c.builder.WriteRune(r)
	return err
}
