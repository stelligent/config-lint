package linter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/stelligent/config-lint/assertion"

	hcl2 "github.com/hashicorp/hcl/v2"
	hcl2parse "github.com/hashicorp/hcl/v2/hclparse"
)

type (
	// Terraform12ResourceLoader converts Terraform configuration files into JSON objects
	Terraform12ResourceLoader struct{}

	// Terraform12LoadResult collects all the returns value for parsing an HCL string
	Terraform12LoadResult struct {
		Resources []interface{}
		Data      []interface{}
		Providers []interface{}
		Modules   []interface{}
		Variables []Variable
		AST       *hcl2.File
	}
)

func loadHCLv2(filename string) (Terraform12LoadResult, error) {
	diags := hcl2.Diagnostics{}
	result := Terraform12LoadResult{
		Resources: []interface{}{},
		Data:      []interface{}{},
		Providers: []interface{}{},
		Modules:   []interface{}{},
		Variables: []Variable{},
	}

	// NEW PARSER FOR HCL V2 (hclparse)
	result.AST, diags = hcl2parse.NewParser().ParseHCLFile(filename)
	if diags.HasErrors() {
		fmt.Printf("ERROR: %v\n", diags)
		return result, diags
	}

	// PRINT OUT STRING CONVERSION OF 'result.AST'
	//fmt.Printf("RESULT AST:\n %v\n", string(result.AST.Bytes))

	// BODY CONTENT IN THE FORMAT OF SCHEMA FROM 'schema.go'
	hcl2BodyContent, _ := result.AST.Body.Content(terraformSchema)
	// fmt.Printf("BODY CONTENT:\n %v\n", hcl2BodyContent)

	// BODY BLOCKS WITHIN BODY CONTENT
	hcl2BodyBlocks := hcl2BodyContent.Blocks.ByType()
	// fmt.Printf("BODY BLOCKS:\n %v\n", hcl2BodyBlocks)

	// RETURNS THE JSON ENCODING OF THE 'hcl2BodyBlocks' map[string]hcl2.blocks (hcl2BodyContent.Blocks.ByType())
	hcl2JSONEncoding, err := json.Marshal(hcl2BodyBlocks)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Printf("JSON ENCODING STRING VALUE CONVERSION FROM BYTE SLICE:\n %v\n", string(hcl2JSONEncoding))

	// TAKES THE BYTE SLICE FORMAT OF THE JSON ENCODING (hcl2JSONEncoding) AND UNMARSHALS IT TO 'var hcl2Data interface{}'
	var hcl2Data interface{}
	err = json.Unmarshal([]byte(hcl2JSONEncoding), &hcl2Data)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return result, err
	}
	// 'hcl2Data interface{}' VALUE AFTER BEING UNMARSHALED TO FROM 'hcl2JSONEncoding'
	// fmt.Printf("VAR HCL2DATA INTERFACE{}:\n %v\n", hcl2Data)

	// *****************************************************************************************************************************
	// BELOW NOT NECESSARY. UNCOMMENT FOR OUTPUTTING READABLE FORMAT OF JSON STRING AND FOR DEBUGGING PURPOSES
	// *****************************************************************************************************************************
	// jsonData, err := json.MarshalIndent(v, "", "  ")
	// if err != nil {
	// 	fmt.Printf("JSON DATA ERR: %v\n", err)
	// 	return result, err
	// }
	// fmt.Printf("jsonData IN []BYTE FORM AFTER MARSHAL INDENT:\n %v\n", jsonData)
	// fmt.Printf("jsonData IN JSON STRING FORMAT:\n %v\n", string(jsonData))
	// assertion.Debugf("LoadHCL: %s\n", string(jsonData))

	// var hcl2Data interface{}
	// err = yaml.Unmarshal(jsonData, &hcl2Data)
	// if err != nil {
	// 	return result, err
	// }
	// fmt.Printf("hcl2Data INTERFACE:\n %v\n", hcl2Data)
	// *****************************************************************************************************************************
	// END
	// *****************************************************************************************************************************

	m := hcl2Data.(map[string]interface{})

	result.Variables = append(tf12LoadVariables(m["variable"]), tf12LoadLocalVariables(m["locals"])...)
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

func tf12LoadVariables(data interface{}) []Variable {
	variables := []Variable{}
	if data == nil {
		return variables
	}
	list := data.([]interface{})
	for _, entry := range list {
		m := entry.(map[string]interface{})
		for key, resource := range m {
			variables = append(variables, Variable{Name: "var." + key, Value: tf12GetVariableValue(key, resource)})
		}
	}
	return variables
}

func tf12LoadLocalVariables(data interface{}) []Variable {
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

func tf12GetVariableValue(key string, resource interface{}) interface{} {
	value := tf12GetVariableFromEnvironment(key)
	if value != "" {
		return value
	}
	return tf12GetVariableDefault(resource)
}

func tf12GetVariableFromEnvironment(key string) interface{} {
	return os.Getenv("TF_VAR_" + key)
}

func tf12GetVariableDefault(resource interface{}) interface{} {
	if list, ok := resource.([]interface{}); ok {
		var defaultValue interface{}
		for _, entry := range list {
			m := entry.(map[string]interface{})
			defaultValue = m["default"]
		}
		return tf12FlattenMaps(defaultValue)
	}
	return ""
}

func tf12FlattenMaps(v interface{}) interface{} {
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

// ** TODO: Create func/logic to make sure this is grabbing the desired resource line based on 'resourceType' and 'resourceID' **
func tf12GetResourceLineNumber(resourceType, resourceID, filename string, root *hcl2.File) int {
	hcl2BodyContent, _ := root.Body.Content(terraformSchema)
	hcl2ResourceBlocks := hcl2BodyContent.Blocks.OfType("resource")
	if len(hcl2ResourceBlocks) > 0 {
		tf12ResourceBlockRange := hcl2ResourceBlocks[0].TypeRange
		assertion.Debugf("Position %s %s:%d\n", resourceID, filename, tf12ResourceBlockRange.Start.Line)
		return tf12ResourceBlockRange.Start.Line
	}
	return 0
}

// Load parses an HCLv2 file into a collection or Resource objects
func (l Terraform12ResourceLoader) Load(filename string) (FileResources, error) {
	loaded := FileResources{
		Resources: []assertion.Resource{},
	}
	result, err := loadHCLv2(filename)
	if err != nil {
		return loaded, err
	}
	loaded.Variables = result.Variables
	loaded.Resources = append(loaded.Resources, tf12GetResources(filename, result.AST, result.Resources, "resource")...)
	loaded.Resources = append(loaded.Resources, tf12GetResources(filename, result.AST, result.Data, "data")...)
	loaded.Resources = append(loaded.Resources, tf12GetResources(filename, result.AST, tf12AddIDToProviders(result.Providers), "provider")...)
	loaded.Resources = append(loaded.Resources, tf12GetResources(filename, result.AST, tf12AddKeyToModules(result.Modules), "module")...)

	assertion.DebugJSON("loaded.Resources", loaded.Resources)

	return loaded, nil
}

// Providers do not have an name, so generate one to make the data format the same as resources
func tf12AddIDToProviders(providers []interface{}) []interface{} {
	resources := []interface{}{}
	for _, provider := range providers {
		resources = append(resources, tf12AddIDToProvider(provider))
	}
	return resources
}

func tf12AddIDToProvider(provider interface{}) interface{} {
	m := provider.(map[string]interface{})
	for providerType, value := range m {
		m[providerType] = tf12AddIDToProviderValue(value)
	}
	return m
}

// TF12Counter used to generate an ID for providers
var TF12Counter = 0

func tf12AddIDToProviderValue(value interface{}) interface{} {
	TF12Counter++
	m := map[string]interface{}{}
	key := fmt.Sprintf("%d", TF12Counter)
	m[key] = value
	return []interface{}{m}
}

// use the source attribute of modules as the key
func tf12AddKeyToModules(modules []interface{}) []interface{} {
	resources := map[string]interface{}{}
	for _, module := range modules {
		resources = tf12AddKeyToModule(resources, module)
	}
	return []interface{}{resources}
}

func tf12AddKeyToModule(resources map[string]interface{}, module interface{}) map[string]interface{} {
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

func tf12GetResources(filename string, ast *hcl2.File, objects []interface{}, category string) []assertion.Resource {
	resources := []assertion.Resource{}
	for _, resource := range objects {
		for resourceType, templateResources := range resource.(map[string]interface{}) {
			if templateResources != nil {
				for _, templateResource := range templateResources.([]interface{}) {
					for resourceID, templateResource := range templateResource.(map[string]interface{}) {
						properties := tf12GetProperties(templateResource)
						lineNumber := tf12GetResourceLineNumber(resourceType, resourceID, filename, ast)
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
func (l Terraform12ResourceLoader) PostLoad(fr FileResources) ([]assertion.Resource, error) {
	for _, resource := range fr.Resources {
		resource.Properties = tf12ReplaceVariables(resource.Properties, fr.Variables)
	}
	for _, resource := range fr.Resources {
		properties, err := tf12ParseJSONDocuments(resource.Properties)
		if err != nil {
			return fr.Resources, err
		}
		resource.Properties = properties
	}
	return fr.Resources, nil
}

func tf12ReplaceVariables(templateResource interface{}, variables []Variable) interface{} {
	switch v := templateResource.(type) {
	case map[string]interface{}:
		return tf12ReplaceVariablesInMap(v, variables)
	case string:
		return interpolate(v, variables)
	default:
		assertion.Debugf("tf12ReplaceVariables cannot process type %T: %v\n", v, v)
		return templateResource
	}
}

func tf12ReplaceVariablesInMap(templateResource map[string]interface{}, variables []Variable) interface{} {
	for key, value := range templateResource {
		switch v := value.(type) {
		case string:
			templateResource[key] = interpolate(v, variables)
		case map[string]interface{}:
			templateResource[key] = tf12ReplaceVariablesInMap(v, variables)
		case []interface{}:
			templateResource[key] = tf12ReplaceVariablesInList(v, variables)
		default:
			assertion.Debugf("tf12ReplaceVariablesInMap cannot process type %T\n", v)
		}
	}
	return templateResource
}

func tf12ReplaceVariablesInList(list []interface{}, variables []Variable) []interface{} {
	result := []interface{}{}
	for _, e := range list {
		result = append(result, tf12ReplaceVariables(e, variables))
	}
	return result
}

func tf12ParseJSONDocuments(resource interface{}) (interface{}, error) {
	properties := resource.(map[string]interface{})
	for _, attribute := range []string{"assume_role_policy", "policy", "container_definitions", "access_policies", "container_properties"} {
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

func tf12GetProperties(templateResource interface{}) map[string]interface{} {
	switch v := templateResource.(type) {
	case []interface{}:
		return v[0].(map[string]interface{})
	default:
		return map[string]interface{}{}
	}
}
