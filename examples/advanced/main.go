package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/containeroo/dynflags"
	flag "github.com/spf13/pflag"
)

// Example version string
var version = "v1.2.3"

// HelpRequested is a custom error to indicate the user requested help or version info.
type HelpRequested struct {
	Message string
}

func (e *HelpRequested) Error() string {
	return e.Message
}

// Is checks if the target is a HelpRequested error.
func (e *HelpRequested) Is(target error) bool {
	_, ok := target.(*HelpRequested)
	return ok
}

// parseFlags orchestrates parsing both global flags (with pflag)
// and dynamic flags (with dynflags). It returns any error that
// indicates parsing failed, or HelpRequested for help/version output.
func parseFlags(args []string, version string, output io.Writer) (*flag.FlagSet, *dynflags.DynFlags, error) {
	// 1. Setup the global pflag FlagSet
	flagSet := setupGlobalFlags()
	flagSet.SetOutput(output) // direct usage output here if needed

	// 2. Setup dynamic flags
	dynFlags := setupDynamicFlags()
	dynFlags.SetOutput(output)
	dynFlags.SortFlags = true

	// 3. Provide custom usage that prints both global and dynamic flags
	setupUsage(flagSet, dynFlags)

	// 4. Parse the dynamic flags first so we can separate known vs. unknown arguments
	if err := dynFlags.Parse(args); err != nil {
		return nil, nil, fmt.Errorf("error parsing dynamic flags: %w", err)
	}

	// Unknown arguments might be pflag or truly unrecognized
	unknownArgs := dynFlags.UnknownArgs()

	// 5. Parse pflag (the known global flags)
	if err := flagSet.Parse(unknownArgs); err != nil {
		return nil, nil, fmt.Errorf("error parsing global flags: %w", err)
	}

	// 6. Handle special flags for help and version
	if err := handleSpecialFlags(flagSet, version); err != nil {
		return nil, nil, err
	}

	return flagSet, dynFlags, nil
}

// setupGlobalFlags defines global flags (pflag) for the application
func setupGlobalFlags() *flag.FlagSet {
	flagSet := flag.NewFlagSet("advancedExample", flag.ContinueOnError)
	flagSet.SortFlags = false

	// Some generic global flags:
	flagSet.Bool("version", false, "Show version and exit.")
	flagSet.BoolP("help", "h", false, "Show help.")
	flagSet.Duration("default-interval", 2*time.Second, "Default interval between checks.")
	return flagSet
}

// setupDynamicFlags defines dynamic flags for HTTP, ICMP, and TCP as an example
func setupDynamicFlags() *dynflags.DynFlags {
	dyn := dynflags.New(dynflags.ContinueOnError)
	dyn.Epilog("For more information, see https://github.com/containeroo/dynflags")
	dyn.SortGroups = true
	dyn.SortFlags = true

	// HTTP group
	http := dyn.Group("http")
	http.String("name", "", "Name of the HTTP checker")
	http.String("method", "GET", "HTTP method to use")
	http.String("address", "", "HTTP target URL")
	http.Duration("interval", 1*time.Second, "Time between HTTP requests (overrides --default-interval if set)")
	http.StringSlices("header", nil, "HTTP headers to send")
	http.Bool("allow-duplicate-headers", false, "Allow duplicate HTTP headers")
	http.String("expected-status-codes", "200", "Expected HTTP status codes")
	http.Bool("skip-tls-verify", false, "Skip TLS verification")
	http.Duration("timeout", 2*time.Second, "Timeout for HTTP requests")

	// ICMP group
	icmp := dyn.Group("icmp")
	icmp.String("name", "", "Name of the ICMP checker")
	icmp.String("address", "", "ICMP target address")
	icmp.Duration("interval", 1*time.Second, "Time between ICMP requests (overrides --default-interval if set)")
	icmp.Duration("read-timeout", 2*time.Second, "Timeout for ICMP read")
	icmp.Duration("write-timeout", 2*time.Second, "Timeout for ICMP write")

	// TCP group
	tcp := dyn.Group("tcp")
	tcp.String("name", "", "Name of the TCP checker")
	tcp.String("address", "", "TCP target address")
	tcp.Duration("interval", 1*time.Second, "Time between TCP requests (overrides --default-interval if set)")
	tcp.Duration("timeout", 2*time.Second, "Timeout for TCP connection")

	return dyn
}

// setupUsage sets a custom usage function that prints both pflag and dynflags usage.
func setupUsage(flagSet *flag.FlagSet, dynFlags *dynflags.DynFlags) {
	flagSet.Usage = func() {
		out := flagSet.Output() // capture writer ONCE

		fmt.Fprintf(out, "Usage: %s [GLOBAL FLAGS...] [DYNAMIC FLAGS...]\n", flagSet.Name()) // nolint:errcheck
		fmt.Fprintln(out, "\nGlobal Flags:")                                                 // nolint:errcheck
		flagSet.PrintDefaults()

		fmt.Fprintln(out, "\nDynamic Flags:") // nolint:errcheck
		dynFlags.PrintDefaults()
	}
}

// handleSpecialFlags checks if --help or --version was requested.
func handleSpecialFlags(flagSet *flag.FlagSet, versionStr string) error {
	helpFlag := flagSet.Lookup("help")
	if helpFlag != nil && helpFlag.Value.String() == "true" {
		// Capture usage output into a buffer, then return a HelpRequested error.
		buffer := &bytes.Buffer{}
		flagSet.SetOutput(buffer)
		flagSet.Usage()
		return &HelpRequested{Message: buffer.String()}
	}

	versionFlag := flagSet.Lookup("version")
	if versionFlag != nil && versionFlag.Value.String() == "true" {
		return &HelpRequested{Message: fmt.Sprintf("%s version %s\n", flagSet.Name(), versionStr)}
	}

	return nil
}

// main is our entry point, showing how to parse and then use the flags.
func main() {
	// We'll pretend our arguments are from the CLI; replace os.Args[1:] in real usage.
	args := []string{
		"--http.default.name", "HTTP Checker Default",
		"--http.default.address", "default.com:80",
		"--http.other.name", "HTTP Checker Other",
		"--http.other.address", "other.com:443",
		"--http.other.method", "POST",
		"--icmp.custom.address", "8.8.4.4",
		"--tcp.testing.address", "example.com:443",
		"--default-interval=5s", // pflag
		//		"--unknownArg", "someValue", // see how it's handled by dynflags
	}

	// Optionally redirect usage/errors to something other than os.Stderr if desired
	output := os.Stdout

	// Parse everything
	flagSet, dynFlags, err := parseFlags(args, version, output)
	if err != nil {
		// If the user requested help or version, print the message and exit
		var hr *HelpRequested
		if errors.As(err, &hr) {
			fmt.Fprint(output, hr.Message) //nolint:errcheck
			return
		}

		fmt.Fprintf(output, "Failed to parse flags: %v\n", err) // nolint:errcheck
		os.Exit(1)
	}

	// If we got here, parse succeeded. Let's show what we got.

	// 1. Print global flags
	fmt.Println("=== Global Flags ===")
	defaultInterval, _ := flagSet.GetDuration("default-interval")
	fmt.Printf("default-interval: %v\n", defaultInterval)

	// 2. Print dynamic flags
	parsedGroups := dynFlags.Parsed()

	fmt.Println("\n=== Dynamic Flags ===")
	for groupName, identifiers := range parsedGroups.Groups() {
		fmt.Printf("Group: %s\n", groupName)
		for identifier, pg := range identifiers {
			fmt.Printf("  Identifier: %s\n", identifier)
			for flagKey, value := range pg.Values {
				fmt.Printf("    %s: %v\n", flagKey, value)
			}
		}
	}

	fmt.Println("\nDone!")
}
