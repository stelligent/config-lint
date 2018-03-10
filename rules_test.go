package main

import (
	"testing"
)

var content = `Rules:
  - id: TEST1
    message: Test message
    resource: aws_instance
    severity: WARNING
    filters:
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
    filters:
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
	r := filterRulesByTag(MustParseRules(content).Rules, tags)
	if len(r) != 1 {
		t.Error("Expected filterRulesByTag to return 1 rule")
	}
	if r[0].Id != "TEST2" {
		t.Error("Expected filterRulesByTag to select correct rule")
	}
}

func TestFilterRulesById(t *testing.T) {
	ids := []string{"TEST2"}
	r := filterRulesById(MustParseRules(content).Rules, ids)
	if len(r) != 1 {
		t.Error("Expected filterRulesById to return 1 rule")
	}
	if r[0].Id != "TEST2" {
		t.Error("Expected filterRulesById to select correct rule")
	}
}
