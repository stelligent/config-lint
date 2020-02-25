package main

import (
	"strconv"
	"testing"

	"github.com/stelligent/config-lint/assertion"
	"github.com/stretchr/testify/assert"
)

func loadRules(t *testing.T, filename string) assertion.RuleSet {
	r, err := loadBuiltInRuleSet(filename)
	assert.Nil(t, err, "Cannot load ruleset file")
	return r
}

type BuiltInTestCase struct {
	Filename     string
	RuleID       string
	WarningCount int
	FailureCount int
}

func numberOfWarnings(violations []assertion.Violation) int {
	n := 0
	for _, v := range violations {
		if v.Status == "WARNING" {
			n++
		}
	}
	return n
}
func numberOfFailures(violations []assertion.Violation) int {
	n := 0
	for _, v := range violations {
		if v.Status == "FAILURE" {
			n++
		}
	}
	return n
}

// String build message for violations. Debug helper
func getViolationsString(severity string, violations []assertion.Violation) string {
	var violationsReported string
	for count, v := range violations {
		if v.Status == severity {
			violationsReported += strconv.Itoa(count+1) + ". Violation:"
			violationsReported += "\n\tRule Message: " + v.RuleMessage
			violationsReported += "\n\tRule Id: " + v.RuleID
			violationsReported += "\n\tResource ID: " + v.ResourceID
			violationsReported += "\n\tResource Type: " + v.ResourceType
			violationsReported += "\n\tCategory: " + v.Category
			violationsReported += "\n\tStatus: " + v.Status
			violationsReported += "\n\tAssertion Message: " + v.AssertionMessage
			violationsReported += "\n\tFilename: " + v.Filename
			violationsReported += "\n\tLine Number: " + strconv.Itoa(v.LineNumber)
			violationsReported += "\n\tCreated At: " + v.CreatedAt + "\n"
		}
	}
	return violationsReported
}
