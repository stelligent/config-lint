package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/lhitchon/config-lint/filter"
	"io/ioutil"
)

type KubernetesLinter struct {
	Log filter.LoggingFunction
}

func loadYAML(filename string, log filter.LoggingFunction) []interface{} {
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

func loadKubernetesResources(filename string, log filter.LoggingFunction) []filter.Resource {
	yamlResources := loadYAML(filename, log)
	resources := make([]filter.Resource, 0)
	for _, resource := range yamlResources {
		m := resource.(map[string]interface{})
		kr := filter.Resource{
			Id:         filename,
			Type:       m["kind"].(string),
			Properties: m,
			Filename:   filename,
		}
		resources = append(resources, kr)
	}
	return resources
}

func filterKubernetesResourcesByType(resources []filter.Resource, resourceType string) []filter.Resource {
	if resourceType == "*" {
		return resources
	}
	filtered := make([]filter.Resource, 0)
	for _, resource := range resources {
		if resource.Type == resourceType {
			filtered = append(filtered, resource)
		}
	}
	return filtered
}

func (l KubernetesLinter) ValidateKubernetesResources(report *filter.ValidationReport, resources []filter.Resource, rules []filter.Rule, tags []string) {
	for _, rule := range filter.FilterRulesByTag(rules, tags) {
		l.Log(fmt.Sprintf("Rule %s: %s", rule.Id, rule.Message))
		for _, resource := range filterKubernetesResourcesByType(resources, rule.Resource) {
			if filter.ExcludeResource(rule, resource) {
				l.Log(fmt.Sprintf("Ignoring resource %s", resource.Id))
			} else {
				l.Log(fmt.Sprintf("Checking resource %s", resource.Id))
				_, violations := filter.ApplyRule(rule, resource, l.Log)
				for _, violation := range violations {
					report.Violations[violation.Status] = append(report.Violations[violation.Status], violation)
				}
			}
		}
	}
}

func (l KubernetesLinter) Validate(filenames []string, ruleSet filter.RuleSet, tags []string, ruleIds []string) filter.ValidationReport {
	report := filter.ValidationReport{
		Violations:   make(map[string]([]filter.Violation), 0),
		FilesScanned: make([]string, 0),
	}
	rules := filter.FilterRulesById(ruleSet.Rules, ruleIds)
	for _, filename := range filenames {
		if filter.ShouldIncludeFile(ruleSet.Files, filename) {
			l.Log(fmt.Sprintf("Processing %s", filename))
			resources := loadKubernetesResources(filename, l.Log)
			l.ValidateKubernetesResources(&report, resources, rules, tags)
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
			v, err := filter.SearchData(searchExpression, resource.Properties)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%s: %s\n", filename, v)
			}
		}
	}
}
