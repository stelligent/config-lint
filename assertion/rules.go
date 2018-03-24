package assertion

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

// LoadRules loads the contents of a YAML file
func LoadRules(filename string) string {
	rules, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(rules)
}

// MustParseRules converts YAML string content to a Result
func MustParseRules(rules string) RuleSet {
	r := RuleSet{}
	err := yaml.Unmarshal([]byte(rules), &r)
	if err != nil {
		panic(err)
	}
	return r
}

// FilterRulesByTag selects a subset of rules based on a tag
func FilterRulesByTag(rules []Rule, tags []string) []Rule {
	filteredRules := make([]Rule, 0)
	for _, rule := range rules {
		if tags == nil || listsIntersect(tags, rule.Tags) {
			filteredRules = append(filteredRules, rule)
		}
	}
	return filteredRules
}

// FilterRulesByID selectes a subset of rules based on ID
func FilterRulesByID(rules []Rule, ruleIDs []string) []Rule {
	if len(ruleIDs) == 0 {
		return rules
	}
	filteredRules := make([]Rule, 0)
	for _, rule := range rules {
		for _, id := range ruleIDs {
			if id == rule.ID {
				filteredRules = append(filteredRules, rule)
			}
		}
	}
	return filteredRules
}

// ResolveRules loads any dynamic values for a collection or rules
func ResolveRules(rules []Rule, valueSource ValueSource, log LoggingFunction) []Rule {
	resolvedRules := make([]Rule, 0)
	for _, rule := range rules {
		resolvedRules = append(resolvedRules, ResolveRule(rule, valueSource, log))
	}
	return resolvedRules
}

// ResolveRule loads any dynamic values for a single Rule
func ResolveRule(rule Rule, valueSource ValueSource, log LoggingFunction) Rule {
	resolvedRule := rule
	resolvedRule.Assertions = make([]Assertion, 0)
	for _, assertion := range rule.Assertions {
		resolvedAssertion := assertion
		resolvedAssertion.Value = valueSource.GetValue(assertion)
		resolvedAssertion.ValueFrom = ValueFrom{}
		resolvedRule.Assertions = append(resolvedRule.Assertions, resolvedAssertion)
	}
	return resolvedRule
}

// CheckRule returns a list of violations for a single Rule applied to a single Resource
func CheckRule(rule Rule, resource Resource, e ExternalRuleInvoker, log LoggingFunction) (string, []Violation) {
	returnStatus := "OK"
	violations := make([]Violation, 0)
	if ExcludeResource(rule, resource) {
		fmt.Println("Ignoring resource:", resource.ID)
		return returnStatus, violations
	}
	if rule.Invoke.URL != "" {
		return e.Invoke(rule, resource)
	}
	for _, ruleAssertion := range rule.Assertions {
		log(fmt.Sprintf("Checking resource %s", resource.ID))
		status := CheckAssertion(rule, ruleAssertion, resource, log)
		if status != "OK" {
			returnStatus = status
			v := Violation{
				RuleID:       rule.ID,
				ResourceID:   resource.ID,
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
