package main

import (
	"fmt"
	"github.com/stelligent/config-lint/assertion"
)

// FileLinter provides implementation for some common functions that are used by multiple Linter implementations
type FileLinter struct {
	Filenames []string
	Log       assertion.LoggingFunction
}

// ValidateFiles validates a collection of filenames using a RuleSet
func (l FileLinter) ValidateFiles(ruleSet assertion.RuleSet, tags []string, ruleIDs []string, loader ResourceLoader) (assertion.ValidationReport, error) {

	report := assertion.ValidationReport{
		FilesScanned:     []string{},
		ResourcesScanned: []assertion.ScannedResource{},
		Violations:       []assertion.Violation{},
	}
	rules := assertion.FilterRulesByTagAndID(ruleSet.Rules, tags, ruleIDs)
	r := ResourceLinter{Log: l.Log}
	for _, filename := range l.Filenames {
		include, err := assertion.ShouldIncludeFile(ruleSet.Files, filename)
		if err == nil && include {
			l.Log(fmt.Sprintf("Processing %s", filename))
			resources, err := loader.Load(filename)
			if err != nil {
				return report, err
			}
			r, err := r.ValidateResources(resources, rules)
			r.FilesScanned = []string{filename}
			report = combineValidationReports(report, r)
			if err != nil {
				return report, err
			}
		}
	}
	return report, nil
}

// SearchFiles evaluates a JMESPath expression against resources in a collection of filenames
func (l FileLinter) SearchFiles(ruleSet assertion.RuleSet, searchExpression string, loader ResourceLoader) {
	for _, filename := range l.Filenames {
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
