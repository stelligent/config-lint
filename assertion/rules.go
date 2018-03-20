package assertion

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
	resolvedRule.Assertions = make([]Assertion, 0)
	for _, assertion := range rule.Assertions {
		resolvedAssertion := assertion
		resolvedAssertion.Value = valueSource.GetValue(assertion)
		resolvedAssertion.ValueFrom = AssertionValueFrom{}
		resolvedRule.Assertions = append(resolvedRule.Assertions, resolvedAssertion)
	}
	return resolvedRule
}

func CheckRule(rule Rule, resource Resource, e ExternalRuleInvoker, log LoggingFunction) (string, []Violation) {
	returnStatus := "OK"
	violations := make([]Violation, 0)
	if ExcludeResource(rule, resource) {
		fmt.Println("Ignoring resource:", resource.Id)
		return returnStatus, violations
	}
	if rule.Invoke.Url != "" {
		return e.Invoke(rule, resource)
	}
	for _, ruleAssertion := range rule.Assertions {
		log(fmt.Sprintf("Checking resource %s", resource.Id))
		status := CheckAssertion(rule, ruleAssertion, resource, log)
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

func ResolveValuesInAssertions(assertions []Assertion, valueSource ValueSource, log LoggingFunction) []Assertion {
	resolved := make([]Assertion, 0)
	for _, assertion := range assertions {
		resolvedAssertion := assertion
		resolvedAssertion.Value = valueSource.GetValue(assertion)
		resolvedAssertion.ValueFrom = AssertionValueFrom{}
		resolved = append(resolved, resolvedAssertion)
	}
	return resolved
}
