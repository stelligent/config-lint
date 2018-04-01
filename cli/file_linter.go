package main

import (
	"fmt"
	"github.com/stelligent/config-lint/assertion"
)

// FileLinter provides implementation for some common functions that are used by multiple Linter implementations
type FileLinter struct {
	Log assertion.LoggingFunction
}

// ValidateFiles validates a collection of filenames using a RuleSet
func (l FileLinter) ValidateFiles(filenames []string, ruleSet assertion.RuleSet, tags []string, ruleIDs []string, loader ResourceLoader) ([]string, []assertion.ScannedResource, []assertion.Violation, error) {
	rules := assertion.FilterRulesByTagAndID(ruleSet.Rules, tags, ruleIDs)
	allViolations := make([]assertion.Violation, 0)
	filesScanned := make([]string, 0)
	resourcesScanned := make([]assertion.ScannedResource, 0)
	r := ResourceLinter{Log: l.Log}
	for _, filename := range filenames {
		include, err := assertion.ShouldIncludeFile(ruleSet.Files, filename)
		if err == nil && include {
			l.Log(fmt.Sprintf("Processing %s", filename))
			resources, err := loader.Load(filename)
			if err != nil {
				return filesScanned, resourcesScanned, allViolations, err
			}
			scanned, violations, err := r.ValidateResources(resources, rules)
			resourcesScanned = append(resourcesScanned, scanned...)
			if err != nil {
				return filesScanned, resourcesScanned, allViolations, err
			}
			allViolations = append(allViolations, violations...)
			filesScanned = append(filesScanned, filename)
		}
	}
	return filesScanned, resourcesScanned, allViolations, nil
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
						fmt.Printf("%s (%s): %s\n", resource.ID, resource.Type, s)
					}
				}
			}
		}
	}
}
