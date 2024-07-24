package internal

import (
	"fmt"
	r "reflect"
	"runtime"
	"strconv"
	"strings"
)

type LeafConverter struct {
	options *Options
}

func NewLeafConverter(o *Options) LeafConverter {
	return LeafConverter{options: o}
}

// Converts basic built-in types to string
func (c *LeafConverter) ConvertToString(val *r.Value) string {
	switch val.Kind() {
	case r.Bool:
		return c.formatBool(val)
	case r.Chan:
		return c.formatChannel(val)
	case r.Complex64:
		return c.formatComplex(32, val)
	case r.Complex128:
		return c.formatComplex(64, val)
	case r.Float32:
		return c.formatFloat(val, 32)
	case r.Float64:
		return c.formatFloat(val, 64)
	case r.Func:
		return c.formatFunc(val)

	case r.Int32:
		if c.options.RuneAsString {
			// Item is part of array of runes
			// that should be written as characters
			return c.formatRune(val)
		} else {
			return c.formatInt(val)
		}

	case r.Int, r.Int8, r.Int16, r.Int64:
		return c.formatInt(val)
	case r.Interface, r.Invalid:
		return c.formatInterface()
	case r.Uint, r.Uint16, r.Uint32, r.Uint64:
		return c.formatUint(val)

	case r.Uint8:
		if c.options.ByteAsString {
			// Item is part of array of bytes
			// that should be written as characters
			return c.formatByte(val)
		} else {
			return c.formatUint(val)
		}

	case r.String:
		return c.formatString(val)
	case r.Uintptr:
		return c.formatUintptr(val)
	case r.UnsafePointer:
		return c.formatUnsafe(val)
	}

	return val.Kind().String()
}

// Formats a bool value
func (c *LeafConverter) formatBool(val *r.Value) string {
	return strconv.FormatBool(val.Bool())
}

// Formats a byte
func (c *LeafConverter) formatByte(val *r.Value) string {
	return string(byte(val.Uint()))
}

// Formats a channel
func (c *LeafConverter) formatChannel(val *r.Value) string {
	return "chan " + val.Type().Elem().String()
}

// Formats a complex number
func (c *LeafConverter) formatComplex(bitSize int, val *r.Value) string {
	complex := val.Complex()
	realPart := floatToString(bitSize, c.options.FloatDecimalPlaces, real(complex))
	imagPart := floatToString(bitSize, c.options.FloatDecimalPlaces, imag(complex))
	return fmt.Sprintf("%s + %si", realPart, imagPart)
}

// Formats a floating-point number
func (c *LeafConverter) formatFloat(val *r.Value, bitSize int) string {
	return floatToString(bitSize, c.options.FloatDecimalPlaces, val.Float())
}

// Formats a function signature
func (c *LeafConverter) formatFunc(val *r.Value) string {
	// Determine the function's name
	name := runtime.FuncForPC(val.Pointer()).Name()
	parts := strings.Split(name, ".")
	lastIx := len(parts) - 1
	funcName := parts[lastIx]

	// Write the function's name
	var builder strings.Builder
	builder.WriteString(funcName)
	builder.WriteString(c.options.FuncStart)

	// Write types of input parameters
	aType := val.Type()
	in := aType.NumIn()

	for i := 0; i < in; i++ {
		if i > 0 {
			builder.WriteString(c.options.FuncSep)
		}

		builder.WriteString(aType.In(i).String())
	}

	// Write end bracket and separator between
	// input and output parameters
	builder.WriteString(c.options.FuncEnd)
	builder.WriteString(c.options.FuncSepInOut)
	out := aType.NumOut()

	if out > 1 {
		// Multiple output parameters are enclosed in brackets
		builder.WriteString(c.options.FuncStart)
	}

	// Write types of output parameters
	for i := 0; i < out; i++ {
		if i > 0 {
			builder.WriteString(c.options.FuncSep)
		}

		builder.WriteString(aType.Out(i).String())
	}

	if out > 1 {
		builder.WriteString(c.options.FuncEnd)
	}

	return builder.String()
}

// Formats a signed integer
func (c *LeafConverter) formatInt(val *r.Value) string {
	return strconv.FormatInt(val.Int(), 10)
}

// Formats an empty interface{}
func (c *LeafConverter) formatInterface() string {
	return "interface{}"
}

// Formats a rune as a character
func (c *LeafConverter) formatRune(val *r.Value) string {
	return string(rune(val.Int()))
}

// Formats a string
func (c *LeafConverter) formatString(val *r.Value) string {
	return val.String()
}

// Formats an unsinged integer
func (c *LeafConverter) formatUint(val *r.Value) string {
	return strconv.FormatUint(val.Uint(), 10)
}

// Formats an unsigned pointer
func (c *LeafConverter) formatUintptr(val *r.Value) string {
	return fmt.Sprintf("0x%X", val.Uint())
}

// Formats an unsafe pointer
func (c *LeafConverter) formatUnsafe(val *r.Value) string {
	return fmt.Sprintf("Ux%X", val.Pointer())
}
