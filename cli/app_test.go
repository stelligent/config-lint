package main

import (
	"bytes"
	"github.com/gobuffalo/packr"
	"github.com/stelligent/config-lint/assertion"
	"github.com/stelligent/config-lint/linter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadTerraformRules(t *testing.T) {
	_, err := loadBuiltInRuleSet("terraform")
	if err != nil {
		t.Errorf("Cannot load built-in Terraform rules")
	}
}

func TestLoadValidateRules(t *testing.T) {
	_, err := loadBuiltInRuleSet("lint-rules.yml")
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

func TestExcludeSubdirectory(t *testing.T) {
	filenames := []string{"file1.tf", "foo/bar/secrets/database.yml"}
	patterns := []string{"foo/bar/secrets/*"}
	filtered := excludeFilenames(filenames, patterns)
	if len(filtered) != 1 {
		t.Errorf("Expecting secrets subdirectory to be excluded, but files are %v", filtered)
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

func TestBuiltRules(t *testing.T) {
	ruleSet, err := loadBuiltInRuleSet("lint-rules.yml")
	if err != nil {
		t.Errorf("Expecting loadBuiltInRuleSet to not return error: %s", err.Error())
	}
	vs := assertion.StandardValueSource{}

	// Get all rule files from the assets box
	box := packr.NewBox("./assets")
	allFilenames := box.List()
	var filenames []string
	for _, filename := range allFilenames {
		if isYamlFile(filename) && !isTestCase(filename) {
			filenames = append(filenames, "assets/"+filename)
		}
	}
	l, err := linter.NewLinter(ruleSet, vs, filenames, "")
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

type MockReportWriter struct {
	Report assertion.ValidationReport
}

func (w MockReportWriter) WriteReport(r assertion.ValidationReport, o LinterOptions) {
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
	w := MockReportWriter{}
	exitCode := applyRules(ruleSets, args, options, w)
	assert.Equal(t, exitCode, 0, "Expecting applyRules to return 0")
	assert.Empty(t, w.Report.Violations, "Expecting empty report")
}

func TestValidateRules(t *testing.T) {
	filenames := []string{"./testdata/has-properties.yml"}
	w := MockReportWriter{}
	validateRules(filenames, w)
	assert.Empty(t, w.Report.Violations, "Expecting empty report for validateRules")
}

func TestResourceMatch(t *testing.T) {
	testRule := []assertion.Rule{
		{
			ID:        "RULE_1",
			Category:  "resource",
			Resources: []string{"aws_instance", "aws_s3_bucket"},
		},
		{
			ID:       "RULE_2",
			Category: "resource",
			Resource: "aws_s3_bucket",
		},
	}
	profileExceptions := []RuleException{
		{
			RuleID:           "RULE_1",
			ResourceCategory: "resource",
			ResourceType:     "aws_instance",
			Comments:         "Testing",
			ResourceID:       "my-special-resource",
		},
		{
			RuleID:           "RULE_2",
			ResourceCategory: "resources",
			ResourceType:     "aws_s3_bucket",
			Comments:         "Testing",
			ResourceID:       "my-special-bucket",
		},
		{
			RuleID:           "RULE_2",
			ResourceCategory: "resources",
			ResourceType:     "aws_vpc",
			Comments:         "Should not match",
			ResourceID:       "my-vpc",
		},
	}

	assert.True(t, resourceMatch(testRule[0], profileExceptions[0]), "Expecting exception resource to be found in rule resources")
	assert.True(t, resourceMatch(testRule[1], profileExceptions[1]), "Expecting one to one match with exception resource and rule resource")
	assert.False(t, resourceMatch(testRule[1], profileExceptions[2]), "Expecting rule and exception to not match")
}

func TestLoadRuleSetsBadFilename(t *testing.T) {
	args := []string{"no-such-file.yml"}
	_, err := loadRuleSets(args)
	assert.NotNil(t, err, "LoadRuleSet with bad filename should return an error")
}

func TestLoadRuleSetsParseErrors(t *testing.T) {
	args := []string{"./testdata/syntax-errors.yml"}
	_, err := loadRuleSets(args)
	assert.NotNil(t, err, "Expecting rules file with syntax errors to fail")
	if err != nil {
		assert.Contains(t, err.Error(), "error unmarshaling JSON")
	}
}

func TestStdinFilename(t *testing.T) {
	filenames := getFilenames([]string{"-"})
	assert.Len(t, filenames, 1, "getFilenames should file 1 file")
	assert.Equal(t, filenames[0], "-", "getFilenames should allow - for stdin")
}

func TestGetFilenamesUsingDirectory(t *testing.T) {
	filenames := getFilenames([]string{"./testdata/dirtest"})
	assert.Len(t, filenames, 2)
	assert.Equal(t, "testdata/dirtest/a.yml", filenames[0])
	assert.Equal(t, "testdata/dirtest/b.yml", filenames[1])
}

func TestLoadFilenamesFromCommandLine(t *testing.T) {
	commandLineFilenames := []string{"command.yml"}
	profileFilenames := []string{"default.yml"}
	result := loadFilenames(commandLineFilenames, profileFilenames)
	assert.Equal(t, result, commandLineFilenames)
}

func TestLoadFilenamesFromProfile(t *testing.T) {
	commandLineFilenames := []string{}
	profileFilenames := []string{"default.yml"}
	result := loadFilenames(commandLineFilenames, profileFilenames)
	assert.Equal(t, result, profileFilenames)
}

func TestArrayFlags(t *testing.T) {
	var f arrayFlags
	assert.Equal(t, "", f.String(), "Default arrayFlags should return empty string")
	f.Set("first")
	f.Set("second")
	assert.Equal(t, arrayFlags{"first", "second"}, f, "Expecting arrayFlags to have two elements")
}

func TestLoadBuiltInRuleSetMissing(t *testing.T) {
	_, err := loadBuiltInRuleSet("missing.yml")
	assert.Contains(t, err.Error(), "File or directory doesnt exist", "loadBuiltInRuleSet should fail for missing file")
}
