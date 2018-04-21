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
