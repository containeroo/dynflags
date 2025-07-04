package dynflags_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/containeroo/dynflags"
	"github.com/stretchr/testify/assert"

	flag "github.com/spf13/pflag"
)

func TestDynflags(t *testing.T) {
	t.Run("Smoke test", func(t *testing.T) {
		df := dynflags.New("test.exe", dynflags.ContinueOnError)

		appGroup := df.Group("app")
		appGroup.String("msg", "Hello, World!", "Message to be displayed.")

		args := []string{
			"--app.default.msg", "Hello default DynFlags!",
			"--app.custom.msg", "Hello custom DynFlags!",
		}

		if err := df.Parse(args); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
			os.Exit(1)
		}

		defaultGroup := df.Parsed("app", "default")
		customGroup := df.Parsed("app", "custom")

		msg1, _ := dynflags.GetAs[string](defaultGroup, "msg")
		msg2, _ := dynflags.GetAs[string](customGroup, "msg")

		fmt.Println("Default:", msg1)
		fmt.Println("Custom:", msg2)
	})

	t.Run("extended smoke tests", func(t *testing.T) {
		df := dynflags.New("test.exe", dynflags.ContinueOnError)

		// HTTP group
		http := df.Group("http")
		http.String("name", "", "Name of the HTTP checker")
		http.String("method", "GET", "HTTP method to use")
		http.String("address", "", "HTTP target URL")
		http.Bool("allow-duplicate-headers", false, "Allow duplicate HTTP headers")
		http.String("expected-status-codes", "200", "Expected HTTP status codes")
		http.Bool("skip-tls-verify", false, "Skip TLS verification")

		// ICMP group
		icmp := df.Group("icmp")
		icmp.String("name", "", "Name of the ICMP checker")
		icmp.String("address", "", "ICMP target address")

		// TCP group
		tcp := df.Group("tcp")
		tcp.String("name", "", "Name of the TCP checker")
		tcp.String("address", "", "TCP target address")

		args := []string{
			"--http.default.name", "HTTP Checker Default",
			"--http.default.address", "default.com:80",
			"--http.other.name", "HTTP Checker Other",
			"--http.other.address", "other.com:443",
			"--http.other.method", "POST",
			"--icmp.custom.address", "8.8.4.4",
			"--tcp.testing.address", "example.com:443",
			"--default-interval=5s", // pflag
		}

		err := df.Parse(args)
		assert.NoError(t, err)

		fmt.Println("\n=== Dynamic Flags ===")
		for groupName, identifiers := range df.Groups() {
			fmt.Printf("Group: %s\n", groupName)
			for identifier, pg := range identifiers {
				t.Logf(" Identifier: %s\n", identifier)
				for flagKey, value := range pg.Values {
					t.Logf("    %s: %v\n", flagKey, value)
				}
			}
		}
	})

	t.Run("help tests", func(t *testing.T) {
		buf := strings.Builder{}

		fs := flag.NewFlagSet("test.exe", flag.ContinueOnError)
		fs.SetOutput(&buf)
		fs.BoolP("help", "h", false, "Show help.")

		df := dynflags.New("test.exe", dynflags.ContinueOnError)
		df.Title("\nsome dynamic flags:")
		df.SetOutput(&buf)

		fs.Usage = func() {
			out := fs.Output() // capture writer ONCE

			fmt.Fprintf(out, "Usage: %s [FLAGS] [DYNAMIC FLAGS..]\n", strings.ToLower(fs.Name())) // nolint:errcheck

			fmt.Fprintln(out, "\nGlobal Flags:") // nolint:errcheck
			fs.SetOutput(out)
			fs.PrintDefaults()

			df.PrintTitle(out)
			df.PrintDescription(out, 80)
			df.PrintDefaults()
			df.PrintEpilog(out, 80)
		}

		// HTTP group
		http := df.Group("http")
		http.String("name", "", "Name of the HTTP checker").Metavar("TESTVAR")
		http.String("method", "GET", "HTTP method to use")
		http.String("address", "", "HTTP target URL").Required()
		http.Bool("allow-duplicate-headers", false, "Allow duplicate HTTP headers")
		http.String("expected-status-codes", "200", "Expected HTTP status codes")
		http.Bool("skip-tls-verify", false, "Skip TLS verification").Deprecated("Use --insecure-skip-tls-verify instead")

		// ICMP group
		icmp := df.Group("icmp")
		icmp.String("name", "", "Name of the ICMP checker")
		icmp.String("address", "", "ICMP target address").Required()

		// TCP group
		tcp := df.Group("tcp")
		tcp.String("name", "", "Name of the TCP checker")
		tcp.String("address", "", "TCP target address").Required()

		args := []string{
			"--help",
		}

		err := df.Parse(args)
		assert.NoError(t, err)
		unknownArgs := df.UnknownArgs()

		err = fs.Parse(unknownArgs)
		assert.NoError(t, err)
		help := fs.Lookup("help")
		if help != nil && help.Value.String() == "true" {
			buf.Reset() // instead of creating a new one
			fs.SetOutput(&buf)
			df.SetOutput(&buf)
			fs.Usage()

			expected := `Usage: test.exe [FLAGS] [DYNAMIC FLAGS..]

Global Flags:
  -h, --help   Show help.

some dynamic flags:
HTTP
  Flag                                                             Usage
  --http.<IDENTIFIER>.name TESTVAR                                 Name of the HTTP checker
  --http.<IDENTIFIER>.method METHOD                                HTTP method to use (default: "GET")
  --http.<IDENTIFIER>.address ADDRESS                              HTTP target URL [required]
  --http.<IDENTIFIER>.allow-duplicate-headers                      Allow duplicate HTTP headers (default: false)
  --http.<IDENTIFIER>.expected-status-codes EXPECTED-STATUS-CODES  Expected HTTP status codes (default: "200")
  --http.<IDENTIFIER>.skip-tls-verify                              Skip TLS verification (default: false) [deprecated: Use --insecure-skip-tls-verify instead]

ICMP
  Flag                                 Usage
  --icmp.<IDENTIFIER>.name NAME        Name of the ICMP checker
  --icmp.<IDENTIFIER>.address ADDRESS  ICMP target address [required]

TCP
  Flag                                Usage
  --tcp.<IDENTIFIER>.name NAME        Name of the TCP checker
  --tcp.<IDENTIFIER>.address ADDRESS  TCP target address [required]

`

			assert.Equal(t, expected, buf.String())
		}
	})
}
