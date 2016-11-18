package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
		arr, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		var id string
		if err = json.Unmarshal(arr, &id); err != nil {
			err = fmt.Errorf("Error decoding response: %s: %s", string(arr), err)
		}
		return id, err
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
