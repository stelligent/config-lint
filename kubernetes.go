package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

type KubernetesLinter struct {
	Log LoggingFunction
}

// TODO - is it really necessary to have two types?
type KubernetesResource = TerraformResource

// TODO duplicates loadTerraformRules
func loadKubernetesRules(filename string) string {
	kubernetesRules, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(kubernetesRules)
}

func loadYAML(filename string, log LoggingFunction) []interface{} {
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

func loadKubernetesResources(filename string, log LoggingFunction) []KubernetesResource {
	yamlResources := loadYAML(filename, log)
	resources := make([]KubernetesResource, 0)
	for _, resource := range yamlResources {
		m := resource.(map[string]interface{})
		kr := KubernetesResource{
			Id:         filename,
			Type:       m["kind"].(string),
			Properties: m,
			Filename:   filename,
		}
		resources = append(resources, kr)
	}
	return resources
}

func filterKubernetesResourcesByType(resources []KubernetesResource, resourceType string) []KubernetesResource {
	if resourceType == "*" {
		return resources
	}
	filtered := make([]KubernetesResource, 0)
	for _, resource := range resources {
		if resource.Type == resourceType {
			filtered = append(filtered, resource)
		}
	}
	return filtered
}

func validateKubernetesResources(report *ValidationReport, resources []KubernetesResource, rules []Rule, tags []string, log LoggingFunction) {
	for _, rule := range filterRulesByTag(rules, tags) {
		log(fmt.Sprintf("Rule %s: %s", rule.Id, rule.Message))
		for _, filter := range rule.Filters {
			for _, resource := range filterKubernetesResourcesByType(resources, rule.Resource) {
				log(fmt.Sprintf("Checking resource %s", resource.Id))
				status := applyFilter(rule, filter, resource, log)
				if status != "OK" {
					v := Violation{
						RuleId:       rule.Id,
						ResourceId:   resource.Id,
						ResourceType: resource.Type,
						Status:       status,
						Message:      rule.Message,
						Filename:     resource.Filename,
					}
					report.Violations[status] = append(report.Violations[status], v)
				}
			}
		}
	}
}

func (l KubernetesLinter) Validate(filenames []string, ruleSet RuleSet, tags []string, ruleIds []string) ValidationReport {
	report := ValidationReport{
		Violations:   make(map[string]([]Violation), 0),
		FilesScanned: make([]string, 0),
	}
	rules := filterRulesById(ruleSet.Rules, ruleIds)
	for _, filename := range filenames {
		if shouldIncludeFile(ruleSet.Files, filename) {
			l.Log(fmt.Sprintf("Processing %s", filename))
			resources := loadKubernetesResources(filename, l.Log)
			validateKubernetesResources(&report, resources, rules, tags, l.Log)
			report.FilesScanned = append(report.FilesScanned, filename)
		}
	}
	return report
}

func (l KubernetesLinter) Search(filenames []string, searchExpression string) {
	for _, filename := range filenames {
		l.Log(fmt.Sprintf("Searching %s", filename))
		resources := loadKubernetesResources(filename, l.Log)
		for _, resource := range resources {
			v, err := searchData(searchExpression, resource.Properties)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%s: %s\n", filename, v)
			}
		}
	}
}
