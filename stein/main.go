package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/beatgammit/stein"
	"net/http"
	"time"
)

func main() {
	projects := make(map[string]map[string]*stein.Suite)

	m := martini.Classic()
	m.Use(martini.Static("build/web"))
	m.Get("/projects", func () string {
		var projs []string
		for key := range projects {
			projs = append(projs, key)
		}
		b, _ := json.Marshal(projs)
		return string(b)
	})
	m.Put("/projects/:project", func (params martini.Params) string {
		projects[params["project"]] = make(map[string]*stein.Suite)
		return "created"
	})

	m.Get("/projects/:project/tests", func (params martini.Params) string {
		proj := params["project"]
		var tests []string
		for test := range projects[proj] {
			tests = append(tests, test)
		}
		b, _ := json.Marshal(tests)
		return string(b)
	});

	m.Post("/projects/:project/tests", func (params martini.Params, r *http.Request) string {
		proj := params["project"]
		if _, ok := projects[proj]; !ok {
			return "project not created"
		}

		id := time.Now().Format(time.RFC3339)
		projects[proj][id] = nil

		if tests, ok := projects[proj]; !ok {
			return "project doesn't exist"
		} else if _, ok := tests[id]; !ok {
			return "test doesn't exist"
		} else {
			s, err := stein.Parse(r.Body)
			if err != nil {
				return err.Error()
			}

			tests[id] = s
			return id
		}
	})
	m.Get("/projects/:project/tests/:test", func (params martini.Params) string {
		proj := params["project"]
		id := params["test"]
		if tests, ok := projects[proj]; !ok {
			return "null"
		} else if test, ok := tests[id]; !ok {
			return "null"
		} else {
			b, _ := json.Marshal(test)
			return string(b)
		}
	})
	m.Run()
}
