package main

import (
	"fmt"
	"github.com/stelligent/config-lint/assertion"
)

// Linter provides the interface for all supported linters
type Linter interface {
	Validate(filenames []string, ruleSet assertion.RuleSet, tags []string, ruleIDs []string) ([]string, []assertion.Violation, error)
	Search(filenames []string, ruleSet assertion.RuleSet, searchExpression string)
}

// ResourceLoader provides the interface that a Linter needs to load a collection of Resource objects
type ResourceLoader interface {
	Load(filename string) ([]assertion.Resource, error)
}

// BaseLinter provides implmenation for some common functions that are used by multiple Linter implementations
type BaseLinter struct {
}

// ValidateResources evaluates a list of Rule objects to a list of Resource objects
func (l BaseLinter) ValidateResources(resources []assertion.Resource, rules []assertion.Rule, log assertion.LoggingFunction) ([]assertion.Violation, error) {

	valueSource := assertion.StandardValueSource{Log: log}
	resolvedRules := assertion.ResolveRules(rules, valueSource, log)
	externalRules := assertion.StandardExternalRuleInvoker{Log: log}

	allViolations := make([]assertion.Violation, 0)
	for _, rule := range resolvedRules {
		log(fmt.Sprintf("Rule %s: %s", rule.ID, rule.Message))
		for _, resource := range assertion.FilterResourcesByType(resources, rule.Resource) {
			if assertion.ExcludeResource(rule, resource) {
				log(fmt.Sprintf("Ignoring resource %s", resource.ID))
			} else {
				_, violations, err := assertion.CheckRule(rule, resource, externalRules, log)
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
func (l BaseLinter) ValidateFiles(filenames []string, ruleSet assertion.RuleSet, tags []string, ruleIDs []string, loader ResourceLoader, log assertion.LoggingFunction) ([]string, []assertion.Violation, error) {
	rules := assertion.FilterRulesByTagAndID(ruleSet.Rules, tags, ruleIDs)
	allViolations := make([]assertion.Violation, 0)
	filesScanned := make([]string, 0)
	for _, filename := range filenames {
		include, _ := assertion.ShouldIncludeFile(ruleSet.Files, filename) // FIXME what about error?
		if include {
			log(fmt.Sprintf("Processing %s", filename))
			resources, err := loader.Load(filename)
			if err != nil {
				return filesScanned, allViolations, err
			}
			violations, err := l.ValidateResources(resources, rules, log)
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
func (l BaseLinter) SearchFiles(filenames []string, ruleSet assertion.RuleSet, searchExpression string, loader ResourceLoader) {
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

func makeLinter(linterType string, log assertion.LoggingFunction) Linter {
	switch linterType {
	case "Kubernetes":
		return KubernetesLinter{Log: log}
	case "Terraform":
		return TerraformLinter{Log: log}
	case "SecurityGroup":
		return SecurityGroupLinter{Log: log}
	case "LintRules":
		return RulesLinter{Log: log}
	default:
		fmt.Printf("Type not supported: %s\n", linterType)
		return nil
	}
}
