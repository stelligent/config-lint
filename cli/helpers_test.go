package main

import (
	"github.com/stelligent/config-lint/assertion"
	"testing"
)

func testLogger(string) {}

func loadRulesForTest(filename string, t *testing.T) assertion.RuleSet {
	rulesContent, err := assertion.LoadRules(filename)
	if err != nil {
		t.Errorf("Unable to load rules file: %s", filename)
		return assertion.RuleSet{}
	}
	ruleSet, err := assertion.ParseRules(rulesContent)
	if err != nil {
		t.Errorf("Unable to parse rules file: %s", filename)
		return assertion.RuleSet{}
	}
	return ruleSet
}
