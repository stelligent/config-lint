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

func loadHCL(template string, log LoggingFunction) []interface{} {
	var v interface{}
	err := hcl.Unmarshal([]byte(template), &v)
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(v, "", "  ")
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

func loadTerraformResources(filename string, hclResources []interface{}) []TerraformResource {
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

func loadTerraformRules() string {
	terraformRules, err := ioutil.ReadFile("./rules/terraform.yml")
	if err != nil {
		panic(err)
	}
	return string(terraformRules)
}

func filterTerraformResourcesByType(resources []TerraformResource, resourceType string) []TerraformResource {
	filtered := make([]TerraformResource, 0)
	for _, resource := range resources {
		if resource.Type == resourceType {
			filtered = append(filtered, resource)
		}
	}
	return filtered
}

func validateTerraformResources(resources []TerraformResource, rules []Rule, tags []string, log LoggingFunction) []ValidationResult {
	results := make([]ValidationResult, 0)
	for _, rule := range filterRulesByTag(rules, tags) {
		log(fmt.Sprintf("Rule %s: %s", rule.Id, rule.Message))
		for _, filter := range rule.Filters {
			for _, resource := range filterTerraformResourcesByType(resources, rule.Resource) {
				log(fmt.Sprintf("Checking resource %s", resource.Id))
				status := applyFilter(rule, filter, resource, log)
				if status != "OK" {
					results = append(results, ValidationResult{
						RuleId:       rule.Id,
						ResourceId:   resource.Id,
						ResourceType: resource.Type,
						Status:       status,
						Message:      rule.Message,
						Filename:     resource.Filename,
					})
				}
			}
		}
	}
	return results
}
