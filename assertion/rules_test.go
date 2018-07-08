package assertion

import (
	"errors"
	"testing"
)

// TestValueSource provides test values

type TestValueSource struct{}

func (t TestValueSource) GetValue(expression Expression) (string, error) {
	if expression.Value != "" {
		return expression.Value, nil
	}
	return "m3.medium", nil
}

func testValueSource() ValueSource {
	return TestValueSource{}
}

// TestValueSourceWithError simulates errors for value provider

type TestValueSourceWithError struct{}

func (t TestValueSourceWithError) GetValue(expression Expression) (string, error) {
	return "", errors.New("GET_VALUE_ERROR")
}

func testValueSourceWithError() ValueSource {
	return TestValueSourceWithError{}
}

// MockExternalRuleInvoker simulates invocation of external endpoints to get values

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
  - id: TEST3
    message: Test message
    resource: aws_ebs_volume
    severity: WARNING
    assertions:
      - key: size
        op: le
        value: 1000
        value_type: integer
    tags:
      - ebs
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
	if len(r.Rules) != 3 {
		t.Error("Expected to parse 3 rules")
	}
}

type FilterTestCase struct {
	Tags          []string
	Ids           []string
	IgnoreIds     []string
	ExpectedRules []string
}

func TestFilterRules(t *testing.T) {

	var emptyTags []string
	var emptyIds []string

	testCases := map[string]FilterTestCase{
		"allRules": FilterTestCase{emptyTags, emptyIds, emptyIds, []string{"TEST1", "TEST2", "TEST3"}},
		"tags":     FilterTestCase{[]string{"s3"}, emptyIds, emptyIds, []string{"TEST2"}},
		"rules":    FilterTestCase{emptyTags, []string{"TEST1"}, emptyIds, []string{"TEST1"}},
		"both":     FilterTestCase{[]string{"s3"}, []string{"TEST1"}, emptyIds, []string{"TEST1", "TEST2"}},
		"overlap":  FilterTestCase{[]string{"s3"}, []string{"TEST2"}, emptyIds, []string{"TEST2"}},
		"exclude":  FilterTestCase{emptyTags, emptyIds, []string{"TEST1"}, []string{"TEST2", "TEST3"}},
	}
	for k, tc := range testCases {
		r := FilterRulesByTagAndID(MustParseRules(content, t).Rules, tc.Tags, tc.Ids, tc.IgnoreIds)
		if len(r) != len(tc.ExpectedRules) {
			t.Errorf("Expected %s to include %d rules not %d\n", k, len(tc.ExpectedRules), len(r))
		}
	}
}

func TestFilterRulesByTagAndID(t *testing.T) {
	tags := []string{"s3"}
	ids := []string{"TEST3"}
	r := FilterRulesByTagAndID(MustParseRules(content, t).Rules, tags, ids, []string{})
	if len(r) != 2 {
		t.Error("Expected filterRulesByTag to return 2 rules")
	}
	for _, rule := range r {
		if rule.ID != "TEST2" && rule.ID != "TEST3" {
			t.Error("Expected filterRulesByTagAndID to select correct rules")
		}
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
	status, violations, err := CheckRule(rules.Rules[0], resource, mockExternalRuleInvoker())
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
	status, violations, err := CheckRule(rules.Rules[0], resource, mockExternalRuleInvoker())
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
	status, violations, err := CheckRule(rules.Rules[0], resource, mockExternalRuleInvoker())
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
	resolved, violations := ResolveRules(rules.Rules, testValueSource())
	if len(violations) != 0 {
		t.Errorf("Expecting ResolveRules to return 0 violations: %v", violations)
	}
	status, violations, err := CheckRule(resolved[0], resource, mockExternalRuleInvoker())
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

func TestResolveRuleError(t *testing.T) {
	rules := MustParseRules(ruleWithValueFrom, t)
	_, violations := ResolveRules(rules.Rules, testValueSourceWithError())
	if len(violations) != 1 {
		t.Errorf("Expecting ResolveRules to return 1 violations: %v", violations)
	} else {
		ruleID := violations[0].RuleID
		if ruleID != "RULE_RESOLVE" {
			t.Errorf("Expected RULE_RESOLVE violation, not %s", ruleID)
		}
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

func TestInvokeRule(t *testing.T) {
	rules := MustParseRules(ruleWithInvoke, t)
	resource := Resource{
		ID:         "a_test_resource",
		Type:       "aws_instance",
		Properties: map[string]interface{}{"instance_type": "m3.medium"},
		Filename:   "test.tf",
	}
	resolved, _ := ResolveRules(rules.Rules, testValueSource())
	counter := mockExternalRuleInvoker()
	CheckRule(resolved[0], resource, counter)
	if *counter != 1 {
		t.Error("Expecting external rule engine to be invoked")
	}
}
