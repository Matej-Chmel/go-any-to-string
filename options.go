package goanytostring

import ite "github.com/Matej-Chmel/go-any-to-string/internal"

// Type alias for formatting options
type Options = ite.Options

// Constructs new formatting options
func NewOptions() *Options {
	return ite.NewOptions()
}
