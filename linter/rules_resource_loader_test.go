package linter

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRulesLinterValidate(t *testing.T) {
	options := Options{
		Tags:    []string{},
		RuleIDs: []string{},
	}
	ruleSet := loadRulesForTest("./testdata/rules/rules.yml", t)
	filenames := []string{"./testdata/rules/rules.yml"}
	loader := RulesResourceLoader{}
	linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: loader}
	report, err := linter.Validate(ruleSet, options)
	assert.Nil(t, err, "Expecting Validate to run without error")
	assert.Equal(t, 2, len(report.ResourcesScanned), "Expecting Validate to scan 2 resources")
	assert.Equal(t, 1, len(report.FilesScanned), "Expecting Validate to scan 1 file")
	assert.Equal(t, 0, len(report.Violations), "Expecting Validate to find 0 violations")
}

func TestRulesLinterSearch(t *testing.T) {
	ruleSet := loadRulesForTest("./testdata/rules/rules.yml", t)
	filenames := []string{"./testdata/rules/rules.yml"}
	loader := RulesResourceLoader{}
	linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: loader}
	var b bytes.Buffer
	linter.Search(ruleSet, "message", &b)
	assert.Contains(t, b.String(), "needs", "Expecting TestRulesLinterSearch to find string in output")
}
