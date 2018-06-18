package main

import (
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
