package xlsx2csv

type config struct {
	// Align empty fields on the end to header length (see with_empty_fields).
	// Default: false.
	align bool

	// A func that returns a sheet for conversation to csv.
	// Default: the first one in the document.
	getSheet SheetSelector

	// A csv delimiter.
	// Default: comma
	comma rune
}

var defaultConfig = config{
	align:    false,
	getSheet: FirstSheet(),
	comma:    ',',
}

type Option func(*config)

func WithAlign() Option {
	return func(c *config) {
		c.align = true
	}
}

func SetSheetSelector(selector SheetSelector) Option {
	return func(c *config) {
		c.getSheet = selector
	}
}

func SetComma(comma rune) Option {
	return func(c *config) {
		c.comma = comma
	}
}
