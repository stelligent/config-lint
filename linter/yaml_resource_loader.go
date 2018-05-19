package linter

import (
	"github.com/stelligent/config-lint/assertion"
)

// YAMLResourceLoader loads a list of Resource objects based on the list of ResourceConfig objects
type YAMLResourceLoader struct {
	Resources []assertion.ResourceConfig
}

func extractResourceID(expression string, properties interface{}) string {
	resourceID := "None"
	result, err := assertion.SearchData(expression, properties)
	if err == nil {
		resourceID, _ = result.(string)
	}
	return resourceID
}

// Load converts a text file into a collection of Resource objects
func (l YAMLResourceLoader) Load(filename string) (FileResources, error) {
	loaded := FileResources{
		Resources: make([]assertion.Resource, 0),
	}
	yamlResources, err := loadYAML(filename)
	if err != nil {
		return loaded, err
	}
	for _, document := range yamlResources {
		for _, resourceConfig := range l.Resources {
			matches, err := assertion.SearchData(resourceConfig.Key, document)
			if err != nil {
				return loaded, nil
			}
			sliceOfProperties, ok := matches.([]interface{})
			if ok {
				for _, element := range sliceOfProperties {
					properties := element.(map[string]interface{})
					resource := assertion.Resource{
						ID:         extractResourceID(resourceConfig.ID, properties),
						Type:       resourceConfig.Type,
						Properties: properties,
						Filename:   filename,
					}
					loaded.Resources = append(loaded.Resources, resource)
				}
			}
		}
	}
	return loaded, nil
}

func (l YAMLResourceLoader) PostLoad(r FileResources) ([]assertion.Resource, error) {
	return r.Resources, nil
}
