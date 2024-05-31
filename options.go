package goanytostring

type Options struct {
	ArrayStart   string
	ArrayEnd     string
	ArraySep     string
	ByteAsString bool
	MapStart     string
	MapEnd       string
	MapSep       string
	RuneAsString bool
	ShowType     bool
}

func NewOptions() Options {
	return Options{
		ArrayStart:   "[",
		ArrayEnd:     "]",
		ArraySep:     " ",
		ByteAsString: false,
		MapStart:     "{",
		MapEnd:       "}",
		MapSep:       " ",
		RuneAsString: false,
		ShowType:     false,
	}
}
