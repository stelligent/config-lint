package main

import (
	"github.com/stelligent/config-lint/assertion"
)

// RulesLinter lints rules files for itself
type RulesLinter struct {
	Filenames   []string
	Log         assertion.LoggingFunction
	ValueSource assertion.ValueSource
}

// RulesResourceLoader converts a YAML configuration file into a collection with Resource objects
type RulesResourceLoader struct {
	Log assertion.LoggingFunction
}

func getAttr(m map[string]interface{}, keys ...string) []interface{} {
	for _, key := range keys {
		if r, ok := m[key].([]interface{}); ok {
			return r
		}
	}
	return []interface{}{}
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
		resources = append(resources, setResource)
		// The LintRuleSet resources already has an attribute called Rules
		// but also adding each of these rules as a separate LintRule resource
		// makes writing rules a lot simpler
		m := ruleSet.(map[string]interface{})
		rules := getAttr(m, "rules", "Rules")
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
	return resources, nil
}

// Validate runs validate on a collection of filenames using a RuleSet
func (l RulesLinter) Validate(ruleSet assertion.RuleSet, tags []string, ruleIDs []string) (assertion.ValidationReport, error) {
	loader := RulesResourceLoader{Log: l.Log}
	f := FileLinter{Filenames: l.Filenames, Log: l.Log, Loader: loader}
	return f.ValidateFiles(ruleSet, tags, ruleIDs)
}

// Search evaluates a JMESPath expression against the resources in a collection of filenames
func (l RulesLinter) Search(ruleSet assertion.RuleSet, searchExpression string) {
	loader := RulesResourceLoader{Log: l.Log}
	f := FileLinter{Filenames: l.Filenames, Log: l.Log, Loader: loader}
	f.SearchFiles(ruleSet, searchExpression)
}
