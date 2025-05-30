package dynflags

import (
	"fmt"
	"strconv"
)

type BoolValue struct {
	Bound *bool
}

func (b *BoolValue) GetBound() any {
	if b.Bound == nil {
		return nil
	}
	return *b.Bound
}

func (b *BoolValue) Parse(value string) (any, error) {
	return strconv.ParseBool(value)
}

func (b *BoolValue) Set(value any) error {
	if val, ok := value.(bool); ok {
		*b.Bound = val
		return nil
	}
	return fmt.Errorf("invalid value type: expected bool")
}

// Bool defines a boolean flag with the specified name, default value, and usage description.
// The flag is added to the group's flag list and returned as a *Flag instance.
func (g *ConfigGroup) Bool(name string, value bool, usage string) *Flag {
	bound := &value
	flag := &Flag{
		Type:    FlagTypeBool,
		Default: value,
		Usage:   usage,
		value:   &BoolValue{Bound: bound},
	}
	g.Flags[name] = flag
	g.flagOrder = append(g.flagOrder, name)
	return flag
}

// GetBool returns the bool value of a flag with the given name
func (pg *ParsedGroup) GetBool(flagName string) (bool, error) {
	value, exists := pg.Values[flagName]
	if !exists {
		return false, fmt.Errorf("flag '%s' not found in group '%s'", flagName, pg.Name)
	}
	if boolVal, ok := value.(bool); ok {
		return boolVal, nil
	}
	return false, fmt.Errorf("flag '%s' is not a bool", flagName)
}
