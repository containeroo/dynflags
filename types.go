package dynflags

import "fmt"

// ConfigGroup represents a static group of dynamic flags (e.g., "http", "db", etc.)
type ConfigGroup struct {
	Name      string
	Flags     map[string]*Flag
	flagOrder []string
	usage     string
}

// Lookup returns a flag by name (e.g., "port", "enabled").
func (g *ConfigGroup) Lookup(name string) *Flag {
	return g.Flags[name]
}

// ParsedGroup holds the resolved values for a specific group + identifier combo.
type ParsedGroup struct {
	Parent *ConfigGroup   // link to the static definition
	Name   string         // identifier (e.g. "main", "us-west", "1")
	Values map[string]any // parsed values keyed by flag name
}

// GroupsMap maps group names to a set of named identifiers.
type GroupsMap map[string]IdentifiersMap

// IdentifiersMap maps dynamic identifiers to parsed groups.
type IdentifiersMap map[string]*ParsedGroup

// Get returns the typed value of a flag.
func GetAs[T any](pg *ParsedGroup, name string) (T, error) {
	val, ok := pg.Values[name]
	if !ok {
		var zero T
		return zero, fmt.Errorf("flag %q not found in group %q", name, pg.Name)
	}
	v, ok := val.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("flag %q is not of expected type", name)
	}
	return v, nil
}
