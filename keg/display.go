package main

import (
	"fmt"
	"strings"

	"github.com/beatgammit/stein"
	"github.com/wsxiaoys/terminal"
)

const (
	timeFormat = "2006-01-02 at 15:04:05"
)

func newKegParser() *stein.Parser {
	return &stein.Parser{
		&kegHandler{
			currentLevel: -1,
		},
	}
}

type kegHandler struct {
	*stein.DefaultHandler
	currentLevel int
	afterCase    bool
}

func (h *kegHandler) getIndentation(level int) string {
	indentation := ""
	for i := 0; i < level; i++ {
		indentation += "    "
	}
	return indentation
}

func (h *kegHandler) HandleSuite(s *stein.Suite) {
	fmt.Printf("Test suite started on %s\n", s.Start.Local().Format(timeFormat))
	if s.Count > 0 {
		fmt.Printf("Running %d tests...\n\n", s.Count)
	}
}

func (h *kegHandler) HandleCase(c *stein.Case) {
	clearLine()
	fmt.Printf("%s%s\n", h.getIndentation(h.currentLevel), c.Label)
	h.currentLevel = c.Level
}

func (h *kegHandler) HandleTest(t *stein.Test) {
	status := t.Status
	if strings.ToLower(status) == "omit" {
		status = "skip"
	}
	clearLine()
	testLine := fmt.Sprintf("%s%s ... %s", h.getIndentation(h.currentLevel+1), status, t.Label)
	printColor(statusColors[status], testLine)
	if *onlyFail {
		if status == "fail" || status == "error" {
			fmt.Print("\n")
		} else {
			fmt.Print("\r")
		}
	} else {
		fmt.Print("\n")
	}
}

func (h *kegHandler) HandleFinal(t *stein.Tally) {
	clearLine()
	fmt.Printf("\nResults: %d pass, %d fail, %d error, %d skip\n",
		t.Counts.Pass, t.Counts.Fail, t.Counts.Error, t.Counts.Omit)
}

func clearLine() {
	if *onlyFail {
		terminal.Stdout.ClearLine()
	}
}

var statusColors = map[string]string{
	"pass":  "g",
	"skip":  "y",
	"fail":  "r",
	"error": "m",
	"todo":  "b",
}

func printColor(color, text string) {
	terminal.Stdout.Color(color).Print(text).Reset()
}
