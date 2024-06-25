package internal

type Options struct {
	ArrayStart   string
	ArrayEnd     string
	ArraySep     string
	ArraySep2D   string
	ArraySep3D   string
	ByteAsString bool
	FuncStart    string
	FuncEnd      string
	FuncSep      string
	FuncSepInOut string
	DimIndent    string
	MapStart     string
	MapEnd       string
	MapSepKey    string
	MapSepVal    string
	RuneAsString bool
	ShowType     bool
	StructStart  string
	StructEnd    string
	StructSep    string
}

func NewOptions() Options {
	return Options{
		ArrayStart:   "[",
		ArrayEnd:     "]",
		ArraySep:     " ",
		ArraySep2D:   "\n",
		ArraySep3D:   "\n\n",
		ByteAsString: false,
		DimIndent:    "    ",
		FuncStart:    "(",
		FuncEnd:      ")",
		FuncSep:      ", ",
		FuncSepInOut: " ",
		MapStart:     "{",
		MapEnd:       "}",
		MapSepKey:    ":",
		MapSepVal:    " ",
		RuneAsString: false,
		ShowType:     false,
		StructStart:  "{",
		StructEnd:    "}",
		StructSep:    " ",
	}
}
