package linter

import (
	"github.com/stelligent/config-lint/assertion"
)

// KubernetesLinter lints resources in Kubernets YAML files
type KubernetesLinter struct {
	Filenames   []string
	ValueSource assertion.ValueSource
}

// KubernetesResourceLoader converts Kubernetes configuration files into a collection of Resource objects
type KubernetesResourceLoader struct{}

func getResourceIDFromMetadata(m map[string]interface{}) (string, bool) {
	if metadata, ok := m["metadata"].(map[string]interface{}); ok {
		if name, ok := metadata["name"].(string); ok {
			return name, true
		}
	}
	return "", false
}

// Load converts a text file into a collection of Resource objects
func (l KubernetesResourceLoader) Load(filename string) (FileResources, error) {
	loaded := FileResources{
		Resources: make([]assertion.Resource, 0),
	}
	yamlResources, err := loadYAML(filename)
	if err != nil {
		return loaded, err
	}
	for _, resource := range yamlResources {
		m := resource.(map[string]interface{})
		var resourceID string
		if name, ok := getResourceIDFromMetadata(m); ok {
			resourceID = name
		} else {
			resourceID = getResourceIDFromFilename(filename)
		}
		kr := assertion.Resource{
			ID:         resourceID,
			Type:       m["kind"].(string),
			Properties: m,
			Filename:   filename,
		}
		loaded.Resources = append(loaded.Resources, kr)
	}
	return loaded, nil
}

func (l KubernetesResourceLoader) ReplaceVariables(resources []assertion.Resource, variables []Variable) ([]assertion.Resource, error) {
	return resources, nil
}
