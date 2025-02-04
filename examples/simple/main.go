package main

import (
	"fmt"
	"os"

	"github.com/containeroo/dynflags"
)

func main() {
	// 1. Create a new DynFlags instance.
	df := dynflags.New(dynflags.ContinueOnError)

	// Optional metadata for your CLI help text.
	df.Title("My Example CLI")
	df.Description("This application demonstrates how to use dynflags in a simple program.")
	df.Epilog("For more information, visit https://example.com.")

	// 2. Define a group named "app" and add a flag called "msg".
	//    The third parameter here is the default value, and the fourth is the help text.
	appGroup := df.Group("app")
	appGroup.String("msg", "Hello, World!", "Message to be displayed.")

	// For demonstration, we hard-code example arguments.
	// In a real program, you would typically use: args := os.Args[1:]
	args := []string{
		"--app.default.msg", "Hello default DynFlags!",
		"--app.custom.msg", "Hello custom DynFlags!",
	}

	// 3. Parse the command-line arguments.
	if err := df.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	// If no arguments were provided, show usage and exit.
	if len(args) < 2 {
		df.Usage()
		os.Exit(0)
	}

	// 4. Access the parsed values.
	parsedGroups := df.Parsed()

	// Look up the "app" group if it was populated.
	appParsed := parsedGroups.Lookup("app")
	if appParsed == nil {
		fmt.Println("No 'app' flags found.")
		return
	}

	// Iterate over all groups returned by df.Parsed().Groups() (in this case, only "app")
	// and then over each identifier (e.g., "default", "custom") within that group.
	for groupName, identifiers := range parsedGroups.Groups() {
		for identifierName, parsedGroup := range identifiers {
			msg, err := parsedGroup.GetString("msg") // Custom method you may have added
			if err != nil {
				fmt.Printf("Error getting flag 'msg' in group %q (identifier %q): %v\n", groupName, identifierName, err)
				continue
			}
			fmt.Printf("Group %q, identifier %q => msg: %q\n", groupName, identifierName, msg)
		}
	}

	// If any arguments were unrecognized or invalid, print them here.
	unparsed := df.UnknownArgs()
	if len(unparsed) > 0 {
		fmt.Println("Unknown arguments:", unparsed)
	}
}
