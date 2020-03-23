package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/gobuffalo/packr"
	"github.com/stelligent/config-lint/assertion"
	"github.com/stelligent/config-lint/linter"
	"github.com/stretchr/testify/assert"
)

type TestSuite struct {
	Version     string
	Description string
	Type        string
	Files       []string
	Tests       []TestCase
	RootPath    string
}

type TestCase struct {
	RuleId   string
	Warnings int
	Failures int
	Tags     []string
}

/*
*  Determine if the given filename is a test case
 */
func isTestCase(filename string) bool {
	if strings.Contains(filename, "test") && isYamlFile(filename) {
		return true
	} else {
		return false
	}
}

/*
* Given filepath to a YAML test suite, return a TestSuite object
 */
func loadTestSuite(filename string) (TestSuite, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		assertion.Debugf("Failed to load Test Suite file: %v\n", err)
		return TestSuite{}, err
	}
	ts := TestSuite{}
	err = yaml.Unmarshal([]byte(yamlFile), &ts)
	if err != nil {
		assertion.Debugf("Failed to unmarshall YAML Test Suite: %v\n", err)
		return TestSuite{}, err
	}

	ts.RootPath = strings.TrimRight(filename, "tests/test.yml")

	return ts, nil
}

func getTestResources(directory string) ([]string, error) {
	var testResources []string
	box := packr.NewBox(directory)
	filesInBox := box.List()
	configPatterns := []string{"*.tf"}
	for _, f := range filesInBox {
		match, _ := assertion.ShouldIncludeFile(configPatterns, f)
		if match {
			absolutePath, err := filepath.Abs(directory + "/" + f)
			if err != nil {
				return []string{}, err
			}
			testResources = append(testResources, absolutePath)
		}
	}
	return testResources, nil
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// Run each test case in a test suite
func runTestSuite(t *testing.T, ts TestSuite) {
	// Load only the rule for this test suite
	ruleConfigPath := strings.Split(ts.RootPath, "config-lint/cli/assets/")[1] + "/rule.yml"
	ruleSet, err := loadBuiltInRuleSet(ruleConfigPath)
	if err != nil {
		assert.Nil(t, err, "Cannot load built-in Terraform rule")
	}

	for _, tc := range ts.Tests {
		options := linter.Options{
			RuleIDs: []string{tc.RuleId},
		}
		vs := assertion.StandardValueSource{}

		// validate the rule set
		if contains(tc.Tags, "terraform11") {
			// Load the test resources for this test suite
			testResourceDirectory := strings.Split(ts.RootPath, "config-lint/cli/")[1] + "/tests/terraform11/"
			testResources, err := getTestResources(testResourceDirectory)
			if err != nil {
				assert.Nil(t, err, "Cannot load built-in Terraform 11 test resources")

			}
			// Defining 'tf11' for the Parser type
			l, err := linter.NewLinter(ruleSet, vs, testResources, "tf11")

			report, err := l.Validate(ruleSet, options)
			assert.Nil(t, err, "Validate failed for file")

			warningViolationsReported := getViolationsString("WARNING", report.Violations)
			warningMessage := fmt.Sprintf("Expecting %d warnings for rule %s:\n %s", tc.Warnings, tc.RuleId, warningViolationsReported)
			assert.Equal(t, tc.Warnings, numberOfWarnings(report.Violations), warningMessage)

			failureViolationsReported := getViolationsString("FAILURE", report.Violations)
			failureMessage := fmt.Sprintf("Expecting %d failures for rule %s:\n %s", tc.Failures, tc.RuleId, failureViolationsReported)
			assert.Equal(t, tc.Failures, numberOfFailures(report.Violations), failureMessage)
		}

		if contains(tc.Tags, "terraform12") {
			// Load the test resources for this test suite
			testResourceDirectory := strings.Split(ts.RootPath, "config-lint/cli/")[1] + "/tests/terraform12/"
			testResources, err := getTestResources(testResourceDirectory)
			if err != nil {
				assert.Nil(t, err, "Cannot load built-in Terraform 12 test resources")

			}
			// Defining 'tf11' for the Parser type
			l, err := linter.NewLinter(ruleSet, vs, testResources, "tf12")

			report, err := l.Validate(ruleSet, options)
			assert.Nil(t, err, "Validate failed for file")

			warningViolationsReported := getViolationsString("WARNING", report.Violations)
			warningMessage := fmt.Sprintf("Expecting %d warnings for rule %s:\n %s", tc.Warnings, tc.RuleId, warningViolationsReported)
			assert.Equal(t, tc.Warnings, numberOfWarnings(report.Violations), warningMessage)

			failureViolationsReported := getViolationsString("FAILURE", report.Violations)
			failureMessage := fmt.Sprintf("Expecting %d failures for rule %s:\n %s", tc.Failures, tc.RuleId, failureViolationsReported)
			assert.Equal(t, tc.Failures, numberOfFailures(report.Violations), failureMessage)
		}
	}
}

/*
* Given resource type and tags
* Run all defined tests per resource type and subtype
 */
func RunBuiltinTests(t *testing.T, resourceType string) {
	// Get a list of all test cases
	box := packr.NewBox("./assets/" + resourceType)
	filesInBox := box.List()
	for _, file := range filesInBox {
		if isTestCase(file) {
			absolutePath, _ := filepath.Abs("./assets/" + resourceType + "/" + file)
			ts, err := loadTestSuite(absolutePath)
			if err != nil {
				assert.Nil(t, err, "Cannot load test case")
			}
			runTestSuite(t, ts)
		}
	}
}

// Run built in rules against Terraform parser
func TestTerraformBuiltInRules(t *testing.T) {
	RunBuiltinTests(t, "terraform")
}
