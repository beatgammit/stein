package main

import (
	"encoding/json"
	"github.com/beatgammit/stein"
	"os"
	"path"
)

type FS struct {
	basedir string
}

func NewFileStore(basedir string) (DB, error) {
	return FS{basedir}, os.MkdirAll(basedir, 0777)
}

func (fs FS) Create(project string) error {
	return os.Mkdir(path.Join(fs.basedir, project), 0777)
}

func (fs FS) GetProjects() ([]string, error) {
	f, err := os.Open(fs.basedir)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Readdirnames(0)
}

func (fs FS) ProjectExists(proj string) (bool, error) {
	f, err := os.Open(path.Join(fs.basedir, proj))
	if err != nil {
		return false, nil
	}
	f.Close()
	return true, nil
}

func (fs FS) GetTests(project string) ([]string, error) {
	f, err := os.Open(path.Join(fs.basedir, project, "tests"))
	if err != nil {
		return []string{}, nil
	}
	defer f.Close()
	return f.Readdirnames(0)
}

func (fs FS) parseTypes(project string) (map[string][]string, error) {
	types := make(map[string][]string)
	f, err := os.Open(path.Join(fs.basedir, project, "types.json"))
	if err != nil {
		return types, err
	}
	defer f.Close()

	return types, json.NewDecoder(f).Decode(&types)
}

func (fs FS) GetTestTypes(project string) ([]string, error) {
	ret := []string{}

	types, _ := fs.parseTypes(project)
	for k := range types {
		ret = append(ret, k)
	}
	return ret, nil
}

func (fs FS) GetTestsByType(project, typ string) ([]string, error) {
	types, err := fs.parseTypes(project)
	if err != nil {
		return []string{}, nil
	}
	return types[typ], nil
}

func (fs FS) GetTest(project, test string) (*stein.Suite, error) {
	f, err := os.Open(path.Join(fs.basedir, project, "tests", test))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var s stein.Suite
	return &s, json.NewDecoder(f).Decode(&s)
}

func (fs FS) Save(project, test string, s *stein.Suite) error {
	if err := os.MkdirAll(path.Join(fs.basedir, project, "tests"), 0777); err != nil {
		return err
	}
	f, err := os.Create(path.Join(fs.basedir, project, "tests", test))
	if err != nil {
		return err
	}
	defer f.Close()

	types, _ := fs.parseTypes(project)
	types[s.Type] = append(types[s.Type], test)
	typesFile, err := os.Create(path.Join(fs.basedir, project, "types.json"))
	if err == nil {
		defer typesFile.Close()
		json.NewEncoder(typesFile).Encode(types)
	}

	return json.NewEncoder(f).Encode(s)
}
