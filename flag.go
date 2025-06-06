package dynflags

type FlagType string

const (
	FlagTypeStringSlice   FlagType = "..STRINGs"
	FlagTypeString        FlagType = "STRING"
	FlagTypeInt           FlagType = "INT"
	FlagTypeIntSlice      FlagType = "..INTs"
	FlagTypeBool          FlagType = "BOOL"
	FlagTypeBoolSlice     FlagType = "..BOOLs"
	FlagTypeDuration      FlagType = "DURATION"
	FlagTypeDurationSlice FlagType = "..DURATIONs"
	FlagTypeFloat         FlagType = "FLOAT"
	FlagTypeFloatSlice    FlagType = "..FLOATs"
	FlagTypeIP            FlagType = "IP"
	FlagTypeIPSlice       FlagType = "..IPs"
	FlagTypeURL           FlagType = "URL"
	FlagTypeURLSlice      FlagType = "..URLs"
)

// Flag represents a single configuration flag
type Flag struct {
	Default any       // Default value for the flag
	Type    FlagType  // Type of the flag
	Usage   string    // Description for usage
	metaVar string    // MetaVar for flag
	value   FlagValue // Encapsulated parsing and value-setting logic
}

func (f *Flag) MetaVar(metaVar string) {
	f.metaVar = metaVar
}

// FlagValue interface encapsulates parsing and value-setting logic
type FlagValue interface {
	// Parse parses the given string value into the flag's value type
	Parse(value string) (any, error)
	// Set sets the flag's value to the given value
	Set(value any) error
	// GetBound returns the bound value of the flag.
	GetBound() any
}

// Value returns the current value of the flag.
func (f *Flag) GetValue() any {
	if f == nil || f.value == nil {
		return nil
	}
	return f.value.GetBound()
}
