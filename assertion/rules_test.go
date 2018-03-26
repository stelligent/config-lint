package assertion

import (
	"testing"
)

type TestValueSource struct {
}

func (t TestValueSource) GetValue(assertion Assertion) (string, error) {
	if assertion.Value != "" {
		return assertion.Value, nil
	}
	return "m3.medium", nil
}

func testValueSource() ValueSource {
	return TestValueSource{}
}

type MockExternalRuleInvoker int

func mockExternalRuleInvoker() *MockExternalRuleInvoker {
	var m MockExternalRuleInvoker
	return &m
}

func (e *MockExternalRuleInvoker) Invoke(Rule, Resource) (string, []Violation, error) {
	*e++
	noViolations := make([]Violation, 0)
	return "OK", noViolations, nil
}

var content = `Rules:
  - id: TEST1
    message: Test message
    resource: aws_instance
    severity: WARNING
    assertions:
      - key: instance_type
        op: in
        value: t2.micro
    tags:
      - ec2
  - id: TEST2
    message: Test message
    resource: aws_s3_bucket
    severity: WARNING
    assertions:
      - key: name
        op: eq
        value: bucket1
    tags:
      - s3
`

func MustParseRules(content string, t *testing.T) RuleSet {
	r, err := ParseRules(content)
	if err != nil {
		t.Error("Unable to parse:" + content)
	}
	return r
}

func TestParseRules(t *testing.T) {
	r := MustParseRules(content, t)
	if len(r.Rules) != 2 {
		t.Error("Expected to parse 1 rule")
	}
}

func TestFilterRulesByTag(t *testing.T) {
	tags := []string{"s3"}
	r := FilterRulesByTag(MustParseRules(content, t).Rules, tags)
	if len(r) != 1 {
		t.Error("Expected filterRulesByTag to return 1 rule")
	}
	if r[0].ID != "TEST2" {
		t.Error("Expected filterRulesByTag to select correct rule")
	}
}

func TestFilterRulesByID(t *testing.T) {
	ids := []string{"TEST2"}
	r := FilterRulesByID(MustParseRules(content, t).Rules, ids)
	if len(r) != 1 {
		t.Error("Expected filterRulesByID to return 1 rule")
	}
	if r[0].ID != "TEST2" {
		t.Error("Expected filterRulesByID to select correct rule")
	}
}

var ruleWithMultipleFilters = `Rules:
  - id: TEST1
    message: Test message
    resource: aws_instance
    severity: FAILURE
    assertions:
      - key: instance_type
        op: eq
        value: t2.micro
      - key: ami
        op: eq
        value: ami-000000
`

func TestRuleWithMultipleFilter(t *testing.T) {
	rules := MustParseRules(ruleWithMultipleFilters, t)
	resource := Resource{
		ID:         "a_test_resource",
		Type:       "aws_instance",
		Properties: map[string]interface{}{"instance_type": "t2.micro", "ami": "ami-000000"},
		Filename:   "test.tf",
	}
	status, violations, err := CheckRule(rules.Rules[0], resource, mockExternalRuleInvoker(), testLogging)
	if err != nil {
		t.Error("Error in CheckRule:" + err.Error())
	}
	if status != "OK" {
		t.Error("Expecting multiple rule to match")
	}
	if len(violations) != 0 {
		t.Error("Expecting multiple rule to have zero violations")
	}
}

func TestMultipleFiltersWithSingleFailure(t *testing.T) {
	rules := MustParseRules(ruleWithMultipleFilters, t)
	resource := Resource{
		ID:         "a_test_resource",
		Type:       "aws_instance",
		Properties: map[string]interface{}{"instance_type": "t2.micro", "ami": "ami-111111"},
		Filename:   "test.tf",
	}
	status, violations, err := CheckRule(rules.Rules[0], resource, mockExternalRuleInvoker(), testLogging)
	if err != nil {
		t.Error("Error in CheckRule:" + err.Error())
	}
	if status != "FAILURE" {
		t.Error("Expecting multiple rule to return FAILURE")
	}
	if len(violations) != 1 {
		t.Error("Expecting multiple rule to have one violation")
	}
}

func TestMultipleFiltersWithMultipleFailures(t *testing.T) {
	rules := MustParseRules(ruleWithMultipleFilters, t)
	resource := Resource{
		ID:         "a_test_resource",
		Type:       "aws_instance",
		Properties: map[string]interface{}{"instance_type": "c3.medium", "ami": "ami-111111"},
		Filename:   "test.tf",
	}
	status, violations, err := CheckRule(rules.Rules[0], resource, mockExternalRuleInvoker(), testLogging)
	if err != nil {
		t.Error("Error in CheckRule:" + err.Error())
	}
	if status != "FAILURE" {
		t.Error("Expecting multiple rule to return FAILURE")
	}
	if len(violations) != 2 {
		t.Error("Expecting multiple rule to have two violations")
	}
}

var ruleWithValueFrom = `Rules:
  - id: FROM1
    message: Test value_from
    severity: FAILURE
    resource: aws_instance
    assertions:
      - key: instance_type
        op: in
        value_from:
          bucket: config-rules-for-lambda
          key: instance-types
`

func TestValueFrom(t *testing.T) {
	rules := MustParseRules(ruleWithValueFrom, t)
	resource := Resource{
		ID:         "a_test_resource",
		Type:       "aws_instance",
		Properties: map[string]interface{}{"instance_type": "m3.medium"},
		Filename:   "test.tf",
	}
	resolved := ResolveRules(rules.Rules, testValueSource(), testLogging)
	status, violations, err := CheckRule(resolved[0], resource, mockExternalRuleInvoker(), testLogging)
	if err != nil {
		t.Error("Error in CheckRule:" + err.Error())
	}
	if status != "OK" {
		t.Error("Expecting value_from to match")
	}
	if len(violations) != 0 {
		t.Error("Expecting value_from test to have 0 violations")
	}
}

var ruleWithInvoke = `Rules:
  - id: FROM1
    message: Test value_from
    severity: FAILURE
    resource: aws_instance
    invoke:
      url: http://localhost
`

func TestInvoke(t *testing.T) {
	rules := MustParseRules(ruleWithInvoke, t)
	resource := Resource{
		ID:         "a_test_resource",
		Type:       "aws_instance",
		Properties: map[string]interface{}{"instance_type": "m3.medium"},
		Filename:   "test.tf",
	}
	resolved := ResolveRules(rules.Rules, testValueSource(), testLogging)
	counter := mockExternalRuleInvoker()
	CheckRule(resolved[0], resource, counter, testLogging)
	if *counter != 1 {
		t.Error("Expecting external rule engine to be invoked")
	}
}
