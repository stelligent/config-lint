package main

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/hcl"
	"io/ioutil"
)

type TerraformResource struct {
	Id         string
	Type       string
	Properties interface{}
	Filename   string
}

func loadHCL(filename string, log LoggingFunction) []interface{} {
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

func loadTerraformResources(filename string, log LoggingFunction) []TerraformResource {
	hclResources := loadHCL(filename, log)

	resources := make([]TerraformResource, 0)
	for _, resource := range hclResources {
		for resourceType, templateResources := range resource.(map[string]interface{}) {
			if templateResources != nil {
				for _, templateResource := range templateResources.([]interface{}) {
					for resourceId, resource := range templateResource.(map[string]interface{}) {
						tr := TerraformResource{
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

func loadTerraformRules(filename string) string {
	terraformRules, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(terraformRules)
}

func filterTerraformResourcesByType(resources []TerraformResource, resourceType string) []TerraformResource {
	if resourceType == "*" {
		return resources
	}
	filtered := make([]TerraformResource, 0)
	for _, resource := range resources {
		if resource.Type == resourceType {
			filtered = append(filtered, resource)
		}
	}
	return filtered
}

func validateTerraformResources(report *ValidationReport, resources []TerraformResource, rules []Rule, tags []string, log LoggingFunction) {
	for _, rule := range filterRulesByTag(rules, tags) {
		log(fmt.Sprintf("Rule %s: %s", rule.Id, rule.Message))
		for _, filter := range rule.Filters {
			for _, resource := range filterTerraformResourcesByType(resources, rule.Resource) {
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
					report.AllViolations = append(report.AllViolations, v)
					if status == "WARNING" {
						report.Warnings = append(report.Warnings, v)
					}
					if status == "FAILURE" {
						report.Failures = append(report.Failures, v)
					}
				}
			}
		}
	}
}

func terraform(filenames []string, rulesFilename string, tags []string, ruleIds []string, log LoggingFunction) ValidationReport {
	var report ValidationReport
	ruleSet := MustParseRules(loadTerraformRules(rulesFilename))
	rules := filterRulesById(ruleSet.Rules, ruleIds)
	for _, filename := range filenames {
		if shouldIncludeFile(ruleSet.Files, filename) {
			resources := loadTerraformResources(filename, log)
			validateTerraformResources(&report, resources, rules, tags, log)
			report.FilesScanned = append(report.FilesScanned, filename)
		}
	}
	return report
}

func terraformSearch(filenames []string, searchExpression string, log LoggingFunction) {
	for _, filename := range filenames {
		log(fmt.Sprintf("Searching %s", filename))
		resources := loadTerraformResources(filename, log)
		for _, resource := range resources {
			v := searchData(searchExpression, resource.Properties)
			if v != "null" {
				fmt.Printf("%s: %s\n", filename, v)
			}
		}
	}
}
