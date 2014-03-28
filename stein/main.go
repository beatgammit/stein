package main

import (
	"encoding/json"
	"github.com/beatgammit/stein"
	"github.com/codegangsta/martini"
	"net/http"
	"time"
)

func main() {
	db, err := NewFileStore("file_store")
	if err != nil {
		panic(err)
	}

	m := martini.Classic()
	m.Use(martini.Static("build/web"))
	m.Get("/projects", func() (string, int) {
		projs, err := db.GetProjects()
		if err != nil {
			return err.Error(), 500
		}
		b, _ := json.Marshal(projs)
		return string(b), 200
	})
	m.Get("/projects/:project", func(params martini.Params) (string, int) {
		if exists, err := db.ProjectExists(params["project"]); err != nil {
			return err.Error(), 500
		} else if exists {
			return "exists", 200
		} else {
			return "doesn't exist", 404
		}
	})
	m.Put("/projects/:project", func(params martini.Params) (string, int) {
		err := db.Create(params["project"])
		if err != nil {
			return err.Error(), 500
		}
		return "created", 200
	})

	m.Get("/projects/:project/tests", func(params martini.Params) (string, int) {
		tests, err := db.GetTests(params["project"])
		if err != nil {
			return err.Error(), 500
		}
		b, _ := json.Marshal(tests)
		return string(b), 200
	})

	m.Post("/projects/:project/tests", func(params martini.Params, r *http.Request) (string, int) {
		id := time.Now().Format(time.RFC3339)
		s, err := stein.Parse(r.Body)
		if err != nil {
			return err.Error(), 500
		}

		err = db.Save(params["project"], id, s)
		if err != nil {
			return err.Error(), 500
		}
		return id, 200
	})
	m.Get("/projects/:project/tests/:test", func(params martini.Params) (string, int) {
		s, err := db.GetTest(params["project"], params["test"])
		if err != nil {
			return err.Error(), 500
		}
		b, _ := json.Marshal(s)
		return string(b), 200
	})
	m.Run()
}
