package goanytostring

type Options struct {
	ByteAsString bool
	RuneAsString bool
}

func NewOptions() Options {
	return Options{
		ByteAsString: false,
		RuneAsString: false,
	}
}
