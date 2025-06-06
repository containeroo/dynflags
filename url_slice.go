package dynflags

import (
	"fmt"
	"net/url"
	"strings"
)

type URLSlicesValue struct {
	Bound *[]*url.URL
}

func (u *URLSlicesValue) GetBound() any {
	if u.Bound == nil {
		return nil
	}
	return *u.Bound
}

func (u *URLSlicesValue) Parse(value string) (any, error) {
	parsedURL, err := url.Parse(value)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %s, error: %w", value, err)
	}
	return parsedURL, nil
}

func (u *URLSlicesValue) Set(value any) error {
	if parsedURL, ok := value.(*url.URL); ok {
		*u.Bound = append(*u.Bound, parsedURL)
		return nil
	}
	return fmt.Errorf("invalid value type: expected *url.URL")
}

// URLSlices defines a URL slice flag with the specified name, default value, and usage description.
// The flag is added to the group's flag list and returned as a *Flag instance.
func (g *ConfigGroup) URLSlices(name string, value []*url.URL, usage string) *Flag {
	bound := &value
	defaultValue := make([]string, len(value))
	for i, u := range value {
		defaultValue[i] = u.String()
	}

	flag := &Flag{
		Type:    FlagTypeURLSlice,
		Default: strings.Join(defaultValue, ","),
		Usage:   usage,
		value:   &URLSlicesValue{Bound: bound},
	}
	g.Flags[name] = flag
	g.flagOrder = append(g.flagOrder, name)
	return flag
}

// GetURLSlices returns the []*url.URL value of a flag with the given name
func (pg *ParsedGroup) GetURLSlices(flagName string) ([]*url.URL, error) {
	value, exists := pg.Values[flagName]
	if !exists {
		return nil, fmt.Errorf("flag '%s' not found in group '%s'", flagName, pg.Name)
	}

	if urlSlice, ok := value.([]*url.URL); ok {
		return urlSlice, nil
	}

	if u, ok := value.(*url.URL); ok {
		return []*url.URL{u}, nil
	}

	return nil, fmt.Errorf("flag '%s' is not a []*url.URL", flagName)
}
