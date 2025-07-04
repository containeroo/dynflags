package dynflags

// FlagBuilder is the base builder for scalar and slice dynamic flags.
type FlagBuilder[T any] struct {
	df  *DynFlags
	bf  *Flag
	ptr *T
}

// Required marks the flag as required.
func (b *FlagBuilder[T]) Required() *FlagBuilder[T] {
	b.bf.required = true
	return b
}

// Metavar sets the metavar used in help output for this flag.
func (b *FlagBuilder[T]) Metavar(s string) *FlagBuilder[T] {
	b.bf.metavar = s
	return b
}

// Deprecated marks the flag as deprecated with a reason.
func (b *FlagBuilder[T]) Deprecated(reason string) *FlagBuilder[T] {
	b.bf.deprecated = reason
	return b
}

// SliceFlagBuilder extends FlagBuilder with slice-specific options.
type SliceFlagBuilder[T any] struct {
	FlagBuilder[T]
}

// Delimiter sets the delimiter used to split string input for slice flags.
func (b *SliceFlagBuilder[T]) Delimiter(sep string) *SliceFlagBuilder[T] {
	if d, ok := b.bf.value.(DelimiterSetter); ok {
		d.SetDelimiter(sep)
	}
	return b
}
