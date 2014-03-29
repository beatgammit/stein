package main

import "github.com/beatgammit/stein"

type DB interface {
	// GetProjects gets all projects with at least one test
	GetProjects() ([]string, error)
	// GetTests gets all tests for a product.
	GetTests(project string) ([]string, error)
	// GetTest gets a specific test.
	GetTest(project, test string) (*stein.Suite, error)
	// Save stores a test in the database.
	Save(project, test string, s *stein.Suite) error
}
