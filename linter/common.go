package linter

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/stelligent/config-lint/assertion"
	"io/ioutil"
	"os"
	"path/filepath"
)

func loadYAML(filename string) ([]interface{}, error) {
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

// CombineValidationReports merges results from two separate Validate runs
func CombineValidationReports(r1, r2 assertion.ValidationReport) assertion.ValidationReport {
	return assertion.ValidationReport{
		FilesScanned:     append(r1.FilesScanned, r2.FilesScanned...),
		ResourcesScanned: append(r1.ResourcesScanned, r2.ResourcesScanned...),
		Violations:       append(r1.Violations, r2.Violations...),
	}
}
