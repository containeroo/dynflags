package dynflags

import (
	"fmt"
	"strings"
)

// BaseSliceValue is a generic value holder for slice flags (e.g., []int, []string).
type BaseSliceValue[T any] struct {
	ptr       *[]T                    // target storage
	def       []T                     // default slice
	changed   bool                    // true if Set was called
	parse     func(string) (T, error) // element parser
	format    func(T) string          // element formatter
	delimiter string                  // input split separator
}

// NewBaseSliceValue constructs a new BaseSliceValue.
func NewBaseSliceValue[T any](
	ptr *[]T,
	def []T,
	parseFn func(string) (T, error),
	formatFn func(T) string,
	delimiter string,
) *BaseSliceValue[T] {
	*ptr = append([]T(nil), def...) // defensive copy
	return &BaseSliceValue[T]{
		ptr:       ptr,
		def:       def,
		parse:     parseFn,
		format:    formatFn,
		delimiter: delimiter,
	}
}

// Set parses and appends new elements from input string.
func (v *BaseSliceValue[T]) Set(s string) error {
	if !v.changed {
		*v.ptr = nil // clear default only on first use
	}
	items := strings.Split(s, v.delimiter)
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		val, err := v.parse(trimmed)
		if err != nil {
			return fmt.Errorf("invalid value %q: %w", trimmed, err)
		}
		*v.ptr = append(*v.ptr, val)
	}
	v.changed = true
	return nil
}

func (v *BaseSliceValue[T]) Get() any {
	return *v.ptr
}

func (v *BaseSliceValue[T]) Default() string {
	formatted := make([]string, len(v.def))
	for i, item := range v.def {
		formatted[i] = v.format(item)
	}
	return strings.Join(formatted, v.delimiter)
}

func (v *BaseSliceValue[T]) IsChanged() bool {
	return v.changed
}

// Implements SliceFlag marker interface.
func (v *BaseSliceValue[T]) isSlice() {}

// SetDelimiter allows custom input delimiters like ":" or ";".
func (v *BaseSliceValue[T]) SetDelimiter(d string) {
	v.delimiter = d
}
