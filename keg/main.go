package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	flag "github.com/ogier/pflag"
)

var (
	onlyFail  = flag.BoolP("onlyfail", "f", false, "Only display failed tests")
	steinHost = flag.StringP("stein", "s", "", "Address of Stein server to send results to")
	project   = flag.StringP("project", "p", "default", "Project the test results belong to")
	testType  = flag.StringP("type", "t", "", "Test result type")
	quiet     = flag.BoolP("quiet", "q", false, "Quiet mode")
)

func main() {
	flag.Parse()
	if len(*steinHost) > 0 {
		if err := validatePost(*steinHost); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	buffer := &bytes.Buffer{}
	_, err := newKegParser().Parse(io.TeeReader(os.Stdin, buffer))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing tap input:", err)
		os.Exit(1)
	}
	if len(*steinHost) > 0 {
		if id, err := postSuite(buffer, *steinHost, *project, *testType); err != nil {
			fmt.Fprintln(os.Stderr, "Error posting test results to Stein:", err)
			os.Exit(1)
		} else if *steinHost != "" {
			if !*quiet {
				fmt.Println("\nTest finished, results can be viewed at:")
			}
			fmt.Printf("http://%s/#/projects/%s/tests/%s\n", *steinHost, *project, id)
		}
	}
}
