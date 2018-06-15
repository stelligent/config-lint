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
