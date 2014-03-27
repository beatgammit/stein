package stein

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gopkg.in/v1/yaml"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

var docSep *regexp.Regexp
var tapTestLine *regexp.Regexp
var tapPlan *regexp.Regexp
var tapVersion *regexp.Regexp

func init() {
	docSep = regexp.MustCompile(`(?m)^---$|^\.\.\.$`)
	tapTestLine = regexp.MustCompile(`(?m)^(ok|not ok)(?: (\d+))?( [^#]+)?(\s*#.*)$`)
	tapPlan = regexp.MustCompile(`(?m)^\d..(\d+|N).*$`)
	tapVersion = regexp.MustCompile(`(?m)^TAP version \d+$`)
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
		return t.parse(v) == nil
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

func (s *Suite) parse(dec decoder, parts []string) (err error) {
	var lastCase *Case
	for _, doc := range parts {
		if strings.TrimSpace(doc) == "" {
			continue
		}

		if s.Start == nil {
			if err = dec.Unmarshal([]byte(doc), &s); err != nil {
				return
			}
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
				return
			}

			if c.Level == 0 {
				s.Cases = append(s.Cases, &c)
				lastCase = &c
			} else if lastCase == nil {
				fmt.Println("No parent found for test case:", c)
				return
			} else {
				// find the parent
				for c.Level <= lastCase.Level {
					lastCase = lastCase.parent
				}
				// TODO: strict check to ensure c.Level == lastCase.Level + 1
				c.parent = lastCase
				lastCase.Subcases = append(lastCase.Subcases, &c)
			}
		case "test":
			var t Test
			if err = dec.Unmarshal([]byte(doc), &t); err != nil {
				return
			}
			if lastCase != nil {
				lastCase.Tests = append(lastCase.Tests, &t)
			} else {
				s.Tests = append(s.Tests, &t)
			}
		case "note":
			var n Note
			if err = dec.Unmarshal([]byte(doc), &n); err != nil {
				return
			}
			s.Notes = append(s.Notes, n)
		case "tally":
			var t Tally
			if err = dec.Unmarshal([]byte(doc), &t); err != nil {
				return
			}
			fmt.Println("tally:", t)
		case "final":
			var t Tally
			if err = dec.Unmarshal([]byte(doc), &t); err != nil {
				return
			}
			s.Final = t
		default:
			err = fmt.Errorf("Invalid type: %s", typ)
			return
		}
	}
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
	Type    string
	Subtype string
	Label   string
	Level   int
	Extra   interface{} `yaml:",omitempty" json:",omitempty"`

	Tests    []*Test `yaml:"-" json:"-"`
	Subcases []*Case `yaml:"-" json:"-"`

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

type Test struct {
	Type     string
	Subtype  string
	Status   string
	Setup    string
	Label    string
	Expected interface{}
	Returned interface{}
	File     string
	Line     int
	Source   string
	Snippet  Snippet
	Coverage struct {
		File string
		Line interface{}
		Code string
	}
	Exception struct {
		Message   string
		File      string
		Line      int
		Source    string
		Snippet   Snippet
		Backtrace interface{}
	}
	Stdout string `yaml:",omitempty" json:",omitempty"`
	Stderr string `yaml:",omitempty" json:",omitempty"`
	Time   float64
	Extra  interface{} `yaml:",omitempty" json:",omitempty"`
}

type Note struct {
	Type  string
	Text  string
	Extra interface{} `yaml:",omitempty" json:",omitempty"`
}

type Tally struct {
	Type   string
	Time   float64
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

type decoder string

func (dec decoder) Type(b []byte) (string, error) {
	m := make(map[string]interface{})
	switch dec {
	case "json":
		if err := json.Unmarshal(b, &m); err != nil {
			return "", fmt.Errorf("Invalid JSON: %s", err)
		}
	case "yaml":
		if err := yaml.Unmarshal(b, m); err != nil {
			return "", fmt.Errorf("Invalid YAML: %s", err)
		}
	default:
		panic("Unsupported decoder")
	}

	if typ, ok := m["type"].(string); !ok {
		return "", fmt.Errorf("Missing 'type' key in document")
	} else {
		return typ, nil
	}
}
func (dec decoder) Unmarshal(b []byte, v interface{}) error {
	switch dec {
	case "json":
		return json.Unmarshal(b, v)
	case "yaml":
		return yaml.Unmarshal(b, v)
	}
	panic("Unsupported decoder")
}

func Parse(r io.Reader) (*Suite, error) {
	rd := bufio.NewReader(r)
	b, err := rd.Peek(1)
	if err != nil {
		panic(err)
	}

	switch b[0] {
	case '{':
		return ParseJSON(rd)
	case '-':
		return ParseYaml(rd)
	default:
		return ParseTap(rd)
	}

}

func ParseJSON(r io.Reader) (*Suite, error) {
	rd := bufio.NewReader(r)
	var parts []string
	for {
		var line string
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		parts = append(parts, line)
	}

	s := new(Suite)
	return s, s.parse(decoder("json"), parts)
}

func ParseYaml(r io.Reader) (*Suite, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	parts := docSep.Split(string(b), -1)

	s := new(Suite)
	return s, s.parse(decoder("yaml"), parts)
}

func ParseTap(r io.Reader) (*Suite, error) {
	return nil, nil
}
