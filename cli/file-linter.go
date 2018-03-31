package main

import (
	"fmt"
	"github.com/stelligent/config-lint/assertion"
)

// FileLinter provides implmenation for some common functions that are used by multiple Linter implementations
type FileLinter struct {
	Log assertion.LoggingFunction
}

// ValidateResources evaluates a list of Rule objects to a list of Resource objects
func (l FileLinter) ValidateResources(resources []assertion.Resource, rules []assertion.Rule) ([]assertion.Violation, error) {

	valueSource := assertion.StandardValueSource{Log: l.Log}
	resolvedRules := assertion.ResolveRules(rules, valueSource, l.Log)
	externalRules := assertion.StandardExternalRuleInvoker{Log: l.Log}

	allViolations := make([]assertion.Violation, 0)
	for _, rule := range resolvedRules {
		l.Log(fmt.Sprintf("Rule %s: %s", rule.ID, rule.Message))
		for _, resource := range assertion.FilterResourcesByType(resources, rule.Resource) {
			if assertion.ExcludeResource(rule, resource) {
				l.Log(fmt.Sprintf("Ignoring resource %s", resource.ID))
			} else {
				_, violations, err := assertion.CheckRule(rule, resource, externalRules, l.Log)
				if err != nil {
					return allViolations, err
				}
				allViolations = append(allViolations, violations...)
			}
		}
	}
	return allViolations, nil
}

// ValidateFiles validates a collection of filenames using a RuleSet
func (l FileLinter) ValidateFiles(filenames []string, ruleSet assertion.RuleSet, tags []string, ruleIDs []string, loader ResourceLoader) ([]string, []assertion.Violation, error) {
	rules := assertion.FilterRulesByTagAndID(ruleSet.Rules, tags, ruleIDs)
	allViolations := make([]assertion.Violation, 0)
	filesScanned := make([]string, 0)
	for _, filename := range filenames {
		include, _ := assertion.ShouldIncludeFile(ruleSet.Files, filename) // FIXME what about error?
		if include {
			l.Log(fmt.Sprintf("Processing %s", filename))
			resources, err := loader.Load(filename)
			if err != nil {
				return filesScanned, allViolations, err
			}
			violations, err := l.ValidateResources(resources, rules)
			if err != nil {
				return filesScanned, allViolations, err
			}
			allViolations = append(allViolations, violations...)
			filesScanned = append(filesScanned, filename)
		}
	}
	return filesScanned, allViolations, nil
}

// SearchFiles evaluates a JMESPath expression against resources in a collection of filenames
func (l FileLinter) SearchFiles(filenames []string, ruleSet assertion.RuleSet, searchExpression string, loader ResourceLoader) {
	for _, filename := range filenames {
		include, _ := assertion.ShouldIncludeFile(ruleSet.Files, filename) // FIXME what about error?
		if include {
			fmt.Printf("Searching %s:\n", filename)
			resources, err := loader.Load(filename)
			if err != nil {
				fmt.Println("Error for file:", filename)
				fmt.Println(err.Error())
			}
			for _, resource := range resources {
				v, err := assertion.SearchData(searchExpression, resource.Properties)
				if err != nil {
					fmt.Println(err)
				} else {
					s, err := assertion.JSONStringify(v)
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Printf("%s: %s\n", resource.ID, s)
					}
				}
			}
		}
	}
}
