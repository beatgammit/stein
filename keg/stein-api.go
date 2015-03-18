package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func postSuite(buf io.Reader, host, project, testType string) (string, error) {
	addr := &url.URL{Scheme: "http", Host: host, Path: fmt.Sprintf("/projects/%s", project)}
	if testType != "" {
		addr.Path += "/types/" + testType
	}
	resp, err := http.Post(addr.String(), "application/tap", buf)
	if err == nil {
		var id string
		return id, json.NewDecoder(resp.Body).Decode(&id)
	}
	return "", err
}

func validatePost(host string) error {
	addr := &url.URL{Scheme: "http", Host: host, Path: "/projects"}
	if _, err := http.Get(addr.String()); err != nil {
		return fmt.Errorf("Could not connect to server: %s", err)
	}
	return nil
}
