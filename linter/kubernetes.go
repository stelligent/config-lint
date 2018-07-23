package linter

import (
	"github.com/stelligent/config-lint/assertion"
	"path/filepath"
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
		properties := resource.(map[string]interface{})
		var resourceID string
		if name, ok := getResourceIDFromMetadata(properties); ok {
			resourceID = name
		} else {
			resourceID = getResourceIDFromFilename(filename)
		}
		properties["__file__"] = filename
		properties["__dir__"] = filepath.Dir(filename)
		kr := assertion.Resource{
			ID:         resourceID,
			Type:       properties["kind"].(string),
			Properties: properties,
			Filename:   filename,
		}
		loaded.Resources = append(loaded.Resources, kr)
	}
	return loaded, nil
}

// PostLoad does no additional processing for a KubernetesLoader
func (l KubernetesResourceLoader) PostLoad(r FileResources) ([]assertion.Resource, error) {
	return r.Resources, nil
}
