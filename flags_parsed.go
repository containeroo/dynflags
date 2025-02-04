package dynflags

// GroupsMap is a map of group name -> IdentifiersMap.
type GroupsMap map[string]IdentifiersMap

// IdentifiersMap is a map of identifier -> ParsedGroup pointer.
type IdentifiersMap map[string]*ParsedGroup

// ParsedGroup represents a runtime group with parsed values.
type ParsedGroup struct {
	Parent *ConfigGroup           // Reference to the parent static group.
	Name   string                 // Identifier for the child group (e.g., "IDENTIFIER1").
	Values map[string]interface{} // Parsed values for the group's flags.
}

// Lookup retrieves the value of a flag in the parsed group.
func (g *ParsedGroup) Lookup(flagName string) interface{} {
	if g == nil {
		return nil
	}
	return g.Values[flagName]
}

// ParsedGroups represents all parsed groups with lookup and iteration support.
type ParsedGroups struct {
	groups GroupsMap // Nested map of group name -> IdentifiersMap.
}

// Lookup retrieves a group by its name.
func (g *ParsedGroups) Lookup(groupName string) *ParsedIdentifiers {
	if g == nil {
		return nil
	}
	if identifiers, exists := g.groups[groupName]; exists {
		return &ParsedIdentifiers{Name: groupName, identifiers: identifiers}
	}
	return nil
}

// Groups returns the underlying GroupsMap for direct iteration.
func (g *ParsedGroups) Groups() GroupsMap {
	return g.groups
}

// ParsedIdentifiers provides lookup for identifiers within a group.
type ParsedIdentifiers struct {
	Name        string
	identifiers IdentifiersMap
}

// Lookup retrieves a specific identifier within a group.
func (i *ParsedIdentifiers) Lookup(identifier string) *ParsedGroup {
	if i == nil {
		return nil
	}
	return i.identifiers[identifier]
}

// Parsed returns a ParsedGroups instance for the dynflags instance.
func (f *DynFlags) Parsed() *ParsedGroups {
	parsed := make(GroupsMap)
	for groupName, groups := range f.parsedGroups {
		identifierMap := make(IdentifiersMap)
		for _, group := range groups {
			identifierMap[group.Name] = group
		}
		parsed[groupName] = identifierMap
	}
	return &ParsedGroups{groups: parsed}
}
