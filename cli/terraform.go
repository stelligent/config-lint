package main

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/hcl"
	"github.com/stelligent/config-lint/assertion"
	"io/ioutil"
)

// TerraformLinter implements a Linter for Terraform configuration files
type TerraformLinter struct {
	BaseLinter
	Log assertion.LoggingFunction
}

// TerraformResourceLoader converts Terraform configuration files into JSON objects
type TerraformResourceLoader struct {
	Log assertion.LoggingFunction
}

func parsePolicy(resource assertion.Resource) assertion.Resource {
	if resource.Properties != nil {
		properties := resource.Properties.(map[string]interface{})
		for _, attribute := range []string{"assume_role_policy", "policy"} {
			if policyAttribute, hasPolicyString := properties[attribute]; hasPolicyString {
				if policyString, isString := policyAttribute.(string); isString {
					var policy interface{}
					err := json.Unmarshal([]byte(policyString), &policy)
					if err != nil {
						panic(err)
					}
					properties[attribute] = policy
				}
			}
		}
	}
	return resource
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

// Load parses an HCL file into a collection or Resource objects
func (l TerraformResourceLoader) Load(filename string) []assertion.Resource {
	hclResources := loadHCL(filename, l.Log)

	resources := make([]assertion.Resource, 0)
	for _, resource := range hclResources {
		for resourceType, templateResources := range resource.(map[string]interface{}) {
			if templateResources != nil {
				for _, templateResource := range templateResources.([]interface{}) {
					for resourceID, resource := range templateResource.(map[string]interface{}) {
						tr := assertion.Resource{
							ID:         resourceID,
							Type:       resourceType,
							Properties: resource.([]interface{})[0],
							Filename:   filename,
						}
						resources = append(resources, parsePolicy(tr))
					}
				}
			}
		}
	}
	return resources
}

// Validate uses a RuleSet to validate resources in a collection of Terraform configuration files
func (l TerraformLinter) Validate(filenames []string, ruleSet assertion.RuleSet, tags []string, ruleIDs []string) ([]string, []assertion.Violation) {
	loader := TerraformResourceLoader{Log: l.Log}
	return l.ValidateFiles(filenames, ruleSet, tags, ruleIDs, loader, l.Log)
}

// Search applies a JMESPath expression to the resources in a collection of Terraform configuration files
func (l TerraformLinter) Search(filenames []string, ruleSet assertion.RuleSet, searchExpression string) {
	loader := TerraformResourceLoader{Log: l.Log}
	l.SearchFiles(filenames, ruleSet, searchExpression, loader)
}
