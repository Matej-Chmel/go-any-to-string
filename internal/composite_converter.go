package internal

import (
	r "reflect"
	"strings"
	"unsafe"

	gs "github.com/Matej-Chmel/go-generic-stack"
)

// Processes given stack and writes each item to builder.
// In the end, writes all contents from builder to writer.
type CompositeConverter struct {
	builder strings.Builder
	LeafConverter
	stack gs.Stack[*Item]
}

// Constructs new Converter with val as the first item in the stack
func NewCompositeConverter(o *Options, val *r.Value) CompositeConverter {
	c := CompositeConverter{
		builder:       strings.Builder{},
		LeafConverter: NewLeafConverter(o),
		stack:         gs.Stack[*Item]{},
	}
	c.push(None, 0, val)
	return c
}

// Run the whole conversion from start to finish
func (c *CompositeConverter) ConvertStackToString() string {
	if c.stack.Empty() {
		return ""
	}

	firstItem, _ := c.stack.Top()

	if c.options.ShowType {
		c.write(formatCompositeType(firstItem.val))
		c.writeRune(' ')
	}

	for c.stack.HasItems() {
		top, _ := c.stack.Top()
		c.processItem(top)
	}

	return c.builder.String()
}

// Processes arrays and slices. Items with 2 or 3 dimensions
// are delegated to method processArray2D3D.
func (c *CompositeConverter) processArray(it *Item) {
	var currentDim uint32

	if it.dim == 0 {
		// Count dimensions for the most outer layer
		currentDim = countDimensions(it.val)
		it.SetCurrentDim(currentDim)
		it.SetOriginalDim(currentDim)
	} else {
		currentDim = it.GetCurrentDim()
	}

	// Delegate work for 2D and 3D
	if currentDim >= 2 && currentDim <= 3 {
		c.processArray2D3D(it, currentDim)
		return
	}

	// Compute indentation
	var indentLen int

	if origDim := it.GetOriginalDim(); origDim <= 3 {
		indentLen = 0
	} else {
		indentLen = int(origDim - max(currentDim, 3))
	}

	length := it.val.Len()

	if it.ix == 0 {
		// First item

		c.writeIndent(indentLen)

		if currentDim >= 4 {
			// 4D and higher, write start and newlines
			c.write(c.options.ArrayStart)
			c.write(c.options.ArraySep2D)
		} else if currentDim == 1 && it.flag != InnerDim {
			// 1D standalone, write start
			c.write(c.options.ArrayStart)
		}
	} else if it.ix > 0 && it.ix < length {
		// Items other than first one

		if currentDim >= 4 {
			// 4D and higher, write separators
			c.write(c.options.ArraySep2D)
			c.writeIndent(indentLen)
			c.write(c.options.ArrayEnd)
			c.write(c.options.ArraySep2D)
			c.writeIndent(indentLen)
			c.write(c.options.ArrayStart)
			c.write(c.options.ArraySep2D)
		} else {
			// 1D standalone or inner, only 1 separator
			c.write(c.options.ArraySep)
		}
	} else if it.ix == length {
		// End of the array

		if currentDim >= 4 {
			// 4D and higher, Write separator and end
			c.write(c.options.ArraySep2D)
			c.writeIndent(indentLen)
			c.write(c.options.ArrayEnd)
		} else if currentDim == 1 && it.flag != InnerDim {
			// 1D standalone, write end
			c.write(c.options.ArrayEnd)
		}
	}

	if it.ix == length {
		// End of the array, pop item from stack
		c.stack.Pop()
	} else {
		// Otherwise, push lower layer onto stack
		c.pushArrayItem(it, currentDim)
	}
}

// Processes a 2D and 3D arrays or slices
func (c *CompositeConverter) processArray2D3D(it *Item, currentDim uint32) {
	if l := it.val.Len(); it.ix > 0 && it.ix < l {
		// Items other than first one

		if currentDim == 2 {
			c.write(c.options.ArraySep2D)
		} else {
			c.write(c.options.ArraySep3D)
		}
	} else if it.ix == l {
		// End of the array, pop item from stack
		c.stack.Pop()
		return
	}

	// Push lower layer onto stack
	c.pushArrayItem(it, currentDim)
}

// Process an array or slice of bytes as a string
func (c *CompositeConverter) processBytes(it *Item) {
	l := it.val.Len()

	if it.ix == l {
		c.stack.Pop()
		return
	}

	elem := it.val.Index(it.ix)
	c.push(it.flag, 0, &elem)

	// Move index onto the next byte or rune
	it.ix++
}

// Process composite types.
// Returns true if Item it is of composite type and was processed
func (c *CompositeConverter) processComposites(it *Item, kind r.Kind) bool {
	switch kind {
	case r.Array, r.Slice:
		if it.flag == Bytes || it.flag == Runes {
			// Item already flagged as string
			c.processBytes(it)
			return true
		}

		if it.ix == 0 {
			// Item is at the first stage of processing
			// Kind of underlying element type
			elemKind := it.val.Type().Elem().Kind()

			if c.options.ByteAsString && elemKind == r.Uint8 {
				// Array of bytes that should be processed as a string
				it.flag = Bytes
				c.processBytes(it)
				return true
			}

			if c.options.RuneAsString && elemKind == r.Int32 {
				// Array of runes that should be processed as a string
				it.flag = Runes
				c.processBytes(it)
				return true
			}
		}

		// Process as standard array
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

// Attempts to write custom string representation of a struct
// by finding and calling its String() string method
func (c *CompositeConverter) processCustomMethod(it *Item) bool {
	if method := it.val.MethodByName("String"); method.IsValid() {
		res := method.Call(nil)

		if len(res) == 1 && res[0].Kind() == r.String {
			c.write(res[0].String())
			return true
		}
	}

	return false
}

// If Item it is a part of a byte or rune array, its value is written
// as character and true is returned.
func (c *CompositeConverter) processFlaggedBytes(it *Item) bool {
	if it.flag == Bytes {
		c.write(c.formatByte(it.val))
	} else if it.flag == Runes {
		c.write(c.formatRune(it.val))
	} else {
		return false
	}

	return true
}

// Determines kind of Item it and writes its value into
// a builder. Item may be popped from the stack if fully processed.
func (c *CompositeConverter) processItem(it *Item) {
	kind := it.val.Kind()

	if processed := c.processComposites(it, kind); processed {
		return
	}

	// Item is not a composite, one pass will suffice,
	// pop the item from the stack
	c.stack.Pop()

	if processed := c.processFlaggedBytes(it); processed {
		return
	}

	c.write(c.ConvertToString(it.val))
}

// Processes a map
func (c *CompositeConverter) processMap(it *Item) {
	if it.flag == None && it.ix == 0 {
		// First stage, save keys so that order doesn't change
		it.flag = KeyNext
		it.keys = it.val.MapKeys()
		c.write(c.options.MapStart)
	} else if it.flag == KeyNext && it.ix < it.val.Len() {
		// Write separator between two key-value pairs
		c.write(c.options.MapSepVal)
	} else if it.ix == it.val.Len() {
		// End of map, pop item from the stack
		c.write(c.options.MapEnd)
		c.stack.Pop()
		return
	}

	if key := it.keys[it.ix]; it.flag == KeyNext {
		// Process a key next
		c.push(None, 0, &key)
		it.flag = ValueNext
	} else if it.flag == ValueNext {
		// Write separator between key and value
		c.write(c.options.MapSepKey)

		// Find and process a value next
		val := it.val.MapIndex(key)
		c.push(None, 0, &val)
		it.flag = KeyNext

		// Move index onto the next key-value pair
		it.ix++
	}
}

// Processes a pointer
func (c *CompositeConverter) processPointer(it *Item) {
	elem := it.val.Elem()

	if elem.Kind() == r.Struct {
		// Pointer points to a struct,
		// method processStruct doesn't need to create the pointer.
		it.flag = StructData

		// If the struct has custom String() string, use it
		if c.processCustomMethod(it) {
			c.stack.Pop()
			return
		}
	}

	c.write("&")
	it.val = &elem
}

// Processes a struct
func (c *CompositeConverter) processStruct(it *Item) {
	if it.flag != StructData {
		// First stage, create pointer to the struct
		// so that unexported fields can be addressed
		tmp := r.New(it.val.Type())
		tmp.Elem().Set(*it.val)
		elem := tmp.Elem()
		it.flag = StructData
		it.val = &elem

		// If the struct has custom String() string, use it
		if c.processCustomMethod(it) {
			c.stack.Pop()
			return
		}
	}

	if it.ix == 0 {
		c.write(c.options.StructStart)
	} else if it.ix < it.val.NumField() {
		c.write(c.options.StructSep)
	} else if it.ix == it.val.NumField() {
		// End of struct, pop item from the stack
		c.write(c.options.StructEnd)
		c.stack.Pop()
		return
	}

	if field := it.val.Field(it.ix); field.CanInterface() {
		// Current field is exported, simply process it
		c.push(None, 0, &field)
	} else {
		// Current field is unexported, retrieve its value from a memory address
		addr := unsafe.Pointer(field.UnsafeAddr())
		data := r.NewAt(field.Type(), addr).Elem()
		c.push(None, 0, &data)
	}

	// Move index onto the next field
	it.ix++
}

// Push new Item onto the stack
func (c *CompositeConverter) push(flag uint, index int, val *r.Value) {
	c.stack.Push(NewItem(flag, index, val))
}

// Push the next element from array or slice represented by the Item it
// onto the stack
func (c *CompositeConverter) pushArrayItem(it *Item, currentDim uint32) {
	elem := it.val.Index(it.ix)

	// The next layer is an inner dimension
	newItem := NewItem(InnerDim, 0, &elem)
	newItem.dim = it.dim
	newItem.SetCurrentDim(currentDim - 1)
	c.stack.Push(newItem)

	// Move index onto the next array element
	it.ix++
}

// Write a string to builder
func (c *CompositeConverter) write(s string) {
	c.builder.WriteString(s)
}

// Write a byte to builder
func (c *CompositeConverter) writeByte(b byte) {
	c.builder.WriteByte(b)
}

// Write indentation to builder
func (c *CompositeConverter) writeIndent(length int) {
	for i := 0; i < length; i++ {
		c.write(c.options.ArrayIndent)
	}
}

// Write a rune to builder
func (c *CompositeConverter) writeRune(r rune) {
	c.builder.WriteRune(r)
}
