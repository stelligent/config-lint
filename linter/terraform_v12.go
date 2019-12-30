package linter

import (
	"fmt"
	"github.com/zclconf/go-cty/cty"

	//"github.com/ghodss/yaml"
	"github.com/stelligent/config-lint/assertion"

	"github.com/hashicorp/hcl/v2"
	//hclsyntax "github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stelligent/config-lint/linter/tf12parser"
)

type (
	// Terraform12ResourceLoader converts Terraform configuration files into JSON objects
	Terraform12ResourceLoader struct{}

	// Terraform12LoadResult collects all the returns value for parsing an HCL string
	Terraform12LoadResult struct {
		Resources []assertion.Resource
		Data      []interface{}
		Providers []interface{}
		Modules   []interface{}
		Variables []Variable
		AST       *hcl.File
	}
)

// Load parses an HCLv2 file into a collection or Resource objects
func (l Terraform12ResourceLoader) Load(filename string) (FileResources, error) {
	loaded := FileResources{
		Resources: []assertion.Resource{},
	}
	result, err := loadHCLv2(filename)
	if err != nil {
		return loaded, err
	}
	//TODO: MN -- Need to iterate over range here to append all? Seems like I'm missing a GoLang idiom
	for _, element := range result.Resources {
		loaded.Resources = append(loaded.Resources, element)
	}

	assertion.DebugJSON("loaded.Resources", loaded.Resources)

	return loaded, nil
}

type (
	Instance struct {
		Ami string `cty:"ami"`
		InstanceType string `cty:"instance_type"`
		Tags map[string]string `cty:"tags"`
	}
)

func loadHCLv2(filename string) (Terraform12LoadResult, error) {
	result := Terraform12LoadResult{
		Resources: []assertion.Resource{},
		Data:      []interface{}{},
		Providers: []interface{}{},
		Modules:   []interface{}{},
		Variables: []Variable{},
	}

	parser := *tf12parser.New()
	context, err := parser.ParseFile(filename)
	if err != nil {
		fmt.Println("Boo!")
	}


	resources := context.Variables["resource"]
	mapResources(resources)
	//resourceStruct := new(Instance)
	//err = gocty.FromCtyValue(resources, resourceStruct)
	//Note: values are not processing in a consistent order. If there's any error, the entire result is invalid
	if err != nil {
		fmt.Println("Boo!")
	}

	//fmt.Println(resourceStruct.Ami)
	//fmt.Println(resourceStruct.InstanceType)

	assertion.Debugf("LoadHCL Variables: %v\n", result.Variables)
	return result, nil
}

func mapResources(value cty.Value) {
	var resources []assertion.Resource
	it := value.ElementIterator()
	for it.Next() {
		key, _ := it.Element()
		resource := assertion.Resource{
			ID:         "",
			Type:       key.AsString(),
			Category:   "resource",
			Properties: nil,
			Filename:   "",
			LineNumber: 0,
		}
		resources = append(resources, resource)
	}
	//for v := range value {
	//
	//}
	//var resources []assertion.Resource
	//
	//cty.Walk(value, func(path cty.Path, value cty.Value) (bool, error) {
	//	level := len(path)
	//	if level == 1 {
	//		resources = append(resources, assertion.Resource{
	//			Type: path[len(path)-1].(cty.GetAttrStep).Name,
	//		})
	//	}
	//	return true, nil
	//})
	//fmt.Println(resources)
}

func getNameAtLevel(path cty.Path, value cty.Value) (b bool, err error) {
	level := len(path)
	if level == 1 {
		fmt.Println(path[len(path)-1].(cty.GetAttrStep).Name)
	}
	return true, nil
}

// PostLoad resolves variable expressions
func (l Terraform12ResourceLoader) PostLoad(inputResources FileResources) ([]assertion.Resource, error) {
	//for _, resource := range inputResources.Resources {
	//	resource.Properties = tf12ReplaceVariables(resource.Properties, inputResources.Variables)
	//}
	//for _, resource := range inputResources.Resources {
	//	properties, err := tf12ParseJSONDocuments(resource.Properties)
	//	if err != nil {
	//		return inputResources.Resources, err
	//	}
	//	resource.Properties = properties
	//}
	return inputResources.Resources, nil
}
