package dynflags

// Flag represents a dynamic flag definition (type, value, metadata, etc.)
type Flag struct {
	name       string
	usage      string
	value      Value
	required   bool
	deprecated string
	metavar    string
	defaultSet bool
}

// Value is implemented by all concrete flag value holders.
type Value interface {
	Set(string) error // parses and sets from string
	Get() any         // returns the parsed value
	Default() string  // stringified default
	IsChanged() bool  // whether the value was explicitly set
}

// BoolFlag marks a flag as --flag (true) shorthand.
type BoolFlag interface {
	IsBoolFlag() bool
}

// SliceFlag marks slice flags for internal classification.
type SliceFlag interface {
	isSlice()
}

// DelimiterSetter is implemented by slice values that support custom delimiters.
type DelimiterSetter interface {
	SetDelimiter(string)
}
