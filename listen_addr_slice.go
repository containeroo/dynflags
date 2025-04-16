package dynflags

import (
	"fmt"
	"net"
	"strings"
)

type ListenAddrSlicesValue struct {
	Bound *[]string
}

func (s *ListenAddrSlicesValue) GetBound() any {
	if s.Bound == nil {
		return nil
	}
	return *s.Bound
}

func (s *ListenAddrSlicesValue) Parse(value string) (any, error) {
	_, err := net.ResolveTCPAddr("tcp", value)
	if err != nil {
		return nil, fmt.Errorf("invalid listen address: %w", err)
	}
	return value, nil
}

func (s *ListenAddrSlicesValue) Set(value any) error {
	if addr, ok := value.(string); ok {
		*s.Bound = append(*s.Bound, addr)
		return nil
	}
	return fmt.Errorf("invalid value type: expected string listen address")
}

// ListenAddrSlices defines a slice-of-listen-address flag with the specified name, default values, and usage.
func (g *ConfigGroup) ListenAddrSlices(name string, value []string, usage string) *Flag {
	bound := &value
	defaultValue := strings.Join(value, ",")

	// Validate all default addresses
	for _, v := range value {
		if _, err := net.ResolveTCPAddr("tcp", v); err != nil {
			panic(fmt.Sprintf("%s has an invalid default listen address '%s': %v", name, v, err))
		}
	}

	flag := &Flag{
		Type:    FlagTypeStringSlice,
		Default: defaultValue,
		Usage:   usage,
		value:   &ListenAddrSlicesValue{Bound: bound},
	}
	g.Flags[name] = flag
	g.flagOrder = append(g.flagOrder, name)
	return flag
}

// GetListenAddrSlices returns the []string value of a listen address slice flag.
func (pg *ParsedGroup) GetListenAddrSlices(flagName string) ([]string, error) {
	value, exists := pg.Values[flagName]
	if !exists {
		return nil, fmt.Errorf("flag '%s' not found in group '%s'", flagName, pg.Name)
	}

	if list, ok := value.([]string); ok {
		return list, nil
	}

	if str, ok := value.(string); ok {
		return []string{str}, nil
	}

	return nil, fmt.Errorf("flag '%s' is not a []string listen address slice", flagName)
}
