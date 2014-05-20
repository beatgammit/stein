package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func postSuite(buf io.Reader, host, project, testType string) error {
	addr := &url.URL{Scheme: "http", Host: host, Path: fmt.Sprintf("/projects/%s", project)}
	if testType != "" {
		addr.Path += "/types/" + testType
	}
	_, err := http.Post(addr.String(), "application/tap", buf)
	return err
}

func validatePost(host string) error {
	addr := &url.URL{Scheme: "http", Host: host, Path: "/projects"}
	if _, err := http.Get(addr.String()); err != nil {
		return fmt.Errorf("Could not connect to server: %s", err)
	}
	return nil
}
