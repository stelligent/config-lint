package linter

import (
	"fmt"
	"github.com/stelligent/config-lint/assertion"
	"io"
)

type (
	// AWSResourceLoader uses the AWS SDK to get resource information about existing resources
	AWSResourceLoader interface {
		Load() ([]assertion.Resource, error)
	}

	// AWSResourceLinter implements a Linter for data returned by the calls to the AWS SDK
	AWSResourceLinter struct {
		Loader      AWSResourceLoader
		ValueSource assertion.ValueSource
	}
)

// Validate applies a Ruleset to all SecurityGroups
func (l AWSResourceLinter) Validate(ruleSet assertion.RuleSet, options Options) (assertion.ValidationReport, error) {
	rules := assertion.FilterRulesByTagAndID(ruleSet.Rules, options.Tags, options.RuleIDs)
	resources, err := l.Loader.Load()
	if err != nil {
		return assertion.ValidationReport{}, err
	}
	r := ResourceLinter{ValueSource: l.ValueSource}
	return r.ValidateResources(resources, rules)
}

// Search applies a JMESPath to all SecurityGroups
func (l AWSResourceLinter) Search(ruleSet assertion.RuleSet, searchExpression string, w io.Writer) {
	resources, _ := l.Loader.Load()
	for _, resource := range resources {
		v, err := assertion.SearchData(searchExpression, resource.Properties)
		if err != nil {
			fmt.Fprintln(w, err)
		} else {
			s, err := assertion.JSONStringify(v)
			if err != nil {
				fmt.Fprintln(w, err)
			} else {
				fmt.Fprintf(w, "%s: %s\n", resource.ID, s)
			}
		}
	}
}
