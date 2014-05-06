package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func postSuite(buf io.Reader, host string, project, testType string) error {
	if len(project) == 0 {
		project = "default"
	}

	addr := &url.URL{Scheme: "http", Host: host, Path: fmt.Sprintf("/projects/%s/tests", project)}
	_, err := http.Post(addr.String(), "application/tap", buf)
	if err != nil {
		return err
	}
	return nil
}

func validatePost(host string) error {
	addr := &url.URL{Scheme: "http", Host: host, Path: "/projects"}
	if _, err := http.Get(addr.String()); err != nil {
		return fmt.Errorf("Could not connect to server: %s", err)
	}
	return nil
}
