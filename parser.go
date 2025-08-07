package dynflags

import (
	"fmt"
	"strings"
)

// Parse parses the provided CLI arguments.
func (df *DynFlags) Parse(args []string) error {
	for i := 0; i < len(args); i++ {
		arg := args[i]

		key, val, err := df.extractKeyValue(arg, args, &i)
		if err != nil {
			if df.parseBehavior == ExitOnError {
				return err
			}
			df.unparsedArgs = append(df.unparsedArgs, arg)
			continue
		}

		groupName, ident, flagName, err := df.splitKey(key)
		if err != nil {
			if df.parseBehavior == ExitOnError {
				return err
			}
			df.unparsedArgs = append(df.unparsedArgs, arg)
			continue
		}

		err = df.setFlag(groupName, ident, flagName, val)
		if err != nil {
			if df.parseBehavior == ExitOnError {
				return err
			}
			df.unparsedArgs = append(df.unparsedArgs, arg)
		}
	}
	return nil
}

// extractKeyValue supports --key=value or --key value syntax.
func (df *DynFlags) extractKeyValue(arg string, args []string, i *int) (key, value string, err error) {
	if !strings.HasPrefix(arg, "--") {
		return "", "", fmt.Errorf("invalid flag: %s", arg)
	}

	arg = strings.TrimPrefix(arg, "--")
	if strings.Contains(arg, "=") {
		parts := strings.SplitN(arg, "=", 2)
		return parts[0], parts[1], nil
	}

	// try next argument
	if *i+1 < len(args) && !strings.HasPrefix(args[*i+1], "--") {
		*i++
		return arg, args[*i], nil
	}

	return "", "", fmt.Errorf("missing value for --%s", arg)
}

// splitKey expects group.identifier.flag pattern.
func (df *DynFlags) splitKey(full string) (group, ident, flag string, err error) {
	parts := strings.Split(full, ".")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid flag key: --%s (must be --<group>.<identifier>.<flag>)", full)
	}
	return parts[0], parts[1], parts[2], nil
}

// setFlag resolves and sets the value of a known flag.
func (df *DynFlags) setFlag(groupName, ident, flagName, val string) error {
	group, ok := df.groups[groupName]
	if !ok {
		return fmt.Errorf("group %q not defined", groupName)
	}

	flag := group.Lookup(flagName)
	if flag == nil {
		return fmt.Errorf("flag %q not found in group %q", flagName, groupName)
	}

	err := flag.value.Set(val)
	if err != nil {
		return fmt.Errorf("failed to parse --%s.%s.%s: %w", groupName, ident, flagName, err)
	}

	pg := df.getParsedGroup(group, ident)
	pg.Values[flagName] = flag.value.Get()
	return nil
}

// getParsedGroup initializes or retrieves a parsed group for the identifier.
func (df *DynFlags) getParsedGroup(group *ConfigGroup, ident string) *ParsedGroup {
	if _, ok := df.parsed[group.Name]; !ok {
		df.parsed[group.Name] = make(IdentifiersMap)
	}
	if pg, ok := df.parsed[group.Name][ident]; ok {
		return pg
	}

	pg := &ParsedGroup{
		Parent: group,
		Name:   ident,
		Values: make(map[string]any),
	}
	df.parsed[group.Name][ident] = pg
	return pg
}
