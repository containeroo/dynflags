package dynflags

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
)

type FlagPrintMode int

const (
	PrintShort FlagPrintMode = iota // Only short flags: -p
	PrintLong                       // Only long flags: --port
	PrintBoth                       // Both: -p|--port
)

// PrintDefaults prints the help text for all dynamic groups and flags.
func (df *DynFlags) PrintDefaults() {
	out := df.Output()
	w := tabwriter.NewWriter(out, 0, 8, 2, ' ', 0)

	groupOrder := df.groupOrder
	if df.sortGroups {
		groupOrder = append([]string(nil), groupOrder...) // shallow copy
		sort.Strings(groupOrder)
	}

	for _, groupName := range groupOrder {
		group := df.groups[groupName]
		if group == nil {
			continue
		}
		df.printGroup(w, groupName, group)
	}

	w.Flush() // nolint:errcheck
}

func (df *DynFlags) PrintTitle(w io.Writer) {
	if df.title != "" {
		fmt.Fprintln(w, df.title) // nolint:errcheck
	}
}

func (df *DynFlags) PrintDescription(w io.Writer, width int) {
	if df.desc != "" {
		fmt.Fprintln(w, wrapText(df.desc, width)) // nolint:errcheck

		fmt.Fprintln(w) // nolint:errcheck
	}
}

func (df *DynFlags) PrintUsage(w io.Writer) {
	fmt.Fprint(w, "Usage: "+df.Name()) // nolint:errcheck
}

func (df *DynFlags) PrintEpilog(w io.Writer, width int) {
	if df.epilog != "" {
		fmt.Fprintln(w)
		fmt.Fprintln(w, wrapText(df.epilog, width)) // nolint:errcheck

	}
}

func (df *DynFlags) printGroup(w io.Writer, groupName string, group *ConfigGroup) {
	// Header
	if group.usage != "" {
		fmt.Fprintln(w, group.usage) // nolint:errcheck
	} else {
		fmt.Fprintf(w, "%s\n", strings.ToUpper(groupName)) // nolint:errcheck
	}

	flagOrder := group.flagOrder
	if df.sortFlags {
		flagOrder = append([]string(nil), flagOrder...)
		sort.Strings(flagOrder)
	}

	if len(flagOrder) == 0 {
		fmt.Fprintln(w, "  (no flags)") // nolint:errcheck
		return
	}

	fmt.Fprintln(w, "  Flag\tUsage") // nolint:errcheck

	for _, flagName := range flagOrder {
		flag := group.Flags[flagName]
		if flag == nil {
			continue
		}
		df.printFlag(w, groupName, flag)
	}

	fmt.Fprintln(w) // nolint:errcheck
}

func (df *DynFlags) printFlag(w io.Writer, group string, flag *Flag) {
	name := fmt.Sprintf("  --%s.<IDENTIFIER>.%s", group, flag.name)
	if meta := formatMetavar(flag); meta != "" {
		name += " " + meta
	}

	var usageParts []string
	if flag.usage != "" {
		usageParts = append(usageParts, flag.usage)
	}

	if flag.defaultSet && !(flag.required && flag.value.Default() == "") {
		usageParts = append(usageParts, fmt.Sprintf("(default: %s)", flag.value.Default()))
	}

	if flag.required {
		usageParts = append(usageParts, "[required]")
	}

	if flag.deprecated != "" {
		usageParts = append(usageParts, "[deprecated: "+flag.deprecated+"]")
	}

	fmt.Fprintf(w, "%s\t%s\n", name, strings.Join(usageParts, " ")) // nolint:errcheck
}

func formatMetavar(flag *Flag) string {
	if bf, ok := flag.value.(BoolFlag); ok && bf.IsBoolFlag() {
		return ""
	}

	meta := flag.metavar
	if meta == "" {
		meta = strings.ToUpper(flag.name)
	}
	if _, ok := flag.value.(SliceFlag); ok {
		meta += "..."
	}
	return meta
}
