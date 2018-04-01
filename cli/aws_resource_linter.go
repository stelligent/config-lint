package main

import (
	"fmt"
	"github.com/stelligent/config-lint/assertion"
)

type (
	// AWSResourceLoader uses the AWS SDK to get resource information about existing resources
	AWSResourceLoader interface {
		Load() ([]assertion.Resource, error)
	}

	// AWSResourceLinter implements a Linter for data returned by the calls to the AWS SDK
	AWSResourceLinter struct {
		Loader AWSResourceLoader
		Log    assertion.LoggingFunction
	}
)

// Validate applies a Ruleset to all SecurityGroups
func (l AWSResourceLinter) Validate(filenames []string, ruleSet assertion.RuleSet, tags []string, ruleIDs []string) ([]string, []assertion.ScannedResource, []assertion.Violation, error) {
	noFilenames := []string{}
	noScannedResources := []assertion.ScannedResource{}
	rules := assertion.FilterRulesByTagAndID(ruleSet.Rules, tags, ruleIDs)
	resources, err := l.Loader.Load()
	if err != nil {
		return noFilenames, noScannedResources, []assertion.Violation{}, err
	}
	r := ResourceLinter{Log: l.Log}
	scannedResources, violations, err := r.ValidateResources(resources, rules)
	return noFilenames, scannedResources, violations, err
}

// Search applies a JMESPath to all SecurityGroups
func (l AWSResourceLinter) Search(filenames []string, ruleSet assertion.RuleSet, searchExpression string) {
	resources, _ := l.Loader.Load()
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
