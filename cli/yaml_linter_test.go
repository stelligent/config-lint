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
	report, err := linter.Validate(filenames, ruleSet, emptyTags, emptyIds)
	if err != nil {
		t.Error("Expecting TestYAMLLinter to not return an error")
	}
	if len(report.ResourcesScanned) != 17 {
		t.Errorf("TestTerraformLinter scanned %d resources, expecting 17", len(report.ResourcesScanned))
	}
	if len(report.FilesScanned) != 1 {
		t.Errorf("TestYAMLLinter scanned %d files, expecting 1", len(report.FilesScanned))
	}
	if len(report.Violations) != 3 {
		t.Errorf("TestYAMLLinter returned %d violations, expecting 3", len(report.Violations))
	}
}
