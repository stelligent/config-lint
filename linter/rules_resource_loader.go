package linter

import (
	"github.com/stelligent/config-lint/assertion"
)

// RulesResourceLoader converts a YAML configuration file into a collection with Resource objects
type RulesResourceLoader struct{}

func getAttr(m map[string]interface{}, keys ...string) []interface{} {
	for _, key := range keys {
		if r, ok := m[key].([]interface{}); ok {
			return r
		}
	}
	return []interface{}{}
}

// Load converts a text file into a collection of Resource objects
func (l RulesResourceLoader) Load(filename string) (FileResources, error) {
	loaded := FileResources{
		Resources: make([]assertion.Resource, 0),
	}
	yamlResources, err := loadYAML(filename)
	if err != nil {
		return loaded, err
	}
	for _, ruleSet := range yamlResources {
		setResource := assertion.Resource{
			ID:         getResourceIDFromFilename(filename),
			Type:       "LintRuleSet",
			Properties: ruleSet,
			Filename:   filename,
		}
		loaded.Resources = append(loaded.Resources, setResource)
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
			loaded.Resources = append(loaded.Resources, ruleResource)
		}
	}
	return loaded, nil
}

func (l RulesResourceLoader) ReplaceVariables(resources []assertion.Resource, variables []Variable) ([]assertion.Resource, error) {
	return resources, nil
}
