package main

import (
	"testing"
)

func TestYAMLLinter(t *testing.T) {
	emptyTags := []string{}
	emptyIds := []string{}
	linter := YAMLLinter{Log: testLogger}
	ruleSet := loadRulesForTest("./testdata/rules/generic.yml", t)
	filenames := []string{"./testdata/resources/generic.config"}
	filesScanned, violations, err := linter.Validate(filenames, ruleSet, emptyTags, emptyIds)
	if err != nil {
		t.Error("Expecting TestYAMLLinter to not return an error")
	}
	if len(filesScanned) != 1 {
		t.Errorf("TestYAMLLinter scanned %d files, expecting 1", len(filesScanned))
	}
	if len(violations) != 3 {
		t.Errorf("TestYAMLLinter returned %d violations, expecting 3", len(violations))
	}
}
