package main

import (
	"fmt"
	"github.com/stelligent/config-lint/assertion"
)

// Linter provides the interface for all supported linters
type Linter interface {
	Validate(filenames []string, ruleSet assertion.RuleSet, tags []string, ruleIDs []string) ([]string, []assertion.Violation)
	Search(filenames []string, ruleSet assertion.RuleSet, searchExpression string)
}

// ResourceLoader provides the interface that a Linter needs to load a collection of Resource objects
type ResourceLoader interface {
	Load(filename string) []assertion.Resource
}

// BaseLinter provides implmenation for some common functions that are used by multiple Linter implementations
type BaseLinter struct {
}

// ValidateResources evaluates a list of Rule objects to a list of Resource objects
func (l BaseLinter) ValidateResources(resources []assertion.Resource, rules []assertion.Rule, tags []string, log assertion.LoggingFunction) []assertion.Violation {

	valueSource := assertion.StandardValueSource{Log: log}
	filteredRules := assertion.FilterRulesByTag(rules, tags)
	resolvedRules := assertion.ResolveRules(filteredRules, valueSource, log)
	externalRules := assertion.StandardExternalRuleInvoker{Log: log}

	allViolations := make([]assertion.Violation, 0)
	for _, rule := range resolvedRules {
		log(fmt.Sprintf("Rule %s: %s", rule.ID, rule.Message))
		for _, resource := range assertion.FilterResourcesByType(resources, rule.Resource) {
			if assertion.ExcludeResource(rule, resource) {
				log(fmt.Sprintf("Ignoring resource %s", resource.ID))
			} else {
				_, violations := assertion.CheckRule(rule, resource, externalRules, log)
				allViolations = append(allViolations, violations...)
			}
		}
	}
	return allViolations
}

// ValidateFiles validates a collection of filenames using a RuleSet
func (l BaseLinter) ValidateFiles(filenames []string, ruleSet assertion.RuleSet, tags []string, ruleIDs []string, loader ResourceLoader, log assertion.LoggingFunction) ([]string, []assertion.Violation) {
	rules := assertion.FilterRulesByID(ruleSet.Rules, ruleIDs)
	allViolations := make([]assertion.Violation, 0)
	filesScanned := make([]string, 0)
	for _, filename := range filenames {
		if assertion.ShouldIncludeFile(ruleSet.Files, filename) {
			log(fmt.Sprintf("Processing %s", filename))
			resources := loader.Load(filename)
			violations := l.ValidateResources(resources, rules, tags, log)
			allViolations = append(allViolations, violations...)
			filesScanned = append(filesScanned, filename)
		}
	}
	return filesScanned, allViolations
}

// SearchFiles evaluates a JMESPath expression against resources in a collection of filenames
func (l BaseLinter) SearchFiles(filenames []string, ruleSet assertion.RuleSet, searchExpression string, loader ResourceLoader) {
	for _, filename := range filenames {
		if assertion.ShouldIncludeFile(ruleSet.Files, filename) {
			fmt.Printf("Searching %s:\n", filename)
			resources := loader.Load(filename)
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

func makeLinter(linterType string, log assertion.LoggingFunction) Linter {
	switch linterType {
	case "Kubernetes":
		return KubernetesLinter{Log: log}
	case "Terraform":
		return TerraformLinter{Log: log}
	case "SecurityGroup":
		return SecurityGroupLinter{Log: log}
	default:
		fmt.Printf("Type not supported: %s\n", linterType)
		return nil
	}
}
