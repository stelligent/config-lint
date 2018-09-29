package linter

import (
	"encoding/csv"
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/stelligent/config-lint/assertion"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func readContent(filename string) ([]byte, error) {
	if filename == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(filename)
}

func loadYAML(filename string) ([]interface{}, error) {
	empty := []interface{}{}
	content, err := readContent(filename)
	if err != nil {
		return empty, err
	}

	var yamlData interface{}
	err = yaml.Unmarshal(content, &yamlData)
	if err != nil {
		return empty, err
	}
	m := yamlData.(map[string]interface{})
	return []interface{}{m}, nil
}

func loadJSON(filename string) ([]interface{}, error) {
	empty := []interface{}{}
	content, err := readContent(filename)
	if err != nil {
		return empty, err
	}

	var jsonData interface{}
	err = json.Unmarshal(content, &jsonData)
	if err != nil {
		return empty, err
	}
	m := jsonData.(map[string]interface{})
	return []interface{}{m}, nil
}

func loadCSV(filename string) ([][]string, error) {
	content, err := readContent(filename)
	if err != nil {
		return [][]string{}, err
	}
	return csv.NewReader(strings.NewReader(string(content))).ReadAll()
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
