package main

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/hcl"
	"github.com/lhitchon/config-lint/assertion"
	"io/ioutil"
)

type TerraformLinter struct {
	Log assertion.LoggingFunction
}

func loadHCL(filename string, log assertion.LoggingFunction) []interface{} {
	template, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var v interface{}
	err = hcl.Unmarshal([]byte(template), &v)
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	log(string(jsonData))

	var hclData interface{}
	err = yaml.Unmarshal(jsonData, &hclData)
	if err != nil {
		panic(err)
	}
	m := hclData.(map[string]interface{})
	results := make([]interface{}, 0)
	for _, key := range []string{"resource", "data"} {
		if m[key] != nil {
			log(fmt.Sprintf("Adding %s", key))
			results = append(results, m[key].([]interface{})...)
		}
	}
	return results
}

func loadTerraformResources(filename string, log assertion.LoggingFunction) []assertion.Resource {
	hclResources := loadHCL(filename, log)

	resources := make([]assertion.Resource, 0)
	for _, resource := range hclResources {
		for resourceType, templateResources := range resource.(map[string]interface{}) {
			if templateResources != nil {
				for _, templateResource := range templateResources.([]interface{}) {
					for resourceId, resource := range templateResource.(map[string]interface{}) {
						tr := assertion.Resource{
							Id:         resourceId,
							Type:       resourceType,
							Properties: resource.([]interface{})[0],
							Filename:   filename,
						}
						resources = append(resources, tr)
					}
				}
			}
		}
	}
	return resources
}

func filterTerraformResourcesByType(resources []assertion.Resource, resourceType string) []assertion.Resource {
	if resourceType == "*" {
		return resources
	}
	filtered := make([]assertion.Resource, 0)
	for _, resource := range resources {
		if resource.Type == resourceType {
			filtered = append(filtered, resource)
		}
	}
	return filtered
}

func (l TerraformLinter) ValidateTerraformResources(report *assertion.ValidationReport, resources []assertion.Resource, rules []assertion.Rule, tags []string) {

	valueSource := assertion.StandardValueSource{Log: l.Log}
	filteredRules := assertion.FilterRulesByTag(rules, tags)
	resolvedRules := assertion.ResolveRules(filteredRules, valueSource, l.Log)

	for _, rule := range resolvedRules {
		l.Log(fmt.Sprintf("Rule %s: %s", rule.Id, rule.Message))
		for _, resource := range filterTerraformResourcesByType(resources, rule.Resource) {
			if assertion.ExcludeResource(rule, resource) {
				l.Log(fmt.Sprintf("Ignoring resource %s", resource.Id))
			} else {
				_, violations := assertion.CheckRule(rule, resource, l.Log)
				for _, violation := range violations {
					report.Violations[violation.Status] = append(report.Violations[violation.Status], violation)
				}
			}
		}
	}
}

func (l TerraformLinter) Validate(filenames []string, ruleSet assertion.RuleSet, tags []string, ruleIds []string) assertion.ValidationReport {
	report := assertion.ValidationReport{
		Violations:   make(map[string]([]assertion.Violation), 0),
		FilesScanned: make([]string, 0),
	}
	rules := assertion.FilterRulesById(ruleSet.Rules, ruleIds)
	for _, filename := range filenames {
		if assertion.ShouldIncludeFile(ruleSet.Files, filename) {
			resources := loadTerraformResources(filename, l.Log)
			l.ValidateTerraformResources(&report, resources, rules, tags)
			report.FilesScanned = append(report.FilesScanned, filename)
		}
	}
	return report
}

func (l TerraformLinter) Search(filenames []string, searchExpression string) {
	for _, filename := range filenames {
		l.Log(fmt.Sprintf("Searching %s", filename))
		resources := loadTerraformResources(filename, l.Log)
		for _, resource := range resources {
			v, err := assertion.SearchData(searchExpression, resource.Properties)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%s: %s\n", filename, v)
			}
		}
	}
}
