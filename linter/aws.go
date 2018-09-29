package linter

import (
	"fmt"
	"github.com/stelligent/config-lint/assertion"
	"io"
)

type (
	// AWSResourceLoader provides the interface that a Linter needs to load a collection or Resource objects
	AWSResourceLoader interface {
		Load(ruleSet assertion.RuleSet) ([]assertion.Resource, error)
	}
	AWSAPIResourceLoader interface {
		Load() ([]assertion.Resource, error)
	}

	// AWSLazyLoader provides the
	AWSLazyLoader struct{}

	// AWSResourceLinter implements a Linter for data returned by the calls to the AWS SDK
	AWSResourceLinter struct {
		RuleSet     assertion.RuleSet
		Loader      AWSResourceLoader
		ValueSource assertion.ValueSource
	}
)

// Load uses the AWS API to load resources
func (l AWSLazyLoader) Load(ruleSet assertion.RuleSet) ([]assertion.Resource, error) {
	types := map[string]bool{}
	for _, r := range ruleSet.Rules {
		if r.Resource != "" {
			types[r.Resource] = true
		}
		if len(r.Resources) > 0 {
			for _, t := range r.Resources {
				types[t] = true
			}
		}
	}
	resources := []assertion.Resource{}
	loaders := map[string]AWSAPIResourceLoader{
		"AWS::IAM::User":          IAMUserLoader{},
		"AWS::IAM::Role":          IAMRoleLoader{},
		"AWS::IAM::Group":         IAMGroupLoader{},
		"AWS::EC2::SecurityGroup": SecurityGroupLoader{},
	}
	for t, _ := range types {
		if loader, ok := loaders[t]; ok {
			fmt.Println("Load:", t)
			r, err := loader.Load()
			if err == nil {
				resources = append(resources, r...)
			}
		} else {
			fmt.Println("Load:", t, "not implemented")
		}
	}
	return resources, nil
}

// Validate applies a Ruleset to all SecurityGroups
func (l AWSResourceLinter) Validate(ruleSet assertion.RuleSet, options Options) (assertion.ValidationReport, error) {
	rules := assertion.FilterRulesByTagAndID(ruleSet.Rules, options.Tags, options.RuleIDs, options.IgnoreRuleIDs)
	rl := ResourceLinter{ValueSource: l.ValueSource}
	resources, err := l.Loader.Load(ruleSet)
	if err != nil {
		return assertion.ValidationReport{}, err
	}
	return rl.ValidateResources(resources, rules)
}

// Search applies a JMESPath to all SecurityGroups
func (l AWSResourceLinter) Search(ruleSet assertion.RuleSet, searchExpression string, w io.Writer) {
	resources, _ := l.Loader.Load(ruleSet)
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
