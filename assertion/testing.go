package assertion

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"testing"
)

type (
	// FixtureTestCases is used to read a set of test cases from a YAML file
	FixtureTestCases struct {
		Description string
		TestCases   []FixtureTestCase `json:"test_cases"`
	}

	// FixtureTestCase describes a single test case
	FixtureTestCase struct {
		Name     string
		Rule     Rule
		Resource Resource
		Result   string
	}
)

// These function are made public so package uses can also use them for writing tests

// TestLogging suppress log message when running tests
func TestLogging(s string) {
}

// FailTestIfError is a helper to check err and call test Error if it is not nil
func FailTestIfError(err error, message string, t *testing.T) {
	if err != nil {
		t.Error(message + ":" + err.Error())
	}
}

// LoadTestCasesFromFixture reads YAML data describing test cases
func LoadTestCasesFromFixture(filename string, t *testing.T) FixtureTestCases {
	var testCases FixtureTestCases
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("Unable to read fixture file: %s", filename)
		return testCases
	}
	err = yaml.Unmarshal(content, &testCases)
	if err != nil {
		t.Errorf("Unable to parse fixture file: %s", filename)
		return testCases
	}
	return testCases
}

// RunTestCasesFromFixture loads a YAML file describing test cases and runs them
func RunTestCasesFromFixture(filename string, t *testing.T) {
	fixture := LoadTestCasesFromFixture(filename, t)
	for _, testCase := range fixture.TestCases {
		status, err := CheckAssertion(testCase.Rule, testCase.Rule.Assertions[0], testCase.Resource, TestLogging)
		FailTestIfError(err, testCase.Name, t)
		if status != testCase.Result {
			t.Errorf("Test case %s returned %s expecting %s", testCase.Name, status, testCase.Result)
		}
	}
}
