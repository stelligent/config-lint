package linter

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYAMLLinterValidate(t *testing.T) {
	options := Options{
		Tags:    []string{},
		RuleIDs: []string{},
	}
	ruleSet := loadRulesForTest("./testdata/rules/generic-yaml.yml", t)
	filenames := []string{"./testdata/resources/generic.config"}
	loader := YAMLResourceLoader{Resources: ruleSet.Resources}
	linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: loader}
	report, err := linter.Validate(ruleSet, options)
	assert.Nil(t, err, "Expecting Validate to run without error")
	assert.Equal(t, 17, len(report.ResourcesScanned), "Expecting Validate to scan 17 resources")
	assert.Equal(t, 1, len(report.FilesScanned), "Expecting Validate to scan 1 file")
	assert.Equal(t, 3, len(report.Violations), "Expecting Validate to find 3 violations")
}

func TestYAMLLinterSearch(t *testing.T) {
	ruleSet := loadRulesForTest("./testdata/rules/generic-yaml.yml", t)
	filenames := []string{"./testdata/resources/generic.config"}
	loader := YAMLResourceLoader{Resources: ruleSet.Resources}
	linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: loader}
	var b bytes.Buffer
	linter.Search(ruleSet, "name", &b)
	assert.Contains(t, b.String(), "gadget", "Expecting TestYAMLLinterSearch to find string in output")
}
