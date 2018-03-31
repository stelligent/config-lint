package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/stelligent/config-lint/assertion"
	"io/ioutil"
	"os"
	"path/filepath"
)

func loadYAML(filename string, log assertion.LoggingFunction) ([]interface{}, error) {
	empty := []interface{}{}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, filename, err.Error())
		return empty, err
	}

	var yamlData interface{}
	err = yaml.Unmarshal(content, &yamlData)
	if err != nil {
		fmt.Fprintln(os.Stderr, filename, err.Error())
		return empty, err
	}
	m := yamlData.(map[string]interface{})
	return []interface{}{m}, nil
}

func getResourceIDFromFilename(filename string) string {
	_, resourceID := filepath.Split(filename)
	return resourceID
}
