package dynflags

import (
	"io"
	"os"
)

// ParseBehavior controls how parsing errors are handled.
type ParseBehavior int

const (
	ContinueOnError ParseBehavior = iota
	ExitOnError
	PanicOnError
)

// DynFlags manages dynamic groups and their flags.
type DynFlags struct {
	parseBehavior ParseBehavior

	groups     map[string]*ConfigGroup // All defined groups
	groupOrder []string                // Preserve registration order

	parsed GroupsMap // Parsed values
	output io.Writer // Help output destination
	Usage  func()    // Usage callback
	name   string    // Optional program name
	title  string    // Optional help title
	desc   string    // Optional help description
	epilog string    // Optional help footer

	unparsedArgs []string // Unknown args for later inspection

	sortGroups bool
	sortFlags  bool
}

// New creates a new DynFlags instance.
func New(name string, behavior ParseBehavior) *DynFlags {
	df := &DynFlags{
		name:          name,
		parseBehavior: behavior,
		groups:        make(map[string]*ConfigGroup),
		parsed:        make(GroupsMap),
		output:        os.Stdout,
	}
	df.Usage = func() {
		out := df.Output()
		df.PrintTitle(out)
		df.PrintUsage(out)
		df.PrintDescription(out, 80)
		df.PrintDefaults()
		df.PrintEpilog(out, 80)
	}
	return df
}

// Name returns the program name (for usage header).
func (df *DynFlags) Name() string {
	return df.name
}

// Title sets the optional help title.
func (df *DynFlags) Title(title string) {
	df.title = title
}

// Description sets the optional help description.
func (df *DynFlags) Description(desc string) {
	df.desc = desc
}

// Epilog sets the help footer.
func (df *DynFlags) Epilog(epilog string) {
	df.epilog = epilog
}

// SetOutput overrides the help output destination.
func (df *DynFlags) SetOutput(w io.Writer) {
	df.output = w
}

// Output returns the current help output writer.
func (df *DynFlags) Output() io.Writer {
	return df.output
}

// UnknownArgs returns any unparsed CLI arguments.
func (df *DynFlags) UnknownArgs() []string {
	return df.unparsedArgs
}

// SortGroups enables/disables sorting of group names.
func (df *DynFlags) SortGroups(enable bool) {
	df.sortGroups = enable
}

// SortFlags enables/disables sorting of flags within a group.
func (df *DynFlags) SortFlags(enable bool) {
	df.sortFlags = enable
}

// Group returns an existing group or creates a new one.
func (df *DynFlags) Group(name string) *ConfigGroup {
	if g, ok := df.groups[name]; ok {
		return g
	}
	group := &ConfigGroup{
		Name:  name,
		Flags: make(map[string]*Flag),
	}
	df.groups[name] = group
	df.groupOrder = append(df.groupOrder, name)
	return group
}

// Groups returns all parsed values.
func (df *DynFlags) Groups() GroupsMap {
	return df.parsed
}

// Parsed returns the parsed result for a specific group and identifier.
func (df *DynFlags) Parsed(groupName, identifier string) *ParsedGroup {
	if idMap, ok := df.parsed[groupName]; ok {
		return idMap[identifier]
	}
	return nil
}
