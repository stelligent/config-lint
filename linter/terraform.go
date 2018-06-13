package linter

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/parser"
	"github.com/stelligent/config-lint/assertion"
	"io/ioutil"
	"regexp"
)

type (
	// TerraformResourceLoader converts Terraform configuration files into JSON objects
	TerraformResourceLoader struct{}

	// TerraformLoadResult collects all the returns value for parsing an HCL string
	TerraformLoadResult struct {
		Resources    []interface{}
		Data         []interface{}
		Providers    []interface{}
		Provisioners []interface{}
		Variables    []Variable
		AST          *ast.File
	}
)

func loadHCL(filename string) (TerraformLoadResult, error) {
	result := TerraformLoadResult{
		Resources:    []interface{}{},
		Data:         []interface{}{},
		Providers:    []interface{}{},
		Provisioners: []interface{}{},
		Variables:    []Variable{},
	}
	template, err := ioutil.ReadFile(filename)
	if err != nil {
		return result, err
	}

	result.AST, _ = parser.Parse(template)

	var v interface{}
	err = hcl.Unmarshal([]byte(template), &v)
	if err != nil {
		return result, err
	}
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return result, err
	}
	assertion.Debugf("LoadHCL: %s\n", string(jsonData))

	var hclData interface{}
	err = yaml.Unmarshal(jsonData, &hclData)
	if err != nil {
		return result, err
	}
	m := hclData.(map[string]interface{})
	result.Variables = loadVariables(m["variable"])
	if m["resource"] != nil {
		result.Resources = append(result.Resources, m["resource"].([]interface{})...)
	}
	if m["data"] != nil {
		result.Data = append(result.Data, m["data"].([]interface{})...)
	}
	if m["provider"] != nil {
		result.Providers = append(result.Providers, m["provider"].([]interface{})...)
	}
	if m["provisioner"] != nil {
		result.Provisioners = append(result.Provisioners, m["provisioner"].([]interface{})...)
	}
	assertion.Debugf("LoadHCL Variables: %v\n", result.Variables)
	return result, nil
}

func loadVariables(data interface{}) []Variable {
	variables := []Variable{}
	if data == nil {
		return variables
	}
	list := data.([]interface{})
	for _, entry := range list {
		m := entry.(map[string]interface{})
		for key, value := range m {
			variables = append(variables, Variable{Name: key, Value: extractDefault(value)})
		}
	}
	return variables
}

func extractDefault(value interface{}) interface{} {
	list := value.([]interface{})
	var defaultValue interface{} = ""
	for _, entry := range list {
		m := entry.(map[string]interface{})
		defaultValue = m["default"]
	}
	return defaultValue
}

func getResourceLineNumber(resourceType, resourceID, filename string, root *ast.File) int {
	resourceItems := root.Node.(*ast.ObjectList).Filter("resource", resourceType, resourceID).Items
	if len(resourceItems) > 0 {
		resourcePos := resourceItems[0].Val.Pos()
		assertion.Debugf("Position %s %s:%d\n", resourceID, filename, resourcePos.Line)
		return resourcePos.Line
	}
	return 0
}

// Load parses an HCL file into a collection or Resource objects
func (l TerraformResourceLoader) Load(filename string) (FileResources, error) {
	loaded := FileResources{
		Resources: []assertion.Resource{},
	}
	result, err := loadHCL(filename)
	if err != nil {
		return loaded, err
	}
	loaded.Variables = result.Variables
	loaded.Resources = append(loaded.Resources, getResources(filename, result.AST, result.Resources, "resource")...)
	loaded.Resources = append(loaded.Resources, getResources(filename, result.AST, result.Data, "data")...)
	return loaded, nil
}

func getResources(filename string, ast *ast.File, objects []interface{}, category string) []assertion.Resource {
	resources := []assertion.Resource{}
	for _, resource := range objects {
		for resourceType, templateResources := range resource.(map[string]interface{}) {
			if templateResources != nil {
				for _, templateResource := range templateResources.([]interface{}) {
					for resourceID, templateResource := range templateResource.(map[string]interface{}) {
						properties := getProperties(templateResource)
						lineNumber := getResourceLineNumber(resourceType, resourceID, filename, ast)
						tr := assertion.Resource{
							ID:         resourceID,
							Type:       resourceType,
							Category:   category,
							Properties: properties,
							Filename:   filename,
							LineNumber: lineNumber,
						}
						resources = append(resources, tr)
					}
				}
			}
		}
	}
	return resources
}

func (l TerraformResourceLoader) PostLoad(fr FileResources) ([]assertion.Resource, error) {
	for _, resource := range fr.Resources {
		resource.Properties = replaceVariables(resource.Properties, fr.Variables)
	}
	for _, resource := range fr.Resources {
		properties, err := parsePolicy(resource.Properties)
		if err != nil {
			return fr.Resources, err
		}
		resource.Properties = properties
	}
	return fr.Resources, nil
}

func replaceVariables(templateResource interface{}, variables []Variable) interface{} {
	switch v := templateResource.(type) {
	case map[string]interface{}:
		return replaceVariablesInMap(v, variables)
	default:
		assertion.Debugf("replaceVariables cannot process type %T\n", v)
		return templateResource
	}
}

func replaceVariablesInMap(templateResource map[string]interface{}, variables []Variable) interface{} {
	for key, value := range templateResource {
		switch v := value.(type) {
		case string:
			templateResource[key] = resolveValue(v, variables)
		case map[string]interface{}:
			templateResource[key] = replaceVariablesInMap(v, variables)
		case []interface{}:
			templateResource[key] = replaceVariablesInList(v, variables)
		default:
			assertion.Debugf("replaceVariablesInMap cannot process type %T\n", v)
		}
	}
	return templateResource
}

func replaceVariablesInList(list []interface{}, variables []Variable) []interface{} {
	result := []interface{}{}
	for _, e := range list {
		result = append(result, replaceVariables(e, variables))
	}
	return result
}

func resolveValue(s string, variables []Variable) string {
	pattern := "[$][{]var[.](?P<name>.*)[}]"
	re, _ := regexp.Compile(pattern)
	match := re.FindStringSubmatch(s)
	if len(match) == 0 {
		return s
	}
	for _, v := range variables {
		if v.Name == match[1] {
			if replacementValue, ok := v.Value.(string); ok {
				assertion.Debugf("Replacing %s with %v\n", s, replacementValue)
				return re.ReplaceAllString(s, replacementValue)
			}
		}
	}
	defaultReplacementValue := ""
	return re.ReplaceAllString(s, defaultReplacementValue)
}

func parsePolicy(resource interface{}) (interface{}, error) {
	properties := resource.(map[string]interface{})
	for _, attribute := range []string{"assume_role_policy", "policy"} {
		if policyAttribute, hasPolicyString := properties[attribute]; hasPolicyString {
			if policyString, isString := policyAttribute.(string); isString {
				var policy interface{}
				err := json.Unmarshal([]byte(policyString), &policy)
				if err != nil {
					return properties, err
				}
				properties[attribute] = policy
			}
		}
	}
	return properties, nil
}

func getProperties(templateResource interface{}) map[string]interface{} {
	first := templateResource.([]interface{})[0] // FIXME does this array always have 1 element?
	return first.(map[string]interface{})
}
