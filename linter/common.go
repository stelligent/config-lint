package linter

import (
	"encoding/csv"
	"encoding/json"
	"errors"
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
	var result []interface{}
	docs := strings.Split(string(content), "---")
	for _, doc := range docs {
		if len(doc) == 0 {
			continue // skip empty documents, which can happen when a file starts with ---
		}
		var yamlData interface{}
		err = yaml.Unmarshal([]byte(doc), &yamlData)
		if err != nil {
			return empty, err
		}
		// Expecting every document to be an object (not simply valid YAML, such as a list)
		m, ok := yamlData.(map[string]interface{})
		if ok {
			result = append(result, m)
		} else {
			return []interface{}{}, errors.New("YAML in unexpected format")
		}
	}
	return result, nil
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
