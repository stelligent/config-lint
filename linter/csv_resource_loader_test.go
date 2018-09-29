package linter

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCSVLinterValidate(t *testing.T) {
	options := Options{
		Tags:    []string{},
		RuleIDs: []string{},
	}
	ruleSet := loadRulesForTest("./testdata/rules/generic-csv.yml", t)
	filenames := []string{"./testdata/resources/users.csv"}
	loader := CSVResourceLoader{Columns: ruleSet.Columns}
	linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: loader}
	report, err := linter.Validate(ruleSet, options)
	assert.Nil(t, err, "Expecting Validate to run without error")
	assert.Equal(t, 3, len(report.ResourcesScanned), "Expecting Validate to scan 3 resources")
	assert.Equal(t, 1, len(report.Violations), "Expecting Validate to find 1 violation")
}

func TestCSVLinterSearch(t *testing.T) {
	ruleSet := loadRulesForTest("./testdata/rules/generic-csv.yml", t)
	filenames := []string{"./testdata/resources/users.csv"}
	loader := CSVResourceLoader{Columns: ruleSet.Columns}
	linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: loader}
	var b bytes.Buffer
	linter.Search(ruleSet, "Department", &b)
	assert.Contains(t, b.String(), "Audit", "Expecting TestCSVLinterSearch to find string in output")
}
