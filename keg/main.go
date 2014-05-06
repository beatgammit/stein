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
)

func main() {
	flag.Parse()
	if len(*steinHost) > 0 {
		if err := validatePost(*steinHost); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	buffer := &bytes.Buffer{}
	_, err := newKegParser().Parse(io.TeeReader(os.Stdin, buffer))
	if err != nil {
		fmt.Println("Error parsing tap input:", err)
		os.Exit(1)
	}
	if len(*steinHost) > 0 {
		if err := postSuite(buffer, *steinHost, *project, *testType); err != nil {
			fmt.Println("Error posting test results to Stein:", err)
			os.Exit(1)
		}
	}
}
