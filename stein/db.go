package main

import "github.com/beatgammit/stein"

type DB interface {
	// GetProjects gets all projects with at least one test
	GetProjects() ([]string, error)
	// GetTests gets all tests for a product.
	GetTests(project string) ([]string, error)
	// GetTestTypes gets the test types for a project
	GetTestTypes(project string) ([]string, error)
	// GetTestsByType gets tests for a project given a test type
	GetTestsByType(project, typ string) ([]string, error)
	// GetTest gets a specific test.
	GetTest(project, test string) (*stein.Suite, error)
	// Save stores a test in the database.
	Save(project, id string, s *stein.Suite) error
}
