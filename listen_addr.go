package dynflags

import (
	"fmt"
	"net"
)

type ListenAddrValue struct {
	Bound *string
}

func (l *ListenAddrValue) GetBound() any {
	if l.Bound == nil {
		return nil
	}
	return *l.Bound
}

func (l *ListenAddrValue) Parse(value string) (any, error) {
	_, err := net.ResolveTCPAddr("tcp", value)
	if err != nil {
		return nil, fmt.Errorf("invalid listen address: %w", err)
	}
	return &value, nil
}

func (l *ListenAddrValue) Set(value any) error {
	if str, ok := value.(*string); ok {
		*l.Bound = *str
		return nil
	}
	return fmt.Errorf("invalid value type: expected string pointer for listen address")
}

// ListenAddr defines a flag that validates a TCP listen address (host:port or :port).
func (g *ConfigGroup) ListenAddr(name, defaultValue, usage string) *Flag {
	bound := new(string)
	if defaultValue != "" {
		if _, err := net.ResolveTCPAddr("tcp", defaultValue); err != nil {
			panic(fmt.Sprintf("%s has an invalid default listen address '%s': %v", name, defaultValue, err))
		}
		*bound = defaultValue // Copy the parsed ListenAddr into bound
	}
	flag := &Flag{
		Type:    FlagTypeString,
		Default: defaultValue,
		Usage:   usage,
		value:   &ListenAddrValue{Bound: bound},
	}
	g.Flags[name] = flag
	g.flagOrder = append(g.flagOrder, name)
	return flag
}

// GetListenAddr returns the string value of a validated listen address flag.
func (pg *ParsedGroup) GetListenAddr(flagName string) (string, error) {
	value, exists := pg.Values[flagName]
	if !exists {
		return "", fmt.Errorf("flag '%s' not found in group '%s'", flagName, pg.Name)
	}
	if str, ok := value.(string); ok {
		return str, nil
	}
	return "", fmt.Errorf("flag '%s' is not a string listen address", flagName)
}
