package dynflags

import (
	"fmt"
	"time"
)

type DurationValue struct {
	Bound *time.Duration
}

func (d *DurationValue) GetBound() any {
	if d.Bound == nil {
		return nil
	}
	return *d.Bound
}

func (d *DurationValue) Parse(value string) (any, error) {
	return time.ParseDuration(value)
}

func (d *DurationValue) Set(value any) error {
	if dur, ok := value.(time.Duration); ok {
		*d.Bound = dur
		return nil
	}
	return fmt.Errorf("invalid value type: expected duration")
}

// Duration defines a duration flag with the specified name, default value, and usage description.
// The flag is added to the group's flag list and returned as a *Flag instance.
func (g *ConfigGroup) Duration(name string, value time.Duration, usage string) *Flag {
	bound := &value
	flag := &Flag{
		Type:    FlagTypeDuration,
		Default: value,
		Usage:   usage,
		value:   &DurationValue{Bound: bound},
	}
	g.Flags[name] = flag
	g.flagOrder = append(g.flagOrder, name)
	return flag
}

// GetDuration returns the time.Duration value of a flag with the given name
func (pg *ParsedGroup) GetDuration(flagName string) (time.Duration, error) {
	vaue, exists := pg.Values[flagName]
	if !exists {
		return 0, fmt.Errorf("flag '%s' not found in group '%s'", flagName, pg.Name)
	}
	if durationVal, ok := vaue.(time.Duration); ok {
		return durationVal, nil
	}
	return 0, fmt.Errorf("flag '%s' is not a time.Duration", flagName)
}
