package internal

import (
	"strconv"
	"strings"
)

// Format a floating-point number as a string.
// The maximum number of decimal places can be specified.
// For example, decimalPlaces = 3, will yield:
// 0.1115 -> 0.112, 0.1 -> 0.1, 1.0 -> 1, 0 -> 0.0
func floatToString(bitSize int, decimalPlaces int, val float64) string {
	s := strconv.FormatFloat(val, 'f', decimalPlaces, bitSize)

	if s = strings.TrimRight(s, ".0"); len(s) == 0 {
		return "0.0"
	}

	return s
}
