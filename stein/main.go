package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/beatgammit/stein"
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/render"
	log "github.com/jcelliott/lumber"
)

var (
	steinAddr string
	dbType    string
	dbName    string
	couchAddr string
)

var (
	db DB
)

func init() {
	flag.StringVar(&steinAddr, "addr", "localhost:3000", "address and port where stein will run")
	flag.StringVar(&dbType, "dbtype", "fs", "database type to use: fs, couchdb")
	flag.StringVar(&dbName, "dbname", "test_results", "database to use")
	flag.StringVar(&couchAddr, "couchaddr", "localhost:5984", "address of couchdb")

	flag.Parse()
}

func getProjects(re render.Render) {
	projs, err := db.GetProjects()
	if err != nil {
		re.JSON(500, err.Error())
	} else {
		re.JSON(200, projs)
	}
}

func getTestsByProject(params martini.Params, re render.Render) {
	tests, err := db.GetTests(params["project"])
	if err != nil {
		re.JSON(500, err.Error())
	} else {
		re.JSON(200, tests)
	}
}

func uploadTest(params martini.Params, r *http.Request, re render.Render) {
	id := time.Now().Format(time.RFC3339)
	s, err := stein.Parse(r.Body)
	if err != nil {
		re.JSON(500, err.Error())
		return
	}

	if typ, ok := params["type"]; ok {
		s.Type = typ
	}

	if err = db.Save(params["project"], id, s); err != nil {
		re.JSON(500, err.Error())
	} else {
		re.JSON(200, id)
	}
}

func getTest(params martini.Params, re render.Render) {
	s, err := db.GetTest(params["project"], params["test"])
	if err != nil {
		re.JSON(500, err.Error())
	} else {
		re.JSON(200, s)
	}
}

func getTestTypes(params martini.Params, re render.Render) {
	s, err := db.GetTestTypes(params["project"])
	if err != nil {
		re.JSON(500, err.Error())
	} else {
		re.JSON(200, s)
	}
}

func getTestsByType(params martini.Params, re render.Render) {
	s, err := db.GetTestsByType(params["project"], params["type"])
	if err != nil {
		re.JSON(500, err.Error())
	} else {
		re.JSON(200, s)
	}
}

func main() {
	var err error
	switch dbType {
	case "fs":
		db, err = NewFileStore(dbName)
	case "couchdb":
		db, err = NewCouchDB(couchAddr, dbName, "", "")
	default:
		err = fmt.Errorf("Unsupported database type: %s", dbType)
	}

	if err != nil {
		log.Error("Error initializing database: %s", err)
		os.Exit(1)
		return
	}

	m := martini.Classic()

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	m.Use(martini.Static(path.Join(cwd, "web")))
	m.Use(render.Renderer())

	m.Get("/projects", getProjects)
	m.Get("/projects/:project", getTestsByProject)
	m.Get("/projects/:project/tests", getTestsByProject)
	m.Post("/projects/:project", uploadTest)
	m.Get("/projects/:project/tests/:test", getTest)
	m.Get("/projects/:project/types", getTestTypes)
	m.Get("/projects/:project/types/:type", getTestsByType)
	m.Post("/projects/:project/types/:type", uploadTest)

	fmt.Printf("Server listening: %s\n", steinAddr)
	log.Fatal("%s", http.ListenAndServe(steinAddr, m))
}
