package linter

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/parser"
	"github.com/stelligent/config-lint/assertion"
	"io/ioutil"
	"os"
	"path/filepath"
)

type (
	// TerraformResourceLoader converts Terraform configuration files into JSON objects
	TerraformResourceLoader struct{}

	// TerraformLoadResult collects all the returns value for parsing an HCL string
	TerraformLoadResult struct {
		Resources []interface{}
		Data      []interface{}
		Providers []interface{}
		Modules   []interface{}
		Variables []Variable
		AST       *ast.File
	}
)

func loadHCL(filename string) (TerraformLoadResult, error) {
	result := TerraformLoadResult{
		Resources: []interface{}{},
		Data:      []interface{}{},
		Providers: []interface{}{},
		Modules:   []interface{}{},
		Variables: []Variable{},
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
	result.Variables = append(loadVariables(m["variable"]), loadLocalVariables(m["locals"])...)
	if m["resource"] != nil {
		result.Resources = append(result.Resources, m["resource"].([]interface{})...)
	}
	if m["data"] != nil {
		result.Data = append(result.Data, m["data"].([]interface{})...)
	}
	if m["provider"] != nil {
		result.Providers = append(result.Providers, m["provider"].([]interface{})...)
	}
	if m["module"] != nil {
		result.Modules = append(result.Modules, m["module"].([]interface{})...)
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
		for key, resource := range m {
			variables = append(variables, Variable{Name: "var." + key, Value: getVariableValue(key, resource)})
		}
	}
	return variables
}

func loadLocalVariables(data interface{}) []Variable {
	variables := []Variable{}
	if data == nil {
		return variables
	}
	list := data.([]interface{})
	for _, entry := range list {
		m := entry.(map[string]interface{})
		for key, value := range m {
			variables = append(variables, Variable{Name: "local." + key, Value: value})
		}
	}
	return variables

}

func getVariableValue(key string, resource interface{}) interface{} {
	value := getVariableFromEnvironment(key)
	if value != "" {
		return value
	}
	return getVariableDefault(resource)
}

func getVariableFromEnvironment(key string) interface{} {
	return os.Getenv("TF_VAR_" + key)
}

func getVariableDefault(resource interface{}) interface{} {
	if list, ok := resource.([]interface{}); ok {
		var defaultValue interface{}
		for _, entry := range list {
			m := entry.(map[string]interface{})
			defaultValue = m["default"]
		}
		return flattenMaps(defaultValue)
	}
	return ""
}

func flattenMaps(v interface{}) interface{} {
	// map values are wrapped in an array, WAT?
	if listValue, ok := v.([]interface{}); ok {
		if len(listValue) == 1 {
			if mapValue, ok := listValue[0].(map[string]interface{}); ok {
				return mapValue
			}
		}
	}
	return v
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
	loaded.Resources = append(loaded.Resources, getResources(filename, result.AST, addIDToProviders(result.Providers), "provider")...)
	loaded.Resources = append(loaded.Resources, getResources(filename, result.AST, addKeyToModules(result.Modules), "module")...)

	assertion.DebugJSON("loaded.Resources", loaded.Resources)

	return loaded, nil
}

// Providers do not have an name, so generate one to make the data format the same as resources

func addIDToProviders(providers []interface{}) []interface{} {
	resources := []interface{}{}
	for _, provider := range providers {
		resources = append(resources, addIDToProvider(provider))
	}
	return resources
}

func addIDToProvider(provider interface{}) interface{} {
	m := provider.(map[string]interface{})
	for providerType, value := range m {
		m[providerType] = addIDToProviderValue(value)
	}
	return m
}

// Counter used to generate an ID for providers
var Counter = 0

func addIDToProviderValue(value interface{}) interface{} {
	Counter++
	m := map[string]interface{}{}
	key := fmt.Sprintf("%d", Counter)
	m[key] = value
	return []interface{}{m}
}

// use the source attribute of modules as the key
func addKeyToModules(modules []interface{}) []interface{} {
	resources := map[string]interface{}{}
	for _, module := range modules {
		resources = addKeyToModule(resources, module)
	}
	return []interface{}{resources}
}

func addKeyToModule(resources map[string]interface{}, module interface{}) map[string]interface{} {
	m := module.(map[string]interface{})
	for moduleName, valueList := range m {
		a := valueList.([]interface{})
		for _, value := range a {
			properties := value.(map[string]interface{})
			source := properties["source"].(string)

			inner := []interface{}{properties}

			outer := map[string]interface{}{}
			outer[moduleName] = inner

			existing, ok := resources[source]
			if ok {
				list := existing.([]interface{})
				resources[source] = append(list, outer)
			} else {
				resources[source] = []interface{}{outer}
			}
		}
	}
	return resources
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
						properties["__file__"] = filename
						properties["__dir__"] = filepath.Dir(filename)
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

// PostLoad resolves variable expressions
func (l TerraformResourceLoader) PostLoad(fr FileResources) ([]assertion.Resource, error) {
	for _, resource := range fr.Resources {
		resource.Properties = replaceVariables(resource.Properties, fr.Variables)
	}
	for _, resource := range fr.Resources {
		properties, err := parseJSONDocuments(resource.Properties)
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
	case string:
		return interpolate(v, variables)
	default:
		assertion.Debugf("replaceVariables cannot process type %T: %v\n", v, v)
		return templateResource
	}
}

func replaceVariablesInMap(templateResource map[string]interface{}, variables []Variable) interface{} {
	for key, value := range templateResource {
		switch v := value.(type) {
		case string:
			templateResource[key] = interpolate(v, variables)
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

func parseJSONDocuments(resource interface{}) (interface{}, error) {
	properties := resource.(map[string]interface{})
	for _, attribute := range []string{"assume_role_policy", "policy", "container_definitions", "access_policies"} {
		if policyAttribute, hasPolicyString := properties[attribute]; hasPolicyString {
			if policyString, isString := policyAttribute.(string); isString {
				var policy interface{}
				if policyString != "" {
					err := json.Unmarshal([]byte(policyString), &policy)
					if err != nil {
						assertion.Debugf("Unable to parse '%s' as JSON\n", policyString)
						assertion.Debugf("Error: %v\n", err)
					}
				}
				properties[attribute] = policy
			}
		}
	}
	return properties, nil
}

func getProperties(templateResource interface{}) map[string]interface{} {
	switch v := templateResource.(type) {
	case []interface{}:
		return v[0].(map[string]interface{})
	default:
		return map[string]interface{}{}
	}
}
