package linter

import (
	"bytes"
	"strings"
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
	if err != nil {
		t.Error("Expecting TestYAMLLinter to not return an error")
	}
	if len(report.ResourcesScanned) != 17 {
		t.Errorf("TestYAMLLinter scanned %d resources, expecting 17", len(report.ResourcesScanned))
	}
	if len(report.FilesScanned) != 1 {
		t.Errorf("TestYAMLLinter scanned %d files, expecting 1", len(report.FilesScanned))
	}
	if len(report.Violations) != 3 {
		t.Errorf("TestYAMLLinter returned %d violations, expecting 3", len(report.Violations))
	}
}

func TestYAMLLinterSearch(t *testing.T) {
	ruleSet := loadRulesForTest("./testdata/rules/generic-yaml.yml", t)
	filenames := []string{"./testdata/resources/generic.config"}
	loader := YAMLResourceLoader{Resources: ruleSet.Resources}
	linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: loader}
	var b bytes.Buffer
	linter.Search(ruleSet, "name", &b)
	if !strings.Contains(b.String(), "gadget") {
		t.Error("Expecting TestYAMLLinterSearch to find string in output")
	}
}
