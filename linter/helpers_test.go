package linter

import (
	"github.com/stelligent/config-lint/assertion"
	"io/ioutil"
	"testing"
)

func loadRulesForTest(filename string, t *testing.T) assertion.RuleSet {
	rulesContent, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("Unable to load rules file: %s", filename)
		return assertion.RuleSet{}
	}
	ruleSet, err := assertion.ParseRules(string(rulesContent))
	if err != nil {
		t.Errorf("Unable to parse rules file: %s", filename)
		return assertion.RuleSet{}
	}
	return ruleSet
}

func assertViolationsCount(testName string, count int, violations []assertion.Violation, t *testing.T) {
	if len(violations) != count {
		t.Errorf("TestTerraformDataObject returned %d violations, expecting %d", len(violations), count)
		t.Errorf("Violations: %v", violations)
	}
}

func assertViolationByRuleID(testName string, ruleID string, violations []assertion.Violation, t *testing.T) {
	found := false
	ruleIDsFound := []string{}
	for _, v := range violations {
		ruleIDsFound = append(ruleIDsFound, v.RuleID)
		if v.RuleID == ruleID {
			found = true
		}
	}
	if !found {
		t.Errorf("%s expected violation %s not found in %v", testName, ruleID, ruleIDsFound)
	}
}
