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
	Id       string
	Resource string
	Message  string
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

/*
* Given a rule file and a rule ID
* Return a RuleSet containing that single named Rule from the Rule File
 */
func loadSingleRule(ruleFile string, ruleId string) (assertion.RuleSet, error) {
	// TODO only get the name rule from the given rule file as a ruleSet

	ruleSet, err := loadBuiltInRuleSet(ruleFile)
	return ruleSet, err
}

/*
* Given a resource file and a resource ID
* Create a temporary file containing that single named resource
 */
func createTestResourceFile(filename string, resourceId string) (string, error) {
	// _, err := ioutil.ReadFile(filename)
	// if err != nil {
	// 	assertion.Debugf("Failed to load Test Suite file: %v\n", err)
	// 	return "", err
	// }

	return filename, nil
}

// Run each test case in a test suite
func runTestSuite(t *testing.T, ts TestSuite) {
	for _, tc := range ts.Tests {
		assertion.Debugf("Running Test: %v\n", tc.Id)

		// Load only the rule for this test case
		ruleConfigPath := strings.Split(ts.RootPath, "/cli/")[1] + "/rule.yml"
		ruleSet, err := loadSingleRule(ruleConfigPath, tc.Id)
		if err != nil {
			assert.Nil(t, err, "Cannot load built-in Terraform rule")
		}

		// Load only the resource to be tested
		for _, tag := range tc.Tags {
			testResourceDirectory := strings.Split(ts.RootPath, "/cli/")[1] + "/tests/" + tag + "/" + "container_properties_privileged.tf"
			assertion.Debugf("test resource dir %v\n", testResourceDirectory)
			testResourcePath, err := createTestResourceFile(testResourceDirectory, tc.Resource)

			options := linter.Options{
				RuleIDs: []string{tc.Id},
			}
			vs := assertion.StandardValueSource{}

			// validate the rule set
			var l linter.Linter
			if tag == "terraform11" {
				// Defining 'tf11' for the Parser type
				l, err = linter.NewLinter(ruleSet, vs, []string{testResourcePath}, "tf11")
			} else {
				l, err = linter.NewLinter(ruleSet, vs, []string{testResourcePath}, "tf12")
			}
			report, err := l.Validate(ruleSet, options)
			assert.Nil(t, err, "Validate failed for file")

			warningViolationsReported := getViolationsString("WARNING", report.Violations)
			warningMessage := fmt.Sprintf("Expecting %d warnings for RuleID %s in File %s:\n %s", tc.Warnings, tc.Id, tc.Resource, warningViolationsReported)
			assert.Equal(t, tc.Warnings, numberOfWarnings(report.Violations), warningMessage)

			failureViolationsReported := getViolationsString("FAILURE", report.Violations)
			failureMessage := fmt.Sprintf("Expecting %d failures for RuleID %s in File %s:\n %s", tc.Failures, tc.Id, tc.Resource, failureViolationsReported)
			assert.Equal(t, tc.Failures, numberOfFailures(report.Violations), failureMessage)
		}
	}
}

/*
* Given resource type and tags
* Run all defined tests per resource type and subtype
 */
func RunBuiltinTests(t *testing.T, resourceType string, subtypes []string) {
	assertion.SetDebug(true)
	assertion.Debugf("#######################\n\n\n%v\n\n", resourceType)

	// Get a list of all test cases
	box := packr.NewBox("./assets/" + resourceType)
	filesInBox := box.List()
	for _, file := range filesInBox {
		if isTestCase(file) {
			absolutePath, _ := filepath.Abs("./assets/" + resourceType + "/" + file)
			assertion.Debugf("Loading test case: %v\n", absolutePath)
			ts, err := loadTestSuite(absolutePath)
			if err != nil {
				assert.Nil(t, err, "Cannot load test case")
			}
			runTestSuite(t, ts)
		}
	}

	assertion.SetDebug(false)
	assert.NotNil(t, nil) // fail
}

// Run built in rules against Terraform v0.11 parser
func TestTerraform11BuiltInRules(t *testing.T) {
	RunBuiltinTests(t, "terraform", []string{"terraform11", "terraform12"})
}

func TestTerraform12BuiltInRules(t *testing.T) {

	// Define file to load rules from
	// This file is located under cli/assets/
	ruleSet, err := loadBuiltInRuleSet("terraform")
	if err != nil {
		assert.Nil(t, err, "Cannot load built-in Terraform rules")
	}

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
