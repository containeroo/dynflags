package dynflags

import (
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"
)

// PrintDefaults prints all registered flags
func (df *DynFlags) PrintDefaults() {
	w := tabwriter.NewWriter(df.output, 0, 8, 2, ' ', 0)

	// Print title if present
	if df.title != "" {
		fmt.Fprintln(df.output, df.title) // nolint:errcheck
		fmt.Fprintln(df.output)           // nolint:errcheck
	}

	// Print description if present
	if df.description != "" {
		fmt.Fprintln(df.output, df.description) // nolint:errcheck
		fmt.Fprintln(df.output)                 // nolint:errcheck
	}

	// Sort group names
	if df.SortGroups {
		sort.Strings(df.groupOrder)
	}

	// Iterate over groups in the order they were added
	for _, groupName := range df.groupOrder {
		group := df.configGroups[groupName]

		// Print group usage or fallback to uppercase group name
		if group.usage != "" {
			fmt.Fprintln(w, group.usage) // nolint:errcheck
		} else {
			fmt.Fprintln(w, strings.ToUpper(groupName)) // nolint:errcheck
		}

		// Sort flag names
		if df.SortFlags {
			sort.Strings(group.flagOrder)
		}

		// Print flags for the group
		if len(group.flagOrder) > 0 {
			fmt.Fprintln(w, "  Flag\tUsage") // nolint:errcheck
			for _, flagName := range group.flagOrder {
				flag := group.Flags[flagName]
				usage := flag.Usage
				if flag.Default != nil && flag.Default != "" {
					usage = fmt.Sprintf("%s (default: %v)", flag.Usage, flag.Default)
				}
				metavar := string(flag.Type)
				if flag.metaVar != "" {
					metavar = flag.metaVar
				}

				fmt.Fprintf(w, "  --%s.<IDENTIFIER>.%s %s\t%s\n", groupName, flagName, metavar, usage) // nolint:errcheck
			}
			fmt.Fprintln(w, "") // nolint:errcheck
		}
	}

	// tabwriter buffers output for alignment; flush now to ensure aligned flag output is printed before the epilog
	w.Flush() // nolint:errcheck

	// Print epilog if present
	if df.epilog != "" {
		fmt.Fprintln(df.output)            // nolint:errcheck
		fmt.Fprintln(df.output, df.epilog) // nolint:errcheck
	}
}
