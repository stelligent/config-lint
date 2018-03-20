package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/lhitchon/config-lint/assertion"
	"io/ioutil"
)

type KubernetesLinter struct {
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

func loadKubernetesResources(filename string, log assertion.LoggingFunction) []assertion.Resource {
	yamlResources := loadYAML(filename, log)
	resources := make([]assertion.Resource, 0)
	for _, resource := range yamlResources {
		m := resource.(map[string]interface{})
		kr := assertion.Resource{
			Id:         filename,
			Type:       m["kind"].(string),
			Properties: m,
			Filename:   filename,
		}
		resources = append(resources, kr)
	}
	return resources
}

func (l KubernetesLinter) ValidateKubernetesResources(resources []assertion.Resource, rules []assertion.Rule, tags []string) []assertion.Violation {

	valueSource := assertion.StandardValueSource{Log: l.Log}
	filteredRules := assertion.FilterRulesByTag(rules, tags)
	resolvedRules := assertion.ResolveRules(filteredRules, valueSource, l.Log)
	externalRules := assertion.StandardExternalRuleInvoker{Log: l.Log}

	allViolations := make([]assertion.Violation, 0)
	for _, rule := range resolvedRules {
		l.Log(fmt.Sprintf("Rule %s: %s", rule.Id, rule.Message))
		for _, resource := range assertion.FilterResourcesByType(resources, rule.Resource) {
			if assertion.ExcludeResource(rule, resource) {
				l.Log(fmt.Sprintf("Ignoring resource %s", resource.Id))
			} else {
				_, violations := assertion.CheckRule(rule, resource, externalRules, l.Log)
				allViolations = append(allViolations, violations...)
			}
		}
	}
	return allViolations
}

func (l KubernetesLinter) Validate(filenames []string, ruleSet assertion.RuleSet, tags []string, ruleIds []string) ([]string, []assertion.Violation) {
	rules := assertion.FilterRulesById(ruleSet.Rules, ruleIds)
	allViolations := make([]assertion.Violation, 0)
	filesScanned := make([]string, 0)
	for _, filename := range filenames {
		if assertion.ShouldIncludeFile(ruleSet.Files, filename) {
			l.Log(fmt.Sprintf("Processing %s", filename))
			resources := loadKubernetesResources(filename, l.Log)
			violations := l.ValidateKubernetesResources(resources, rules, tags)
			allViolations = append(allViolations, violations...)
			filesScanned = append(filesScanned, filename)
		}
	}
	return filesScanned, allViolations
}

func (l KubernetesLinter) Search(filenames []string, ruleSet assertion.RuleSet, searchExpression string) {
	for _, filename := range filenames {
		if assertion.ShouldIncludeFile(ruleSet.Files, filename) {
			fmt.Printf("Searching %s:\n", filename)
			resources := loadKubernetesResources(filename, l.Log)
			for _, resource := range resources {
				v, err := assertion.SearchData(searchExpression, resource.Properties)
				if err != nil {
					fmt.Println(err)
				} else {
					s, err := assertion.JSONStringify(v)
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Printf("%s: %s\n", resource.Id, s)
					}
				}
			}
		}
	}
}
