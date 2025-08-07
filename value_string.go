package dynflags

import "strconv"

// String defines a dynamic string flag with default and help text.
func (g *ConfigGroup) String(name string, def string, usage string) *FlagBuilder[string] {
	ptr := new(string)
	val := NewBaseValue(ptr, def, func(s string) (string, error) { return s, nil }, strconv.Quote)

	flag := &Flag{
		name:       name,
		usage:      usage,
		value:      val,
		defaultSet: def != "",
	}

	g.Flags[name] = flag
	g.flagOrder = append(g.flagOrder, name)

	return &FlagBuilder[string]{
		df:  nil,
		bf:  flag,
		ptr: ptr,
	}
}
