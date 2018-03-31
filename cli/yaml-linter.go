package main

import (
	"github.com/stelligent/config-lint/assertion"
)

// YAMLLinter lints rules from a generic YAML file
type YAMLLinter struct {
	Log assertion.LoggingFunction
}

// YAMLResourceLoader loads a list of Resource objects based on the list of ResourceConfig objects
type YAMLResourceLoader struct {
	Log       assertion.LoggingFunction
	Resources []assertion.ResourceConfig
}

// Load converts a text file into a collection of Resource objects
func (l YAMLResourceLoader) Load(filename string) ([]assertion.Resource, error) {
	resources := make([]assertion.Resource, 0)
	yamlResources, err := loadYAML(filename, l.Log)
	if err != nil {
		return resources, err
	}
	for _, document := range yamlResources {
		for _, resourceConfig := range l.Resources {
			matches, err := assertion.SearchData(resourceConfig.Key, document)
			if err != nil {
				return resources, nil
			}
			sliceOfProperties := matches.([]interface{})
			for _, element := range sliceOfProperties {
				properties := element.(map[string]interface{})
				resource := assertion.Resource{
					ID:         properties[resourceConfig.ID].(string), // FIXME should use JMESPath to ID, and make sure it handles numbers
					Type:       resourceConfig.Type,
					Properties: properties,
					Filename:   filename,
				}
				resources = append(resources, resource)
			}
		}
	}
	return resources, nil
}

// Validate runs validate on a collection of filenames using a RuleSet
func (l YAMLLinter) Validate(filenames []string, ruleSet assertion.RuleSet, tags []string, ruleIDs []string) ([]string, []assertion.Violation, error) {
	loader := YAMLResourceLoader{Log: l.Log, Resources: ruleSet.Resources}
	f := FileLinter{Log: l.Log}
	return f.ValidateFiles(filenames, ruleSet, tags, ruleIDs, loader)
}

// Search evaluates a JMESPath expression against the resources in a collection of filenames
func (l YAMLLinter) Search(filenames []string, ruleSet assertion.RuleSet, searchExpression string) {
	loader := YAMLResourceLoader{Log: l.Log, Resources: ruleSet.Resources}
	f := FileLinter{Log: l.Log}
	f.SearchFiles(filenames, ruleSet, searchExpression, loader)
}
