package dynflags

// BaseValue is a generic value holder for scalar dynamic flags.
type BaseValue[T any] struct {
	ptr        *T                      // target storage
	def        T                       // default value
	changed    bool                    // whether Set was called
	parse      func(string) (T, error) // parsing logic
	format     func(T) string          // string formatter
	boolMarker bool                    // identifies bool flags
}

// NewBaseValue constructs a new BaseValue with parser and formatter.
func NewBaseValue[T any](
	ptr *T,
	def T,
	parseFn func(string) (T, error),
	formatFn func(T) string,
) *BaseValue[T] {
	*ptr = def
	return &BaseValue[T]{
		ptr:        ptr,
		def:        def,
		parse:      parseFn,
		format:     formatFn,
		changed:    false,
		boolMarker: isBoolPointer(ptr),
	}
}

func (v *BaseValue[T]) Set(s string) error {
	val, err := v.parse(s)
	if err != nil {
		return err
	}
	*v.ptr = val
	v.changed = true
	return nil
}

func (v *BaseValue[T]) Get() any {
	return *v.ptr
}

func (v *BaseValue[T]) Default() string {
	return v.format(v.def)
}

func (v *BaseValue[T]) IsChanged() bool {
	return v.changed
}

// IsBoolFlag marks flags as --flag shorthand = true, if type is bool
func (v *BaseValue[T]) IsBoolFlag() bool {
	return v.boolMarker
}

// isBoolPointer determines if the provided pointer is *bool
func isBoolPointer(ptr any) bool {
	_, ok := ptr.(*bool)
	return ok
}
