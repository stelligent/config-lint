package main

import (
	"fmt"
	"testing"

	"github.com/stelligent/config-lint/assertion"
	"github.com/stelligent/config-lint/linter"
	"github.com/stretchr/testify/assert"
)

// Run built in rules against Terraform v0.11 parser
// TODO, both this tf built in tests need to be updated for new test locations. WIP, 3/10/2020
func TestTerraform11BuiltInRules(t *testing.T) {

	// Define file to load rules from
	// This file is located under cli/assets/
	ruleSet := loadRules(t, "terraform.yml")

	// Get list of test cases
	testCases := []BuiltInTestCase{
		// AWS
	}

	// Run test cases
	// test files must be included under testdata/builtin/terraform11
	for _, tc := range testCases {
		filenames := []string{"testdata/builtin/terraform11/" + tc.Filename}
		options := linter.Options{
			RuleIDs: []string{tc.RuleID},
		}
		vs := assertion.StandardValueSource{}

		// Defining 'tf11' for the Parser type
		l, err := linter.NewLinter(ruleSet, vs, filenames, "tf11")
		report, err := l.Validate(ruleSet, options)
		assert.Nil(t, err, "Validate failed for file")

		warningViolationsReported := getViolationsString("WARNING", report.Violations)
		warningMessage := fmt.Sprintf("Expecting %d warnings for RuleID %s in File %s:\n %s", tc.WarningCount, tc.RuleID, tc.Filename, warningViolationsReported)
		assert.Equal(t, tc.WarningCount, numberOfWarnings(report.Violations), warningMessage)
		failureViolationsReported := getViolationsString("FAILURE", report.Violations)
		failureMessage := fmt.Sprintf("Expecting %d failures for RuleID %s in File %s:\n %s", tc.FailureCount, tc.RuleID, tc.Filename, failureViolationsReported)
		assert.Equal(t, tc.FailureCount, numberOfFailures(report.Violations), failureMessage)
	}
}
