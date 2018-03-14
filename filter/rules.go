package filter

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

func LoadRules(filename string) string {
	rules, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(rules)
}

func MustParseRules(rules string) RuleSet {
	r := RuleSet{}
	err := yaml.Unmarshal([]byte(rules), &r)
	if err != nil {
		panic(err)
	}
	return r
}

func FilterRulesByTag(rules []Rule, tags []string) []Rule {
	filteredRules := make([]Rule, 0)
	for _, rule := range rules {
		if tags == nil || listsIntersect(tags, rule.Tags) {
			filteredRules = append(filteredRules, rule)
		}
	}
	return filteredRules
}

func FilterRulesById(rules []Rule, ruleIds []string) []Rule {
	if len(ruleIds) == 0 {
		return rules
	}
	filteredRules := make([]Rule, 0)
	for _, rule := range rules {
		for _, id := range ruleIds {
			if id == rule.Id {
				filteredRules = append(filteredRules, rule)
			}
		}
	}
	return filteredRules
}

func ResolveRules(rules []Rule, valueSource ValueSource, log LoggingFunction) []Rule {
	resolvedRules := make([]Rule, 0)
	for _, rule := range rules {
		resolvedRules = append(resolvedRules, ResolveRule(rule, valueSource, log))
	}
	return resolvedRules
}

func ResolveRule(rule Rule, valueSource ValueSource, log LoggingFunction) Rule {
	resolvedRule := rule
	resolvedRule.Filters = make([]Filter, 0)
	for _, filter := range rule.Filters {
		resolvedFilter := filter
		resolvedFilter.Value = valueSource.GetValue(filter)
		resolvedFilter.ValueFrom = FilterValueFrom{}
		resolvedRule.Filters = append(resolvedRule.Filters, resolvedFilter)
	}
	return resolvedRule
}

func ApplyRule(rule Rule, resource Resource, log LoggingFunction) (string, []Violation) {
	returnStatus := "OK"
	violations := make([]Violation, 0)
	if ExcludeResource(rule, resource) {
		return returnStatus, violations
	}
	for _, ruleFilter := range rule.Filters {
		log(fmt.Sprintf("Checking resource %s", resource.Id))
		status := ApplyFilter(rule, ruleFilter, resource, log)
		if status != "OK" {
			returnStatus = status
			v := Violation{
				RuleId:       rule.Id,
				ResourceId:   resource.Id,
				ResourceType: resource.Type,
				Status:       status,
				Message:      rule.Message,
				Filename:     resource.Filename,
			}
			violations = append(violations, v)
		}
	}
	return returnStatus, violations
}

func ResolveValuesInFilters(filters []Filter, valueSource ValueSource, log LoggingFunction) []Filter {
	resolved := make([]Filter, 0)
	for _, filter := range filters {
		resolvedFilter := filter
		resolvedFilter.Value = valueSource.GetValue(filter)
		resolvedFilter.ValueFrom = FilterValueFrom{}
		resolved = append(resolved, resolvedFilter)
	}
	return resolved
}
