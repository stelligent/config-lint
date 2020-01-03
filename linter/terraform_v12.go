package linter

import (
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
			Properties: attributesToMap(*block),
			Filename:   block.Range().Filename,
			LineNumber: block.Range().StartLine,
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
			Properties: attributesToMap(*block),
			Filename:   block.Range().Filename,
			LineNumber: block.Range().StartLine,
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
			Properties: attributesToMap(*block),
			Filename:   block.Range().Filename,
			LineNumber: block.Range().StartLine,
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

func attributesToMap(block tf12parser.Block) map[string]interface{} {
	propertyMap := make(map[string]interface{})
	for _, block := range block.AllBlocks() {
		var toAppend []interface{}
		toAppend = append(toAppend, attributesToMap(*block))
		propertyMap[block.Type()] = toAppend
	}
	attributes := block.GetAttributes()
	for _, attribute := range attributes {
		if attribute.Value().CanIterateElements() {
			var innerArray []interface{}
			innerMap := make(map[string]interface{})
			innerArray = append(innerArray, innerMap)
			propertyMap[attribute.Name()] = innerArray

			iter := attribute.Value().ElementIterator()
			for iter.Next() {
				key, value := iter.Element()
				if value.Type().HasDynamicTypes() {
					innerMap[ctyValueToString(key)] = ""
				} else {
					innerMap[ctyValueToString(key)] = ctyValueToString(value)
				}
			}
		} else {
			propertyMap[attribute.Name()] = ctyValueToString(attribute.Value())
		}
	}
	return propertyMap
}

func ctyValueToString(value cty.Value) string {
	switch value.Type() {
	case cty.NilType:
		return ""
	case cty.Bool:
		if value.True() {
			return "true"
		} else {
			return "false"
		}
	case cty.String:
		return value.AsString()
	case cty.Number:
		if value.RawEquals(cty.PositiveInfinity) || value.RawEquals(cty.NegativeInfinity) {
			panic("cannot convert infinity to string")
		}
		return value.AsBigFloat().Text('f', -1)
	default:
		panic("unsupported primitive type")
	}
}

// PostLoad resolves variable expressions
func (l Terraform12ResourceLoader) PostLoad(inputResources FileResources) ([]assertion.Resource, error) {
	return inputResources.Resources, nil
}
