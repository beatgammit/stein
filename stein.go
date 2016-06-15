package stein

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/v1/yaml"
)

var (
	docSep      = regexp.MustCompile(`(?m)^---$|^\.\.\.$`)
	tapTestLine = regexp.MustCompile(`(?m)^(ok|not ok)(?: (\d+))?( [^#]+)?(?:\s*#(.*))?$`)
	tapPlan     = regexp.MustCompile(`(?m)^\d..(\d+|N).*$`)
	tapVersion  = regexp.MustCompile(`(?m)^TAP version \d+$`)
	tapBailOut  = regexp.MustCompile(`(?m)^Bail out!\s*(.*)$`)
)

func Parse(r io.Reader) (*Suite, error) {
	parser := Parser{&DefaultHandler{}}
	return parser.Parse(r)
}

type Parser struct {
	Handler Handler
}

func (p *Parser) Parse(r io.Reader) (*Suite, error) {
	rd := bufio.NewReader(r)
	b, err := rd.Peek(1)
	if err != nil {
		panic(err)
	}

	if p.Handler == nil {
		p.Handler = &DefaultHandler{}
	}

	s := new(Suite)
	switch b[0] {
	case '{':
		return s, s.parse(newDecoder("json", rd), p.Handler)
	case '-':
		return s, s.parse(newDecoder("yaml", rd), p.Handler)
	default:
		return ParseTap(rd, p.Handler)
	}
}

type Time struct {
	time.Time
}

func (t *Time) parse(val string) (err error) {
	var ti time.Time
	if ti, err = time.Parse("2006-01-02 15:04:05", val); err != nil {
		if ti, err = time.Parse(time.RFC3339, val); err != nil {
			return
		}
	}
	(*t).Time = ti
	return
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	// parse time, stripping the quotes
	return t.parse(string(data[1 : len(data)-1]))
}
func (t *Time) GetYAML() (tag string, value interface{}) {
	return "", t.String()
}
func (t *Time) SetYAML(tag string, value interface{}) bool {
	if v, ok := value.(string); ok {
		return t.parse(strings.Trim(v, `'"`)) == nil
	}
	return false
}

type Suite struct {
	Type  string
	Start *Time
	Count int
	Seed  int
	Rev   int
	Extra interface{} `yaml:",omitempty" json:",omitempty"`

	Tests []*Test
	Cases []*Case
	Notes []Note
	Final Tally
}

func (s *Suite) parse(dec decoder, handler Handler) (err error) {
	handler.HandleBeforeStream()
	var lastCase *Case
	for dec.Scan() {
		doc := dec.Scanner.Text()
		handler.HandleDoc(doc, dec.format)

		if strings.TrimSpace(doc) == "" {
			continue
		}

		if s.Start == nil {
			if err = dec.Unmarshal([]byte(doc), &s); err != nil {
				fmt.Println(err)
				continue
			}
			handler.HandleSuite(s)
			continue
		}

		var typ string
		if typ, err = dec.Type([]byte(doc)); err != nil {
			return
		}

		switch typ {
		case "case":
			var c Case
			if err = dec.Unmarshal([]byte(doc), &c); err != nil {
				fmt.Println(err)
				continue
			}

			if c.Level == 0 {
				s.Cases = append(s.Cases, &c)
				lastCase = &c
			} else if lastCase == nil {
				// TODO: possibly just add this case to the suite
				return fmt.Errorf("No parent found for test case: %v", c)
			} else {
				// find the parent
				for c.Level <= lastCase.Level {
					lastCase = lastCase.parent
				}
				// TODO: strict check to ensure c.Level == lastCase.Level + 1
				c.parent = lastCase
				lastCase.Subcases = append(lastCase.Subcases, &c)
			}
			handler.HandleCase(&c)
		case "test":
			var t Test
			if err = dec.Unmarshal([]byte(doc), &t); err != nil {
				fmt.Println(err)
				continue
			}
			if lastCase != nil {
				lastCase.Tests = append(lastCase.Tests, &t)
			} else {
				s.Tests = append(s.Tests, &t)
			}
			handler.HandleTest(&t)
		case "note":
			var n Note
			if err = dec.Unmarshal([]byte(doc), &n); err != nil {
				fmt.Println(err)
				continue
			}
			s.Notes = append(s.Notes, n)
			handler.HandleNote(&n)
		case "tally":
			var t Tally
			if err = dec.Unmarshal([]byte(doc), &t); err != nil {
				fmt.Println(err)
				continue
			}
			handler.HandleTally(&t)
		case "final":
			var t Tally
			if err = dec.Unmarshal([]byte(doc), &t); err != nil {
				fmt.Println(err)
				continue
			}
			s.Final = t
			handler.HandleFinal(&t)
		default:
			err = fmt.Errorf("Invalid type: %s", typ)
			return
		}
	}
	handler.HandleAfterStream()
	return
}

func (s Suite) Docs() (ret []interface{}) {
	ret = append(ret, s)
	for _, t := range s.Tests {
		ret = append(ret, t)
	}

	for _, c := range s.Cases {
		ret = append(ret, c.Docs()...)
	}

	for _, n := range s.Notes {
		ret = append(ret, n)
	}

	ret = append(ret, s.Final)
	return
}

func (s Suite) ToTapY() ([]byte, error) {
	var ret []byte
	for _, doc := range s.Docs() {
		if b, err := yaml.Marshal(doc); err != nil {
			return nil, err
		} else {
			ret = append(ret, []byte("---\n")...)
			ret = append(ret, b...)
			ret = append(ret, '\n')
		}
	}
	return append(ret, []byte("...")...), nil
}

type Case struct {
	Type  string
	Label string
	Level int

	Subtype string      `yaml:",omitempty" json:",omitempty"`
	Extra   interface{} `yaml:",omitempty" json:",omitempty"`

	// Not part of any TAP protocol
	Tests []*Test `yaml:",omitempty" json:",omitempty"`
	// Not part of any TAP protocol
	Subcases []*Case `yaml:",omitempty" json:",omitempty"`

	parent *Case
}

func (c Case) Docs() (ret []interface{}) {
	ret = append(ret, c)
	for _, t := range c.Tests {
		ret = append(ret, t)
	}
	for _, c := range c.Subcases {
		ret = append(ret, c.Docs()...)
	}
	return
}

// TODO: allow strings via UnmarshalJSON/YAML
type Snippet []map[string]string

type TestCoverage struct {
	File string `yaml:",omitempty" json:",omitempty"`
	// TODO: limit this to a number
	// according to tap-y/j, can be:
	// - line: 11
	// - range: 11..13
	// - list of lines: [11,12,13]
	Line interface{} `yaml:",omitempty" json:",omitempty"`
	Code string      `yaml:",omitempty" json:",omitempty"`
}

type TestException struct {
	Message   string
	File      string      `yaml:",omitempty" json:",omitempty"`
	Line      int         `yaml:",omitempty" json:",omitempty"`
	Source    string      `yaml:",omitempty" json:",omitempty"`
	Snippet   Snippet     `yaml:",omitempty" json:",omitempty"`
	Backtrace interface{} `yaml:",omitempty" json:",omitempty"`
}

type Test struct {
	Type   string
	Status string
	Label  string
	Time   float64

	Subtype   string         `yaml:",omitempty" json:",omitempty"`
	Setup     string         `yaml:",omitempty" json:",omitempty"`
	Expected  interface{}    `yaml:",omitempty" json:",omitempty"`
	Returned  interface{}    `yaml:",omitempty" json:",omitempty"`
	File      string         `yaml:",omitempty" json:",omitempty"`
	Line      int            `yaml:",omitempty" json:",omitempty"`
	Source    string         `yaml:",omitempty" json:",omitempty"`
	Snippet   Snippet        `yaml:",omitempty" json:",omitempty"`
	Coverage  *TestCoverage  `yaml:",omitempty" json:",omitempty"`
	Exception *TestException `yaml:",omitempty" json:",omitempty"`
	Stdout    string         `yaml:",omitempty" json:",omitempty"`
	Stderr    string         `yaml:",omitempty" json:",omitempty"`
	Extra     interface{}    `yaml:",omitempty" json:",omitempty"`
}

type Note struct {
	Type  string
	Text  string
	Extra interface{} `yaml:",omitempty" json:",omitempty"`
}

type Tally struct {
	Type   string
	Time   float64 `yaml:",omitempty" json:",omitempty"`
	Counts struct {
		Total int
		Pass  int
		Fail  int
		Error int
		Omit  int
		Todo  int
	}
	Extra interface{} `yaml:",omitempty" json:",omitempty"`
}

type decoder struct {
	*bufio.Scanner
	format    string
	Unmarshal func([]byte, interface{}) error
}

func scanYAMLDoc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if loc := docSep.FindIndex(data); loc != nil {
		if loc[0] == 0 {
			// just a start document marker
			return loc[1], nil, nil
		}
		// We found a complete document
		return loc[1], data[0:loc[0]], nil
	}
	if atEOF {
		// if we're at EOF without having matched docSep, there's nothing else to get
		return len(data), nil, nil
	}
	// Request more data.
	return 0, nil, nil
}

func newDecoder(typ string, r io.Reader) decoder {
	dec := decoder{Scanner: bufio.NewScanner(r)}
	// TODO: handle plain TAP with a decoder as well
	switch typ {
	case "json":
		dec.format = "json"
		dec.Unmarshal = json.Unmarshal
	case "yaml":
		dec.format = "yaml"
		dec.Unmarshal = yaml.Unmarshal
		dec.Scanner.Split(scanYAMLDoc)
	default:
		panic("Unsupported decoder")
	}
	return dec
}

func (dec decoder) Type(b []byte) (string, error) {
	m := make(map[string]interface{})
	if err := dec.Unmarshal(b, &m); err != nil {
		return "", fmt.Errorf("Invalid document: %s", err)
	}
	if typ, ok := m["type"].(string); !ok {
		return "", fmt.Errorf("Missing 'type' key in document")
	} else {
		return typ, nil
	}
}

func ParseTap(r io.Reader, handler Handler) (*Suite, error) {
	if handler == nil {
		handler = &DefaultHandler{}
	}
	handler.HandleBeforeStream()
	rd := bufio.NewReader(r)
	var s Suite
	s.Start = &Time{time.Now()}
	first := true
	var totalTests int
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		handler.HandleDoc(line, "tap")

		if line[0] == '#' {
			// TODO: do something with diagnostic lines
			continue
		}

		if matches := tapBailOut.FindStringSubmatch(line); matches != nil {
			// TODO: do something with the bail reason
			break
		}

		if first {
			if matches := tapVersion.FindStringSubmatch(line); matches != nil {
				// TAP 13+
				return nil, fmt.Errorf("TAP13 not supported yet")
			}
			first = false
			if matches := tapPlan.FindStringSubmatch(line); matches != nil {
				totalTests, _ = strconv.Atoi(matches[1])
				s.Count = totalTests
				// assuming plan line at beginning of output
				handler.HandleSuite(&s)
				continue
			}
		}

		matches := tapTestLine.FindStringSubmatch(line)
		if matches == nil {
			// TODO: handle extra data
			continue
		}

		t := Test{Type: "test"}
		switch matches[1] {
		case "ok":
			t.Status = "pass"
		case "not ok":
			// TODO: error?
			t.Status = "fail"
		}
		// ignore number
		t.Label = strings.TrimSpace(matches[3])
		directive := strings.TrimSpace(matches[4])
		switch {
		case strings.HasPrefix(strings.ToUpper(directive), "TODO"):
			s.Final.Counts.Todo++
			t.Label = directive
			t.Status = "todo"
		case strings.HasPrefix(strings.ToUpper(directive), "SKIP"):
			s.Final.Counts.Omit++
			t.Label = directive
			t.Status = "omit"
		case t.Status == "pass":
			s.Final.Counts.Pass++
		default:
			s.Final.Counts.Fail++
		}
		handler.HandleTest(&t)
		s.Final.Counts.Total++
		s.Tests = append(s.Tests, &t)
	}

	// fixup missing tests
	if totalTests > 0 && totalTests > s.Final.Counts.Total {
		s.Final.Counts.Fail += (totalTests - s.Final.Counts.Total)
		s.Final.Counts.Total = totalTests
	}

	handler.HandleFinal(&s.Final)
	handler.HandleAfterStream()
	return &s, nil
}
