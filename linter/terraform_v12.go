package linter

import (
	"github.com/zclconf/go-cty/cty"
	"strconv"
	"strings"

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

	// Get all Terraform blocks of a given type and append to the slice of Resources
	blockTypes := []string{"resource", "provider", "data", "module"}
	for _, blockType := range blockTypes {
		result.Resources = append(result.Resources, getBlocksOfType(blocks, blockType)...)
	}

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

// Retrieves Terraform blocks of a specific type
// and places them in a slice of assertion.Resources.
func getBlocksOfType(blocks tf12parser.Blocks, blockType string) []assertion.Resource {
	var id string
	var resources []assertion.Resource

	tfBlocks := blocks.OfType(blockType)
	i := 0

	// If there is no Terraform ID for a block (in the case of Providers and Modules),
	// set the id variable to an auto-incrementing integer.
	// Otherwise, set it to the second item in the block's Labels.
	for _, block := range tfBlocks {
		if len(block.Labels()) > 1 {
			id = block.Labels()[1]
		} else {
			id = strconv.Itoa(i)
			i++
		}
		if block.Type() != "module" {
			resource := assertion.Resource{
				ID: id,
				Type:       block.Labels()[0],
				Category:   blockType,
				Properties: attributesToMap(*block),
				Filename:   block.Range().Filename,
				LineNumber: block.Range().StartLine,
			}
			resources = append(resources, resource)
		} else {
			moduleSource := block.GetAttribute("source")
			resource := assertion.Resource{
				ID: block.Labels()[0],
				Type:       moduleSource.Value().AsString(),
				Category:   blockType,
				Properties: attributesToMap(*block),
				Filename:   block.Range().Filename,
				LineNumber: block.Range().StartLine,
			}
			resources = append(resources, resource)
		}
	}
	return resources
}

func attributesToMap(block tf12parser.Block) map[string]interface{} {
	propertyMap := make(map[string]interface{})
	allBlocks := block.AllBlocks()
	for _, block := range allBlocks {
		var toAppend []interface{}
		toAppend = append(toAppend, attributesToMap(*block))
		if propertyMap[block.Type()] == nil {
			propertyMap[block.Type()] = toAppend
		} else {
			v := propertyMap[block.Type()].([]interface{})
			v = append(v, toAppend[0])
			propertyMap[block.Type()] = v
		}
	}
	attributes := block.GetAttributes()
	for _, attribute := range attributes {
		value := attribute.Value()
		if value.Type().IsTupleType() {
			innerArray := make([]interface{}, 0)

			iter := value.ElementIterator()
			for iter.Next() {
				_, value := iter.Element()
				innerArray = append(innerArray, ctyValueToString(value))
			}
			propertyMap[attribute.Name()] = innerArray
		} else if value.CanIterateElements() {
			iterateElements(propertyMap, attribute.Name(), value)
		} else {
			setValue(propertyMap, attribute.Name(), ctyValueToString(value))
		}
	}
	return propertyMap
}

func iterateElements(propertyMap map[string]interface{}, name string, value cty.Value) {
	var innerArray []interface{}
	innerMap := make(map[string]interface{})
	innerArray = append(innerArray, innerMap)
	propertyMap[name] = innerArray

	iter := value.ElementIterator()
	for iter.Next() {
		key, value := iter.Element()
		if value.CanIterateElements() {
			iterateElements(innerMap, ctyValueToString(key), value)
		} else {
			setValue(innerMap, ctyValueToString(key), ctyValueToString(value))
		}
	}
}

func setValue(m map[string]interface{}, name string, value string) {
	environmentVariable := getVariableFromEnvironment(name)
	if environmentVariable == "" {
		m[name] = value
	} else {
		m[name] = environmentVariable
	}
}

func ctyValueToString(value cty.Value) string {
	// In case the value is nil but the type is not necessarily <nil>, return an empty string
	if value.IsNull() || !value.IsKnown() {
		return ""
	} else {
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
			return strings.Trim(value.AsString(), "\n")
		case cty.Number:
			if value.RawEquals(cty.PositiveInfinity) || value.RawEquals(cty.NegativeInfinity) {
				panic("cannot convert infinity to string")
			}
			return value.AsBigFloat().Text('f', -1)
		default:
			panic("unsupported primitive type")
		}
	}
}

// PostLoad resolves variable expressions
func (l Terraform12ResourceLoader) PostLoad(inputResources FileResources) ([]assertion.Resource, error) {
	return inputResources.Resources, nil
}
