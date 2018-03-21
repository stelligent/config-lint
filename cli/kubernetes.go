package main

import (
	"github.com/ghodss/yaml"
	"github.com/lhitchon/config-lint/assertion"
	"io/ioutil"
	"path/filepath"
)

type KubernetesLinter struct {
	BaseLinter
	Log assertion.LoggingFunction
}

type KubernetesResourceLoader struct {
	Log assertion.LoggingFunction
}

func loadYAML(filename string, log assertion.LoggingFunction) []interface{} {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var yamlData interface{}
	err = yaml.Unmarshal(content, &yamlData)
	if err != nil {
		panic(err)
	}
	m := yamlData.(map[string]interface{})
	return []interface{}{m}
}

func getResourceIdFromMetadata(m map[string]interface{}) (string, bool) {
	if metadata, ok := m["metadata"].(map[string]interface{}); ok {
		if name, ok := metadata["name"].(string); ok {
			return name, true
		}
	}
	return "", false
}

func getResourceIdFromFilename(filename string) string {
	_, resourceId := filepath.Split(filename)
	return resourceId
}

func (l KubernetesResourceLoader) Load(filename string) []assertion.Resource {
	yamlResources := loadYAML(filename, l.Log)
	resources := make([]assertion.Resource, 0)
	for _, resource := range yamlResources {
		m := resource.(map[string]interface{})
		var resourceId string
		if name, ok := getResourceIdFromMetadata(m); ok {
			resourceId = name
		} else {
			resourceId = getResourceIdFromFilename(filename)
		}
		kr := assertion.Resource{
			Id:         resourceId,
			Type:       m["kind"].(string),
			Properties: m,
			Filename:   filename,
		}
		resources = append(resources, kr)
	}
	return resources
}

func (l KubernetesLinter) Validate(filenames []string, ruleSet assertion.RuleSet, tags []string, ruleIds []string) ([]string, []assertion.Violation) {
	loader := KubernetesResourceLoader{Log: l.Log}
	return l.ValidateFiles(filenames, ruleSet, tags, ruleIds, loader, l.Log)
}

func (l KubernetesLinter) Search(filenames []string, ruleSet assertion.RuleSet, searchExpression string) {
	loader := KubernetesResourceLoader{Log: l.Log}
	l.SearchFiles(filenames, ruleSet, searchExpression, loader)
}
