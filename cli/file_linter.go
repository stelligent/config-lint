package main

import (
	"fmt"
	"github.com/stelligent/config-lint/assertion"
)

// ResourceLoader provides the interface that a Linter needs to load a collection of Resource objects
type FileResourceLoader interface {
	Load(filename string) ([]assertion.Resource, error)
}

// FileLinter provides implementation for some common functions that are used by multiple Linter implementations
type FileLinter struct {
	Filenames   []string
	Log         assertion.LoggingFunction
	ValueSource assertion.ValueSource
	Loader      FileResourceLoader
}

// ValidateFiles validates a collection of filenames using a RuleSet
func (fl FileLinter) ValidateFiles(ruleSet assertion.RuleSet, options LinterOptions) (assertion.ValidationReport, error) {

	report := assertion.ValidationReport{
		FilesScanned:     []string{},
		ResourcesScanned: []assertion.ScannedResource{},
		Violations:       []assertion.Violation{},
	}
	rules := assertion.FilterRulesByTagAndID(ruleSet.Rules, options.Tags, options.RuleIDs)
	rl := ResourceLinter{Log: fl.Log, ValueSource: fl.ValueSource}
	for _, filename := range fl.Filenames {
		include, err := assertion.ShouldIncludeFile(ruleSet.Files, filename)
		if err == nil && include {
			fl.Log(fmt.Sprintf("Processing %s", filename))
			resources, err := fl.Loader.Load(filename)
			if err != nil {
				return report, err
			}
			r, err := rl.ValidateResources(resources, rules)
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
func (l FileLinter) SearchFiles(ruleSet assertion.RuleSet, searchExpression string) {
	for _, filename := range l.Filenames {
		include, _ := assertion.ShouldIncludeFile(ruleSet.Files, filename) // FIXME what about error?
		if include {
			fmt.Printf("Searching %s:\n", filename)
			resources, err := l.Loader.Load(filename)
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
