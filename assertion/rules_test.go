package assertion

import (
	"testing"
)

type TestValueSource struct {
}

func (t TestValueSource) GetValue(assertion Assertion) string {
	if assertion.Value != "" {
		return assertion.Value
	}
	return "m3.medium"
}

func testValueSource() ValueSource {
	return TestValueSource{}
}

type MockExternalRuleInvoker struct {
	InvokeCount int
}

func mockExternalRuleInvoker() *MockExternalRuleInvoker {
	return &MockExternalRuleInvoker{InvokeCount: 0}
}

func (e *MockExternalRuleInvoker) Invoke(Rule, Resource) (string, []Violation) {
	e.InvokeCount += 1
	noViolations := make([]Violation, 0)
	return "OK", noViolations
}

var content = `Rules:
  - id: TEST1
    message: Test message
    resource: aws_instance
    severity: WARNING
    assertions:
      - type: value
        key: instance_type
        op: in
        value: t2.micro
    tags:
      - ec2
  - id: TEST2
    message: Test message
    resource: aws_s3_bucket
    severity: WARNING
    assertions:
      - type: value
        key: name
        op: eq
        value: bucket1
    tags:
      - s3
`

func TestParseRules(t *testing.T) {
	r := MustParseRules(content)
	if len(r.Rules) != 2 {
		t.Error("Expected to parse 1 rule")
	}
}

func TestFilterRulesByTag(t *testing.T) {
	tags := []string{"s3"}
	r := FilterRulesByTag(MustParseRules(content).Rules, tags)
	if len(r) != 1 {
		t.Error("Expected filterRulesByTag to return 1 rule")
	}
	if r[0].Id != "TEST2" {
		t.Error("Expected filterRulesByTag to select correct rule")
	}
}

func TestFilterRulesById(t *testing.T) {
	ids := []string{"TEST2"}
	r := FilterRulesById(MustParseRules(content).Rules, ids)
	if len(r) != 1 {
		t.Error("Expected filterRulesById to return 1 rule")
	}
	if r[0].Id != "TEST2" {
		t.Error("Expected filterRulesById to select correct rule")
	}
}

var ruleWithMultipleFilters = `Rules:
  - id: TEST1
    message: Test message
    resource: aws_instance
    severity: FAILURE
    assertions:
      - type: value
        key: instance_type
        op: eq
        value: t2.micro
      - type: value
        key: ami
        op: eq
        value: ami-000000
`

func TestRuleWithMultipleFilter(t *testing.T) {
	rules := MustParseRules(ruleWithMultipleFilters)
	resource := Resource{
		Id:         "a_test_resource",
		Type:       "aws_instance",
		Properties: map[string]interface{}{"instance_type": "t2.micro", "ami": "ami-000000"},
		Filename:   "test.tf",
	}
	status, violations := CheckRule(rules.Rules[0], resource, mockExternalRuleInvoker(), testLogging)
	if status != "OK" {
		t.Error("Expecting multiple rule to match")
	}
	if len(violations) != 0 {
		t.Error("Expecting multiple rule to have zero violations")
	}
}

func TestMultipleFiltersWithSingleFailure(t *testing.T) {
	rules := MustParseRules(ruleWithMultipleFilters)
	resource := Resource{
		Id:         "a_test_resource",
		Type:       "aws_instance",
		Properties: map[string]interface{}{"instance_type": "t2.micro", "ami": "ami-111111"},
		Filename:   "test.tf",
	}
	status, violations := CheckRule(rules.Rules[0], resource, mockExternalRuleInvoker(), testLogging)
	if status != "FAILURE" {
		t.Error("Expecting multiple rule to return FAILURE")
	}
	if len(violations) != 1 {
		t.Error("Expecting multiple rule to have one violation")
	}
}

func TestMultipleFiltersWithMultipleFailures(t *testing.T) {
	rules := MustParseRules(ruleWithMultipleFilters)
	resource := Resource{
		Id:         "a_test_resource",
		Type:       "aws_instance",
		Properties: map[string]interface{}{"instance_type": "c3.medium", "ami": "ami-111111"},
		Filename:   "test.tf",
	}
	status, violations := CheckRule(rules.Rules[0], resource, mockExternalRuleInvoker(), testLogging)
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
      - type: value
        key: instance_type
        op: in
        value_from:
          bucket: config-rules-for-lambda
          key: instance-types
`

func TestValueFrom(t *testing.T) {
	rules := MustParseRules(ruleWithValueFrom)
	resource := Resource{
		Id:         "a_test_resource",
		Type:       "aws_instance",
		Properties: map[string]interface{}{"instance_type": "m3.medium"},
		Filename:   "test.tf",
	}
	resolved := ResolveRules(rules.Rules, testValueSource(), testLogging)
	status, violations := CheckRule(resolved[0], resource, mockExternalRuleInvoker(), testLogging)
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
	rules := MustParseRules(ruleWithInvoke)
	resource := Resource{
		Id:         "a_test_resource",
		Type:       "aws_instance",
		Properties: map[string]interface{}{"instance_type": "m3.medium"},
		Filename:   "test.tf",
	}
	resolved := ResolveRules(rules.Rules, testValueSource(), testLogging)
	e := mockExternalRuleInvoker()
	CheckRule(resolved[0], resource, e, testLogging)
	if e.InvokeCount != 1 {
		t.Error("Expecting external rule engine to be invoked")
	}
}
