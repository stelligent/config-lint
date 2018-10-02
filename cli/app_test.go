package main

import (
	"bytes"
	"github.com/stelligent/config-lint/assertion"
	"github.com/stelligent/config-lint/linter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadTerraformRules(t *testing.T) {
	_, err := loadBuiltInRuleSet("assets/terraform.yml")
	if err != nil {
		t.Errorf("Cannot load built-in Terraform rules")
	}
}

func TestLoadValidateRules(t *testing.T) {
	_, err := loadBuiltInRuleSet("assets/lint-rules.yml")
	if err != nil {
		t.Errorf("Cannot load built-in rules for -validate option")
	}
}

func TestExcludeAll(t *testing.T) {
	filenames := []string{"file1.tf", "file2.tf", "file3.tf"}
	patterns := []string{"*.tf"}
	filtered := excludeFilenames(filenames, patterns)
	if len(filtered) != 0 {
		t.Errorf("Expecting all files to be excluded, but files are %v", filtered)
	}
}

func TestExcludeOnePattern(t *testing.T) {
	filenames := []string{"file1.tf", "file2.tf", "file3.tf"}
	patterns := []string{"*1.tf"}
	filtered := excludeFilenames(filenames, patterns)
	if len(filtered) != 2 {
		t.Errorf("Expecting one file to be excluded, but files are %v", filtered)
	}
}

func TestExcludeMultiplePattern(t *testing.T) {
	filenames := []string{"file1.tf", "file2.tf", "file3.tf"}
	patterns := []string{"*1.tf", "*2.tf"}
	filtered := excludeFilenames(filenames, patterns)
	if len(filtered) != 1 {
		t.Errorf("Expecting two files to be excluded, but files are %v", filtered)
	}
}

func TestExcludeFrom(t *testing.T) {
	excludeFromFilenames := []string{"./testdata/exclude-list"}
	patterns, err := loadExcludePatterns([]string{}, excludeFromFilenames)
	if err != nil {
		t.Errorf("Expecting loadExcludePatterns returned error: %s", err.Error())
	}
	if len(patterns) != 2 {
		t.Errorf("Expecting to load 2 patterns from excludeFromFilenames, not %v", patterns)
	}
	if patterns[0] != "*1.tf" {
		t.Errorf("Expecting first pattern from file to be '*1.tf', not '%s'", patterns[0])
	}
	if patterns[1] != "*2.tf" {
		t.Errorf("Expecting second pattern from file to be '*2.tf', not '%s'", patterns[1])
	}
}

func TestProfileExceptions(t *testing.T) {
	filenames := []string{"./testdata/terraform.yml"}
	ruleSets, err := loadRuleSets(filenames)
	if err != nil {
		t.Errorf("Expecting loadRuleSets to not return error: %s", err.Error())
	}
	profileExceptions := []RuleException{
		{
			RuleID:           "RULE_1",
			ResourceCategory: "resource",
			ResourceType:     "aws_instance",
			Comments:         "Testing",
			ResourceID:       "my-special-resource",
		},
	}
	ruleSets = addExceptions(ruleSets, profileExceptions)
	ruleExceptions := ruleSets[0].Rules[0].Except
	if len(ruleExceptions) != 1 {
		t.Errorf("Expecting Rule.Except to have one ID: %v", ruleExceptions)
		return
	}
	id := ruleExceptions[0]
	if id != "my-special-resource" {
		t.Errorf("Unexpected ResourceID found in Except: %s", id)
	}
}

func TestBuiltInTerraformRules(t *testing.T) {
	ruleSet, err := loadBuiltInRuleSet("assets/lint-rules.yml")
	if err != nil {
		t.Errorf("Expecting loadBuiltInRulesSet to not return error: %s", err.Error())
	}
	vs := assertion.StandardValueSource{}
	filenames := []string{"assets/terraform.yml"}
	l, err := linter.NewLinter(ruleSet, vs, filenames)
	if err != nil {
		t.Errorf("Expecting NewLinter to not return error: %s", err.Error())
	}
	options := linter.Options{}
	report, err := l.Validate(ruleSet, options)
	if err != nil {
		t.Errorf("Expecting Validate to not return error: %s", err.Error())
	}
	if len(report.Violations) != 0 {
		t.Errorf("Expecting Validate for built in rules to not report any violations: %v", report.Violations)
	}
}

func TestBuiltInLinterRules(t *testing.T) {
	ruleSet, err := loadBuiltInRuleSet("assets/lint-rules.yml")
	if err != nil {
		t.Errorf("Expecting loadBuiltInRulesSet to not return error: %s", err.Error())
	}
	vs := assertion.StandardValueSource{}
	filenames := []string{"assets/lint-rules.yml"}
	l, err := linter.NewLinter(ruleSet, vs, filenames)
	if err != nil {
		t.Errorf("Expecting NewLinter to not return error: %s", err.Error())
	}
	options := linter.Options{}
	report, err := l.Validate(ruleSet, options)
	if err != nil {
		t.Errorf("Expecting Validate to not return error: %s", err.Error())
	}
	if len(report.Violations) != 0 {
		t.Errorf("Expecting Validate for built in lint-rules to not report any violations: %v", report.Violations)
	}
}

func TestPrintReport(t *testing.T) {
	r := assertion.ValidationReport{}
	var b bytes.Buffer
	err := printReport(&b, r, "")
	assert.Nil(t, err, "Expecting printReport to run without error")
	assert.Contains(t, b.String(), "FilesScanned\": null")
	assert.Contains(t, b.String(), "ResourcesScanned\": null")
	assert.Contains(t, b.String(), "Violations\": null")
}

func TestPrintReportWithQueryString(t *testing.T) {
	r := assertion.ValidationReport{
		Violations: []assertion.Violation{
			assertion.Violation{RuleMessage: "Houston, we have a problem"},
		},
	}
	var b bytes.Buffer
	err := printReport(&b, r, "Violations[]")
	assert.Nil(t, err, "Expecting printReport to run without error")
	assert.Contains(t, b.String(), "RuleMessage")
	assert.NotContains(t, b.String(), "Violations")
	assert.NotContains(t, b.String(), "FilesScanned")
	assert.NotContains(t, b.String(), "ResourcesScanned")
}

type TestReportWriter struct {
	Report assertion.ValidationReport
}

func (w TestReportWriter) WriteReport(r assertion.ValidationReport, o LinterOptions) {
	w.Report = r
}

func TestApplyRules(t *testing.T) {
	ruleSets := []assertion.RuleSet{
		assertion.RuleSet{
			Type: "JSON",
		},
	}
	args := arrayFlags{}
	options := LinterOptions{}
	w := TestReportWriter{}
	exitCode := applyRules(ruleSets, args, options, w)
	assert.Equal(t, exitCode, 0, "Expecting applyRules to return 0")
	assert.Empty(t, w.Report.Violations, "Expecting empty report")
}

func TestValidateRules(t *testing.T) {
	filenames := []string{"./testdata/has-properties.yml"}
	w := TestReportWriter{}
	validateRules(filenames, w)
	assert.Empty(t, w.Report.Violations, "Expecting empty report for validateRules")
}
