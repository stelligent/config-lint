package main

import (
	"github.com/stelligent/config-lint/assertion"
	"testing"
)

type testInvoker struct{}

func (i testInvoker) Invoke(assertion.Rule, assertion.Resource) (string, []assertion.Violation, error) {
	return "TEST", []assertion.Violation{}, nil
}

func TestCheckCompliance(t *testing.T) {
	rules := []assertion.Rule{}
	configurationItem := ConfigurationItem{}
	externalRules := testInvoker{}
	complianceType, err := checkCompliance(rules, configurationItem, externalRules)
	if err != nil {
		t.Errorf("checkCompliance should not return error for empty list of rules")
	}
	if complianceType != "NOT_APPLICABLE" {
		t.Errorf("checkCompliance returned %s instead of NOT_APPLICABLE", complianceType)
	}
}
