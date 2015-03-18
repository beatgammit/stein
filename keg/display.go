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
	if *quiet {
		return
	}
	fmt.Printf("Test suite started on %s\n", s.Start.Local().Format(timeFormat))
	if s.Count > 0 {
		fmt.Printf("Running %d tests...\n\n", s.Count)
	}
}

func (h *kegHandler) HandleCase(c *stein.Case) {
	if *quiet {
		return
	}
	clearLine()
	fmt.Printf("%s%s\n", h.getIndentation(h.currentLevel), c.Label)
	h.currentLevel = c.Level
}

func (h *kegHandler) HandleTest(t *stein.Test) {
	if *quiet {
		return
	}
	status := t.Status
	if strings.ToLower(status) == "omit" {
		status = "skip"
	}
	clearLine()
	indent := h.getIndentation(h.currentLevel + 1)
	var testLine string
	var postIndent = ""
	for len(status)+len(postIndent) < longestStatus {
		postIndent += " "
	}
	testLine = fmt.Sprintf("%s%s%s ... %s", indent, status, postIndent, t.Label)
	printColor(statusColors[status], testLine)
	if t.Exception != nil && t.Exception.Line > 0 && t.Exception.File != "" {
		for i := 0; i < longestStatus+len(" ... "); i++ {
			indent += " "
		}

		printColor(statusColors[status], fmt.Sprintf("\n%sError: %s", indent+"  ", t.Exception.Message))
		printColor(statusColors["note"], fmt.Sprintf("\n%s%s:%d", indent, t.Exception.File, t.Exception.Line))
	}
	if *format == "onlyfail" {
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
	if *quiet {
		return
	}
	clearLine()
	fmt.Printf("\nResults: %d pass, %d fail, %d error, %d skip\n",
		t.Counts.Pass, t.Counts.Fail, t.Counts.Error, t.Counts.Omit)
}

func clearLine() {
	if *format == "onlyfail" {
		terminal.Stdout.ClearLine()
	}
}

var statusColors = map[string]string{
	"pass":  "g",
	"skip":  "y",
	"fail":  "r",
	"error": "m",
	"todo":  "b",
	"note":  "w",
}

var longestStatus int

func init() {
	for status := range statusColors {
		if len(status) > longestStatus {
			longestStatus = len(status)
		}
	}
}

func printColor(color, text string) {
	if nocolor {
		terminal.Stdout.Print(text)
		return
	}
	terminal.Stdout.Color(color).Print(text).Reset()
}
