package main

import "github.com/beatgammit/stein"

type DB interface {
	Create(project string) error
	GetProjects() ([]string, error)
	GetTests(project string) ([]string, error)
	GetTest(project, test string) (*stein.Suite, error)
	ProjectExists(project string) (bool, error)
	Save(project, test string, s *stein.Suite) error
}
