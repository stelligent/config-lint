package assertion

import (
	"errors"
	"github.com/ghodss/yaml"
)

// ParseRules converts YAML string content to a Result
func ParseRules(rules string) (RuleSet, error) {
	r := RuleSet{}
	err := yaml.Unmarshal([]byte(rules), &r)
	return r, err
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
func FilterRulesByID(rules []Rule, ruleIDs []string, ignoreRuleIDs []string) []Rule {
	if len(ruleIDs) == 0 && len(ignoreRuleIDs) == 0 {
		return rules
	}
	filteredRules := make([]Rule, 0)
	for _, rule := range rules {
		include := false
		for _, id := range ruleIDs {
			if id == rule.ID {
				include = true
			}
		}
		if len(ignoreRuleIDs) > 0 {
			include = true
			for _, id := range ignoreRuleIDs {
				if id == rule.ID {
					include = false
				}
			}
		}
		if include {
			filteredRules = append(filteredRules, rule)
		}
	}
	return filteredRules
}

func uniqueRules(list []Rule) []Rule {
	rules := make([]Rule, 0)
	keys := make(map[string]bool, 0)
	for _, rule := range list {
		if _, ok := keys[rule.ID]; !ok {
			keys[rule.ID] = true
			rules = append(rules, rule)
		}
	}
	return rules
}

// FilterRulesByTagAndID filters by both tag and id
func FilterRulesByTagAndID(rules []Rule, tags []string, ruleIds []string, ignoreRuleIds []string) []Rule {
	if len(tags) == 0 && len(ruleIds) == 0 && len(ignoreRuleIds) == 0 {
		return rules
	}
	if len(tags) == 0 {
		return FilterRulesByID(rules, ruleIds, ignoreRuleIds)
	}
	if len(ruleIds) == 0 {
		return FilterRulesByTag(rules, tags)
	}
	return uniqueRules(append(FilterRulesByID(rules, ruleIds, ignoreRuleIds), FilterRulesByTag(rules, tags)...))
}

// ResolveRules loads any dynamic values for a collection or rules
func ResolveRules(rules []Rule, valueSource ValueSource) ([]Rule, []Violation) {
	resolvedRules := []Rule{}
	violations := []Violation{}
	for _, rule := range rules {
		r, vs := ResolveRule(rule, valueSource)
		resolvedRules = append(resolvedRules, r)
		violations = append(violations, vs...)
	}
	return resolvedRules, violations
}

// ResolveRule loads any dynamic values for a single Rule
func ResolveRule(rule Rule, valueSource ValueSource) (Rule, []Violation) {
	resolvedRule := rule
	resolvedRule.Assertions = []Expression{}
	violations := []Violation{}
	for _, assertion := range rule.Assertions {
		value, err := valueSource.GetValue(assertion)
		if err != nil {
			Debugf("ResolveRule error: %s\n", err.Error())
			violations = append(violations, Violation{
				Category:         "load",
				RuleID:           "RULE_RESOLVE",
				ResourceID:       rule.ID,
				ResourceType:     "rule",
				Status:           "FAILURE",
				RuleMessage:      "Unable to resolve value in rule",
				AssertionMessage: err.Error(),
				CreatedAt:        currentTime(),
			})
		}
		resolvedAssertion := assertion
		resolvedAssertion.Value = value
		resolvedAssertion.ValueFrom = ValueFrom{}
		resolvedRule.Assertions = append(resolvedRule.Assertions, resolvedAssertion)
	}
	return resolvedRule, violations
}

// CheckRule returns a list of violations for a single Rule applied to a single Resource
func CheckRule(rule Rule, resource Resource, e ExternalRuleInvoker) (string, []Violation, error) {
	returnStatus := "OK"
	violations := make([]Violation, 0)
	if ExcludeResource(rule, resource) {
		Debugf("Ignoring resource: %s", resource.ID)
		return returnStatus, violations, nil
	}
	if rule.Invoke.URL != "" {
		return e.Invoke(rule, resource)
	}
	match, err := andExpression(rule.Conditions, resource)
	if err != nil {
		return "FAILURE", violations, err
	}
	if !match.Match {
		return returnStatus, violations, nil
	}
	for _, ruleAssertion := range rule.Assertions {
		Debugf("Checking Category: %s, Type: %s, Id: %s\n", resource.Category, resource.Type, resource.ID)
		expressionResult, err := CheckExpression(rule, ruleAssertion, resource)
		if err != nil {
			return "FAILURE", violations, err
		}
		if expressionResult.Status != "OK" {
			returnStatus = expressionResult.Status
			v := Violation{
				RuleID:           rule.ID,
				ResourceID:       resource.ID,
				ResourceType:     resource.Type,
				Category:         resource.Category,
				Status:           expressionResult.Status,
				RuleMessage:      rule.Message,
				AssertionMessage: expressionResult.Message,
				Filename:         resource.Filename,
				LineNumber:       resource.LineNumber,
				CreatedAt:        currentTime(),
			}
			violations = append(violations, v)
		}
	}
	return returnStatus, violations, nil
}

// Join two RuleSets together
func JoinRuleSets(firstSet RuleSet, secondSet RuleSet) (RuleSet, error) {
	// if one of the sets is empty, return the other
	// if both are empty, an empty set is returned
	if len(firstSet.Rules) == 0 {
		return secondSet, nil
	} else if len(secondSet.Rules) == 0 {
		return firstSet, nil
	}

	// RuleSets must match Type and Version
	// Description will be taken from the first given rule set
	if firstSet.Type != secondSet.Type || firstSet.Version != secondSet.Version {
		return firstSet, errors.New("RuleSet Type and Version must match")
	} else {
		joinedSet := RuleSet{}
		joinedSet.Type = firstSet.Type
		joinedSet.Description = firstSet.Description
		joinedSet.Files = append(firstSet.Files, secondSet.Files...)
		joinedSet.Rules = append(firstSet.Rules, secondSet.Rules...)
		joinedSet.Version = firstSet.Version
		joinedSet.Resources = append(firstSet.Resources, secondSet.Resources...)
		joinedSet.Columns = append(firstSet.Columns, secondSet.Columns...)
		return joinedSet, nil
	}
}
