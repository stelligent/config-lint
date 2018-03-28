package main

import (
	"github.com/stelligent/config-lint/assertion"
)

// RulesLinter lints rules files for itself
type RulesLinter struct {
	BaseLinter
	Log assertion.LoggingFunction
}

// RulesResourceLoader converts a YAML configuration file into a collection with Resource objects
type RulesResourceLoader struct {
	Log assertion.LoggingFunction
}

// Load converts a text file into a collection of Resource objects
func (l RulesResourceLoader) Load(filename string) ([]assertion.Resource, error) {
	resources := make([]assertion.Resource, 0)
	yamlResources, err := loadYAML(filename, l.Log)
	if err != nil {
		return resources, err
	}
	for _, ruleSet := range yamlResources {
		setResource := assertion.Resource{
			ID:         getResourceIDFromFilename(filename),
			Type:       "LintRuleSet",
			Properties: ruleSet,
			Filename:   filename,
		}
		// The LintRuleSet resources already has an attribute calls Rules
		// but also adding each of these rules as a separate LintRule resource
		// makes writing rules a lot simpler
		resources = append(resources, setResource)
		for _, resource := range yamlResources {
			m := resource.(map[string]interface{})
			rules := m["Rules"].([]interface{})
			for _, rule := range rules {
				properties := rule.(map[string]interface{})
				ruleResource := assertion.Resource{
					ID:         properties["id"].(string),
					Type:       "LintRule",
					Properties: properties,
					Filename:   filename,
				}
				resources = append(resources, ruleResource)
			}
		}
	}
	return resources, nil
}

// Validate runs validate on a collection of filenames using a RuleSet
func (l RulesLinter) Validate(filenames []string, ruleSet assertion.RuleSet, tags []string, ruleIDs []string) ([]string, []assertion.Violation, error) {
	loader := RulesResourceLoader{Log: l.Log}
	return l.ValidateFiles(filenames, ruleSet, tags, ruleIDs, loader, l.Log)
}

// Search evaluates a JMESPath expression against the resources in a collection of filenames
func (l RulesLinter) Search(filenames []string, ruleSet assertion.RuleSet, searchExpression string) {
	loader := RulesResourceLoader{Log: l.Log}
	l.SearchFiles(filenames, ruleSet, searchExpression, loader)
}
