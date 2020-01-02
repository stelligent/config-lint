package linter

import (
	"fmt"
	"github.com/zclconf/go-cty/cty"
	"strconv"

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
//TODO: This should be unused, but can't remove due to the interface, I think?
func (l Terraform12ResourceLoader) Load(filename string) (FileResources, error) {
	loaded := FileResources{
		Resources: []assertion.Resource{},
	}
	result, err := loadHCLv2([]string{filename})
	if err != nil {
		return loaded, err
	}
	loaded.Resources = result.Resources

	assertion.DebugJSON("loaded.Resources", loaded.Resources)

	return loaded, nil
}

func (l Terraform12ResourceLoader) LoadMany(filenames []string) (FileResources, error) {
	loaded := FileResources{
		Resources: []assertion.Resource{},
	}
	result, err := loadHCLv2(filenames)
	if err != nil {
		return loaded, err
	}
	loaded.Resources = result.Resources

	assertion.DebugJSON("loaded.Resources", loaded.Resources)

	return loaded, nil
}

func loadHCLv2(paths []string) (Terraform12LoadResult, error) {
	result := Terraform12LoadResult{
		Resources: []assertion.Resource{},
		Data:      []interface{}{},
		Providers: []interface{}{},
		Modules:   []interface{}{},
		Variables: []Variable{},
	}

	parser := *tf12parser.New()
	blocks, err := parser.ParseMany(paths)
	if err != nil {
		return result, err
	}

	resourceBlocks := blocks.OfType("resource")
	for _, block := range resourceBlocks {
		resource := assertion.Resource{
			ID:         block.Labels()[1],
			Type:       block.Labels()[0],
			Category:   "resource",
			Properties: attributesToMap(block.GetAttributes()),
			Filename:   "",
			LineNumber: 0,
		}
		result.Resources = append(result.Resources, resource)
	}

	providerBlocks := blocks.OfType("provider")
	i := 0
	for _, block := range providerBlocks {
		resource := assertion.Resource{
			ID:         strconv.Itoa(i),
			Type:       block.Labels()[0],
			Category:   "provider",
			Properties: attributesToMap(block.GetAttributes()),
			Filename:   "",
			LineNumber: 0,
		}
		result.Resources = append(result.Resources, resource)
		i++
	}

	dataBlocks := blocks.OfType("data")
	for _, block := range dataBlocks {
		resource := assertion.Resource{
			ID:         block.Labels()[1],
			Type:       block.Labels()[0],
			Category:   "data",
			Properties: attributesToMap(block.GetAttributes()),
			Filename:   "",
			LineNumber: 0,
		}
		result.Resources = append(result.Resources, resource)
	}

	//dataBlocks := blocks.OfType("data")
	//for _, block := range dataBlocks {
	//	outerMap := attributesToMap(block.GetAttributes())
	//	for _, elem := range outerMap.(map[string]interface{}) {
	//		result.Data = append(result.Data, elem)
	//	}
	//}

	for _, resource := range result.Resources {
		properties, err := parseJSONDocuments(resource.Properties)
		if err != nil {
			return result, err
		}
		resource.Properties = properties
	}

	assertion.Debugf("LoadHCL Variables: %v\n", result.Variables)
	return result, nil
}

func attributesToMap(attributes []*tf12parser.Attribute) interface{} {
	propertyMap := make(map[string]interface{})
	for _, elem := range attributes {
		if elem.Value().CanIterateElements() {
			var innerArray []interface{}
			innerMap := make(map[string]interface{})
			innerArray = append(innerArray, innerMap)
			propertyMap[elem.Name()] = innerArray

			iter := elem.Value().ElementIterator()
			for iter.Next() {
				key, value := iter.Element()
				if value.Type().HasDynamicTypes() {
					innerMap[key.AsString()] = ""
				} else {
					innerMap[key.AsString()] = value.AsString()
				}
			}
		} else {
			if elem.Type() == cty.NilType {
				propertyMap[elem.Name()] = ""
			} else {
				fmt.Println(elem)
				propertyMap[elem.Name()] = elem.Value().AsString()
			}
		}
	}
	return propertyMap
}

// PostLoad resolves variable expressions
func (l Terraform12ResourceLoader) PostLoad(inputResources FileResources) ([]assertion.Resource, error) {
	return inputResources.Resources, nil
}
