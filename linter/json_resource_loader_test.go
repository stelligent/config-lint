package linter

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestJSONLinterValidate(t *testing.T) {
	options := Options{
		Tags:    []string{},
		RuleIDs: []string{},
	}
	ruleSet := loadRulesForTest("./testdata/rules/generic-json.yml", t)
	filenames := []string{"./testdata/resources/users.json"}
	loader := JSONResourceLoader{Resources: ruleSet.Resources}
	linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: loader}
	report, err := linter.Validate(ruleSet, options)
	assert.Nil(t, err, "Expecting Validate to run without error")
	assert.Equal(t, 3, len(report.ResourcesScanned), "Expecting Validate to scan 3 resources")
	assert.Equal(t, 1, len(report.FilesScanned), "Expecting Validate to scan 1 file")
	assert.Equal(t, 1, len(report.Violations), "Expecting Validate to find 1 violation")
}

func TestJSONLinterSearch(t *testing.T) {
	ruleSet := loadRulesForTest("./testdata/rules/generic-json.yml", t)
	filenames := []string{"./testdata/resources/users.json"}
	loader := JSONResourceLoader{Resources: ruleSet.Resources}
	linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: loader}
	var b bytes.Buffer
	linter.Search(ruleSet, "Department", &b)
	if !strings.Contains(b.String(), "Audit") {
		t.Error("Expecting TestJSONLinterSearch to find string in output")
	}
}
