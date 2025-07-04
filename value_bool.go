package dynflags

import "strconv"

// Bool defines a dynamic bool flag with a given name, default, and usage.
func (g *ConfigGroup) Bool(name string, def bool, usage string) *FlagBuilder[bool] {
	ptr := new(bool)
	val := NewBaseValue(ptr, def, strconv.ParseBool, strconv.FormatBool)

	flag := &Flag{
		name:       name,
		usage:      usage,
		value:      val,
		defaultSet: true, // impossible to not set
	}

	g.Flags[name] = flag
	g.flagOrder = append(g.flagOrder, name)

	return &FlagBuilder[bool]{
		df:  nil, // will be injected later via .Group() or .Env() call
		bf:  flag,
		ptr: ptr,
	}
}
