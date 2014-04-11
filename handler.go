package stein

import (
	"fmt"
)

type Handler interface {
	HandleDoc(string, string)
	HandleBeforeStream()
	HandleAfterStream()
	HandleSuite(*Suite)
	HandleCase(*Case)
	HandleTest(*Test)
	HandleNote(*Note)
	HandleTally(*Tally)
	HandleFinal(*Tally)
}

type DefaultHandler struct{}

func (h *DefaultHandler) HandleDoc(string, string) {}
func (h *DefaultHandler) HandleBeforeStream()      {}
func (h *DefaultHandler) HandleAfterStream()       {}
func (h *DefaultHandler) HandleSuite(*Suite)       {}
func (h *DefaultHandler) HandleCase(*Case)         {}
func (h *DefaultHandler) HandleTest(*Test)         {}
func (h *DefaultHandler) HandleNote(*Note)         {}
func (h *DefaultHandler) HandleTally(*Tally)       {}
func (h *DefaultHandler) HandleFinal(*Tally)       {}

type EchoHandler struct {
	*DefaultHandler
}

func (h *EchoHandler) HandleDoc(doc string, docType string) {
	fmt.Println(doc)
}

type DebugHandler struct {
	*EchoHandler
}

func (h *DebugHandler) HandleBeforeStream() {
	fmt.Println("Before stream")
}
func (h *DebugHandler) HandleAfterStream() {
	fmt.Println("After stream")
}
func (h *DebugHandler) HandleSuite(s *Suite) {
	fmt.Printf("%+v\n", s)
}
func (h *DebugHandler) HandleCase(c *Case) {
	fmt.Printf("%+v\n", c)
}
func (h *DebugHandler) HandleTest(t *Test) {
	fmt.Printf("%+v\n", t)
}
func (h *DebugHandler) HandleNote(n *Note) {
	fmt.Printf("%+v\n", n)
}
func (h *DebugHandler) HandleTally(t *Tally) {
	fmt.Printf("%+v\n", t)
}
func (h *DebugHandler) HandleFinal(t *Tally) {
	fmt.Printf("%+v\n", t)
}
