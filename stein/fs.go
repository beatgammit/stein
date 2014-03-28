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
	f, err := os.Open(path.Join(fs.basedir, project))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Readdirnames(0)
}

func (fs FS) GetTest(project, test string) (*stein.Suite, error) {
	f, err := os.Open(path.Join(fs.basedir, project, test))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var s stein.Suite
	return &s, json.NewDecoder(f).Decode(&s)
}

func (fs FS) Save(project, test string, s *stein.Suite) error {
	f, err := os.Create(path.Join(fs.basedir, project, test))
	if err != nil {
		return err
	}
	return json.NewEncoder(f).Encode(s)
}
