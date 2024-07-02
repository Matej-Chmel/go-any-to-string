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
	builder  strings.Builder
	debugStr string
	options  *Options
	stack    gs.Stack[*Item]
	writer   io.Writer
}

func NewConverter(o *Options, val *r.Value, writer io.Writer) Converter {
	c := Converter{options: o, stack: gs.Stack[*Item]{}, writer: writer}
	c.push(Top, 0, val)
	return c
}

func countDimensions(val *r.Value) (d uint32) {
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

func (c *Converter) displayType(it *Item) {
	if it.flag == Top {
		it.flag = None

		if c.options.ShowType {
			c.write(formatType(it.val.Type(), true))
		}
	}
}

func (c *Converter) processArray(it *Item) {
	var currentDim uint32

	if it.dim == 0 {
		currentDim = countDimensions(it.val)
		it.SetCurrentDim(currentDim)
		it.SetOriginalDim(currentDim)
	} else {
		currentDim = it.GetCurrentDim()
	}

	if currentDim >= 2 && currentDim <= 3 {
		c.processArray2D3D(it, currentDim)
		return
	}

	var indent int

	if origDim := it.GetOriginalDim(); origDim <= 3 {
		indent = 0
	} else {
		indent = int(origDim - max(currentDim, 3))
	}

	length := it.val.Len()

	if it.ix == 0 {
		for i := 0; i < indent; i++ {
			c.write(c.options.DimIndent)
		}

		if it.flag != InnerDim {
			c.write(c.options.ArrayStart)

			if currentDim > 3 {
				c.write(c.options.ArraySep2D)
			}
		}
	} else if it.ix > 0 && it.ix < length {
		var sep string

		if currentDim > 3 {
			sep = c.options.ArraySep2D
		} else {
			sep = c.options.ArraySep
		}

		c.write(sep)

		if currentDim > 3 {
			for i := 0; i < indent; i++ {
				c.write(c.options.DimIndent)
			}

			c.write(c.options.ArrayEnd)
			c.write(c.options.ArraySep2D)

			for i := 0; i < indent; i++ {
				c.write(c.options.DimIndent)
			}

			c.write(c.options.ArrayStart)
			c.write(c.options.ArraySep2D)
		}
	} else if it.ix == length {
		if currentDim > 3 {
			c.write(c.options.ArraySep2D)

			for i := 0; i < indent; i++ {
				c.write(c.options.DimIndent)
			}
		}

		if it.flag != InnerDim {
			c.write(c.options.ArrayEnd)
		}
	}

	if it.ix == length {
		c.stack.Pop()
	} else {
		c.pushArrayItem(it, currentDim)
	}
}

func (c *Converter) processArray2D3D(it *Item, currentDim uint32) {
	if l := it.val.Len(); it.ix > 0 && it.ix < l {
		var sep string

		if currentDim == 2 {
			sep = c.options.ArraySep2D
		} else {
			sep = c.options.ArraySep3D
		}

		c.write(sep)
	} else if it.ix == l {
		c.stack.Pop()
		return
	}

	c.pushArrayItem(it, currentDim)
}

func (c *Converter) processBytes(it *Item) {
	l := it.val.Len()

	if it.ix == l {
		c.stack.Pop()
		return
	}

	elem := it.val.Index(it.ix)
	c.push(it.flag, 0, &elem)

	it.ix++
}

func (c *Converter) processComposites(it *Item, kind r.Kind) bool {
	switch kind {
	case r.Array, r.Slice:
		if it.flag == Bytes || it.flag == Runes {
			c.processBytes(it)
			return true
		}

		if it.ix == 0 {
			elemKind := it.val.Type().Elem().Kind()

			if c.options.ByteAsString && elemKind == r.Uint8 {
				it.flag = Bytes
				c.processBytes(it)
				return true
			}

			if c.options.RuneAsString && elemKind == r.Int32 {
				it.flag = Runes
				c.processBytes(it)
				return true
			}
		}

		c.processArray(it)
	case r.Map:
		c.processMap(it)
	case r.Pointer:
		c.processPointer(it)
	case r.Struct:
		c.processStruct(it)
	default:
		return false
	}

	return true
}

func (c *Converter) processFlaggedBytes(it *Item) bool {
	if it.flag == Bytes {
		c.processByte(it.val)
		return true
	}

	if it.flag == Runes {
		c.processRune(it.val)
		return true
	}

	return false
}

func (c *Converter) processLeaf(it *Item, kind r.Kind) {
	switch kind {
	case r.Bool:
		c.processBool(it.val)
	case r.Chan:
		c.processChan(it.val)
	case r.Complex64, r.Complex128:
		c.processComplex(it.val)
	case r.Float32:
		c.processFloat(it.val, 32)
	case r.Float64:
		c.processFloat(it.val, 64)
	case r.Func:
		c.processFunc(it.val)

	case r.Int32:
		if c.options.RuneAsString {
			c.processRune(it.val)
		} else {
			c.processInt(it.val)
		}

	case r.Int, r.Int8, r.Int16, r.Int64:
		c.processInt(it.val)
	case r.Interface, r.Invalid:
		c.processInterface(it.val)
	case r.Uint, r.Uint16, r.Uint32, r.Uint64:
		c.processUint(it.val)

	case r.Uint8:
		if c.options.ByteAsString {
			c.processByte(it.val)
		} else {
			c.processUint(it.val)
		}

	case r.String:
		c.processString(it.val)
	case r.Uintptr:
		c.processUintptr(it.val)
	case r.UnsafePointer:
		c.processUnsafe(it.val)
	default:
		c.write(kind.String())
	}
}

func (c *Converter) processItem(it *Item) {
	c.displayType(it)
	kind := it.val.Kind()

	if processed := c.processComposites(it, kind); processed {
		return
	}

	c.stack.Pop()

	if processed := c.processFlaggedBytes(it); processed {
		return
	}

	c.processLeaf(it, kind)
}

func (c *Converter) processMap(it *Item) {
	if it.flag == None && it.ix == 0 {
		it.flag = KeyNext
		it.keys = it.val.MapKeys()
		c.write(c.options.MapStart)
	} else if it.flag == KeyNext && it.ix < it.val.Len() {
		c.write(c.options.MapSepVal)
	} else if it.ix == it.val.Len() {
		c.write(c.options.MapEnd)
		c.stack.Pop()
		return
	}

	key := it.keys[it.ix]

	if it.flag == KeyNext {
		c.push(None, 0, &key)
		it.flag = ValueNext
	} else if it.flag == ValueNext {
		c.write(c.options.MapSepKey)
		val := it.val.MapIndex(key)
		c.push(None, 0, &val)
		it.flag = KeyNext
		it.ix++
	}
}

func (c *Converter) processPointer(it *Item) {
	elem := it.val.Elem()

	if elem.Kind() == r.Struct {
		it.flag = StructData

		if c.processStructCustom(it) {
			c.stack.Pop()
			return
		}
	}

	c.write("&")
	it.val = &elem
}

func (c *Converter) processStruct(it *Item) {
	if it.flag != StructData {
		tmp := r.New(it.val.Type())
		tmp.Elem().Set(*it.val)
		elem := tmp.Elem()
		it.flag = StructData
		it.val = &elem

		if c.processStructCustom(it) {
			c.stack.Pop()
			return
		}
	}

	if it.ix == 0 {
		c.write(c.options.StructStart)
	} else if it.ix < it.val.NumField() {
		c.write(c.options.StructSep)
	} else if it.ix == it.val.NumField() {
		c.write(c.options.StructEnd)
		c.stack.Pop()
		return
	}

	if field := it.val.Field(it.ix); field.CanInterface() {
		c.push(None, 0, &field)
	} else {
		addr := unsafe.Pointer(field.UnsafeAddr())
		data := r.NewAt(field.Type(), addr).Elem()
		c.push(None, 0, &data)
	}

	it.ix++
}

func (c *Converter) processStructCustom(it *Item) bool {
	if method := it.val.MethodByName("String"); method.IsValid() {
		res := method.Call(nil)

		if len(res) == 1 && res[0].Kind() == r.String {
			c.write(res[0].String())
			return true
		}
	}

	return false
}

func (c *Converter) processBool(val *r.Value) {
	c.write(strconv.FormatBool(val.Bool()))
}

func (c *Converter) processByte(val *r.Value) {
	c.writeByte(byte(val.Uint()))
}

func (c *Converter) processChan(val *r.Value) {
	c.write("chan ")
	c.write(val.Type().Elem().String())
}

func floatToString(f float64, bitSize int) string {
	s := strconv.FormatFloat(f, 'f', 3, bitSize)
	s = strings.TrimRight(s, "0")
	return strings.TrimRight(s, ".")
}

func (c *Converter) processComplex(val *r.Value) {
	realPart := floatToString(real(val.Complex()), 64)
	imagPart := floatToString(imag(val.Complex()), 64)
	c.write(fmt.Sprintf("%s + %si", realPart, imagPart))
}

func (c *Converter) processFloat(val *r.Value, bitSize int) {
	c.write(floatToString(val.Float(), bitSize))
}

func (c *Converter) processFunc(val *r.Value) {
	typ := val.Type()
	in, out := typ.NumIn(), typ.NumOut()

	name := runtime.FuncForPC(val.Pointer()).Name()
	parts := strings.Split(name, ".")
	lastIx := len(parts) - 1
	funcName := parts[lastIx]

	c.write(funcName)
	c.write(c.options.FuncStart)

	for i := 0; i < in; i++ {
		if i > 0 {
			c.write(c.options.FuncSep)
		}

		c.write(typ.In(i).String())
	}

	c.write(c.options.FuncEnd)
	c.write(c.options.FuncSepInOut)

	if out > 1 {
		c.write(c.options.FuncStart)
	}

	for i := 0; i < out; i++ {
		if i > 0 {
			c.write(c.options.FuncSep)
		}

		c.write(typ.Out(i).String())
	}

	if out > 1 {
		c.write(c.options.FuncEnd)
	}
}

func (c *Converter) processInt(val *r.Value) {
	c.write(strconv.FormatInt(val.Int(), 10))
}

func (c *Converter) processInterface(_ *r.Value) {
	c.write("interface{}")
}

func (c *Converter) processRune(val *r.Value) {
	c.writeRune(rune(val.Int()))
}

func (c *Converter) processString(val *r.Value) {
	c.write(val.String())
}

func (c *Converter) processUint(val *r.Value) {
	c.write(strconv.FormatUint(val.Uint(), 10))
}

func (c *Converter) processUintptr(val *r.Value) {
	c.write(fmt.Sprintf("0x%X", val.Uint()))
}

func (c *Converter) processUnsafe(val *r.Value) {
	c.write(fmt.Sprintf("Ux%X", val.Pointer()))
}

func (c *Converter) push(flag int, index int, val *r.Value) {
	c.stack.Push(NewItem(flag, index, val))
}

func (c *Converter) pushArrayItem(it *Item, currentDim uint32) {
	elem := it.val.Index(it.ix)
	newItem := NewItem(InnerDim, 0, &elem)
	newItem.dim = it.dim
	newItem.SetCurrentDim(currentDim - 1)
	c.stack.Push(newItem)
	it.ix++
}

func (c *Converter) Run() error {
	for c.stack.HasItems() {
		top, _ := c.stack.Top()
		c.processItem(top)
	}

	_, err := c.writer.Write([]byte(c.builder.String()))
	return err
}

func (c *Converter) write(s string) {
	c.builder.WriteString(s)

	if len(c.debugStr) > 12 {
		c.debugStr = ""
	}

	c.debugStr += s
}

func (c *Converter) writeByte(b byte) {
	c.builder.WriteByte(b)
}

func (c *Converter) writeRune(r rune) {
	c.builder.WriteRune(r)
}
