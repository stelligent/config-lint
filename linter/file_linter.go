package linter

import (
	"fmt"
	"github.com/stelligent/config-lint/assertion"
	"io"
)

type (
	Variable struct {
		Name  string
		Value interface{}
	}
	FileResources struct {
		Resources []assertion.Resource
		Variables []Variable
	}
	// FileResourceLoader provides the interface that a Linter needs to load a collection of Resource objects
	FileResourceLoader interface {
		Load(filename string) (FileResources, error)
		ReplaceVariables(resources []assertion.Resource, variables []Variable) ([]assertion.Resource, error)
	}
)

// FileLinter provides implementation for some common functions that are used by multiple Linter implementations
type FileLinter struct {
	Filenames   []string
	ValueSource assertion.ValueSource
	Loader      FileResourceLoader
}

// Validate validates a collection of filenames using a RuleSet
func (fl FileLinter) Validate(ruleSet assertion.RuleSet, options Options) (assertion.ValidationReport, error) {

	rules := assertion.FilterRulesByTagAndID(ruleSet.Rules, options.Tags, options.RuleIDs)
	rl := ResourceLinter{ValueSource: fl.ValueSource}

	resources := []assertion.Resource{}
	variables := []Variable{}
	filesScanned := []string{}

	for _, filename := range fl.Filenames {
		include, err := assertion.ShouldIncludeFile(ruleSet.Files, filename)
		if err == nil && include {
			assertion.Debugf("Processing %s\n", filename)
			loaded, err := fl.Loader.Load(filename)
			if err != nil {
				return assertion.ValidationReport{}, err
			}
			assertion.Debugf("Found variables %v\n", loaded.Variables)
			filesScanned = append(filesScanned, filename)
			resources = append(resources, loaded.Resources...)
			variables = append(variables, loaded.Variables...)
		}
	}
	resolvedResources, err := fl.Loader.ReplaceVariables(resources, variables)
	if err != nil {
		return assertion.ValidationReport{}, err
	}
	report, err := rl.ValidateResources(resolvedResources, rules)
	if err != nil {
		return report, err
	}
	report.FilesScanned = filesScanned
	return report, nil
}

// Search evaluates a JMESPath expression against resources in a collection of filenames
func (fl FileLinter) Search(ruleSet assertion.RuleSet, searchExpression string, w io.Writer) {
	for _, filename := range fl.Filenames {
		include, _ := assertion.ShouldIncludeFile(ruleSet.Files, filename) // FIXME what about error?
		if include {
			fmt.Fprintf(w, "Searching %s:\n", filename)
			loaded, err := fl.Loader.Load(filename)
			if err != nil {
				fmt.Fprintln(w, "Error for file:", filename)
				fmt.Fprintln(w, err.Error())
			}
			for _, resource := range loaded.Resources {
				v, err := assertion.SearchData(searchExpression, resource.Properties)
				if err != nil {
					fmt.Fprintln(w, err)
				} else {
					s, err := assertion.JSONStringify(v)
					if err != nil {
						fmt.Fprintln(w, err)
					} else {
						fmt.Fprintf(w, "%s (%s): %s\n", resource.ID, resource.Type, s)
					}
				}
			}
		}
	}
}
