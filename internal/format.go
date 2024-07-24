package internal

import "strconv"

// Format a floating-point number as a string.
// Trailing zeros will be trimmed except one zero
// to denote that the variable is a floating-point number.
// For example, decimalPlaces = 3, will yield:
// 0.1115 -> 0.112, 0.1 -> 0.1
// However, 1 will display as 1.0 only if addZero = true.
func floatToString(
	addZero bool, bitSize int, decimalPlaces int, val float64,
) string {
	s := strconv.FormatFloat(val, 'f', decimalPlaces, bitSize)
	return trimFloat(addZero, s)
}

// Trims trailing zeros except the last one
func trimFloat(addZero bool, data string) string {
	dot := -1
	nonzero := -1

	for i := len(data) - 1; i > 0; i-- {
		if data[i] == '.' {
			dot = i
			break
		} else if nonzero < 0 && data[i] != '0' {
			nonzero = i
		}
	}

	if dot < 0 {
		if addZero {
			return data + ".0"
		}

		return data
	}

	if nonzero < 0 {
		if addZero {
			return data[:dot] + ".0"
		}

		return data[:dot]
	}

	return data[:nonzero+1]
}
