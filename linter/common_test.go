package linter

import (
	"github.com/stelligent/config-lint/assertion"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestLoadYamlFileError(t *testing.T) {
	_, err := loadYAML("does-not-exist.yml")
	if err == nil {
		t.Errorf("LoadYaml should fail for missing file")
	}
}

func TestLoadYamlParseError(t *testing.T) {
	_, err := loadYAML("./testdata/resources/invalid.yml")
	if err == nil {
		t.Errorf("LoadYaml should fail for file with invalid YAML")
	}
	if !strings.Contains(err.Error(), "error converting YAML to JSON") {
		t.Errorf("Expecting parse error for invalid YAML")
	}
}

func TestLoadYamlUnexpectedFormat(t *testing.T) {
	_, err := loadYAML("./testdata/rules/bad-format.yml")
	assert.NotNil(t, err, "YAML with unexpected format should return error")
	assert.Contains(t, err.Error(), "YAML in unexpected format")
}

func TestGetResourceIDFromFilename(t *testing.T) {
	expected := "resource.yml"
	n := getResourceIDFromFilename("path/to/resource.yml")
	if n != expected {
		t.Errorf("expecting getResourceIDFromFilename returned %s, expected '%s'", n, expected)
	}
}

func TestCombineValidationReports(t *testing.T) {
	r1 := assertion.ValidationReport{FilesScanned: []string{"one"}}
	r2 := assertion.ValidationReport{FilesScanned: []string{"two"}}
	r := CombineValidationReports(r1, r2)
	if len(r.FilesScanned) != 2 {
		t.Errorf("expecting CombineValidationReports to have 2 entries for FilesScanned")
	}
}
