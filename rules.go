package main

import (
	"github.com/ghodss/yaml"
)

func MustParseRules(rules string) Rules {
	r := Rules{}
	err := yaml.Unmarshal([]byte(rules), &r)
	if err != nil {
		panic(err)
	}
	return r
}

func filterRulesByTag(rules []Rule, tags []string) []Rule {
	filteredRules := make([]Rule, 0)
	for _, rule := range rules {
		if tags == nil || listsIntersect(tags, rule.Tags) {
			filteredRules = append(filteredRules, rule)
		}
	}
	return filteredRules
}

func filterRulesById(allRules Rules, ruleIds []string) Rules {
	if len(ruleIds) == 0 {
		return allRules
	}
	filteredRules := make([]Rule, 0)
	for _, rule := range allRules.Rules {
		for _, id := range ruleIds {
			if id == rule.Id {
				filteredRules = append(filteredRules, rule)
			}
		}
	}
	return Rules{Rules: filteredRules}
}
