package internal

type Options struct {
	// Symbol at the end of an array or slice, default "]"
	ArrayEnd string
	// Indentation symbol used to indent layers in multidimensional array
	// or slice, default 4 spaces
	ArrayIndent string
	// Separator between two elements of an array or slice, default " "
	ArraySep string
	// Separator between two elements of a 2D array or slice, default "\n"
	ArraySep2D string
	// Separator between two elements of a 3D array or slice, default "\n\n"
	ArraySep3D string
	// Symbol at the start of an array or slice, default "["
	ArrayStart string
	// Flag indicating whether a byte array or slice should be written
	// as a string, default false
	ByteAsString bool
	// Maximum number of decimal places to write when processing a floating-point
	// number, default 3
	FloatDecimalPlaces int
	// Symbol at the end of a function's parameter list, default ")"
	FuncEnd string
	// Symbol between two parameters of a function, default ", "
	FuncSep string
	// Symbol between input and output parameter lists of a function, default " "
	FuncSepInOut string
	// Symbol at the start of a function's parameter list, default "("
	FuncStart string
	// Symbol at the start of a map, default "}"
	MapEnd string
	// Symbol between key and value of a map, default ":"
	MapSepKey string
	// Symbol between two key-value pairs of a map, default " "
	MapSepVal string
	// Symbol at the start of a map, default "{"
	MapStart string
	// Flag indicating whether a rune array or slice should be written
	// as a string, default false
	RuneAsString bool
	// Flag indicating whether to write a type name before the final string,
	// default false
	ShowType bool
	// Symbol at the end of a struct, default "}"
	StructEnd string
	// Symbol between two fields of a struct, default " "
	StructSep string
	// Symbol at the start of a struct, default "{"
	StructStart string
}

const (
	// Default symbol at the end of an array or slice
	DefaultArrayEnd string = "]"
	// Default indentation symbol used to indent layers in multidimensional array or slice
	DefaultArrayIndent string = "    "
	// Default separator between two elements of an array or slice
	DefaultArraySep string = " "
	// Default separator between two elements of a 2D array or slice
	DefaultArraySep2D string = "\n"
	// Default separator between two elements of a 3D array or slice
	DefaultArraySep3D string = "\n\n"
	// Default symbol at the start of an array or slice
	DefaultArrayStart string = "["
	// Default flag indicating whether a byte array or slice should be written as a string
	DefaultByteAsString bool = false
	// Default maximum number of decimal places to write when processing a floating-point number
	DefaultFloatDecimalPlaces int = 3
	// Default symbol at the end of a function's parameter list
	DefaultFuncEnd string = ")"
	// Default symbol between two parameters of a function
	DefaultFuncSep string = ", "
	// Default symbol between input and output parameter lists of a function
	DefaultFuncSepInOut string = " "
	// Default symbol at the start of a function's parameter list
	DefaultFuncStart string = "("
	// Default symbol at the end of a map
	DefaultMapEnd string = "}"
	// Default symbol between key and value of a map
	DefaultMapSepKey string = ":"
	// Default symbol between two key-value pairs of a map
	DefaultMapSepVal string = " "
	// Default symbol at the start of a map
	DefaultMapStart string = "{"
	// Default flag indicating whether a rune array or slice should be written as a string
	DefaultRuneAsString bool = false
	// Default flag indicating whether to write a type name before the final string
	DefaultShowType bool = false
	// Default symbol at the end of a struct
	DefaultStructEnd string = "}"
	// Default symbol between two fields of a struct
	DefaultStructSep string = " "
	// Default symbol at the start of a struct
	DefaultStructStart string = "{"
)

// Constructs new Options with default values
func NewOptions() *Options {
	return &Options{
		ArrayEnd:           DefaultArrayEnd,
		ArrayIndent:        DefaultArrayIndent,
		ArraySep:           DefaultArraySep,
		ArraySep2D:         DefaultArraySep2D,
		ArraySep3D:         DefaultArraySep3D,
		ArrayStart:         DefaultArrayStart,
		ByteAsString:       DefaultByteAsString,
		FloatDecimalPlaces: DefaultFloatDecimalPlaces,
		FuncEnd:            DefaultFuncEnd,
		FuncSep:            DefaultFuncSep,
		FuncSepInOut:       DefaultFuncSepInOut,
		FuncStart:          DefaultFuncStart,
		MapEnd:             DefaultMapEnd,
		MapSepKey:          DefaultMapSepKey,
		MapSepVal:          DefaultMapSepVal,
		MapStart:           DefaultMapStart,
		RuneAsString:       DefaultRuneAsString,
		ShowType:           DefaultShowType,
		StructEnd:          DefaultStructEnd,
		StructSep:          DefaultStructSep,
		StructStart:        DefaultStructStart,
	}
}
