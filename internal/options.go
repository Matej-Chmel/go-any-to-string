package internal

type Options struct {
	ArrayStart   string
	ArrayEnd     string
	ArraySep     string
	ArraySep2D   string
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
		ArraySep2D:   "\n",
		ByteAsString: false,
		MapStart:     "{",
		MapEnd:       "}",
		MapSep:       " ",
		RuneAsString: false,
		ShowType:     false,
	}
}
