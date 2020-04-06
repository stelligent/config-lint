package linter

import (
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/stelligent/config-lint/assertion"
	"github.com/stelligent/config-lint/linter/tf12parser"
	"github.com/zclconf/go-cty/cty"
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

var (
	blockTypes = []string{
		"data",
		"locals",
		"module",
		"output",
		"provider",
		"resource",
		"terraform",
		"variable",
	}

	blockLabelSyntax = map[string][]string{
		"TypeAndName":  []string{"data", "resource"},
		"TypeOnly":     []string{"provider"},
		"NameOnly":     []string{"module", "output", "variable"},
		"NoTypeNoName": []string{"locals", "terraform"},
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
func getBlocksOfType(blocks tf12parser.Blocks, blockCategory string) []assertion.Resource {
	var blockType string
	var blockName string
	var resources []assertion.Resource

	tfBlocks := blocks.OfType(blockCategory)
	i := 0

	for _, block := range tfBlocks {

		properties := attributesToMap(*block)

		// Terraform has labels between 0 and 2 for each block e.g `locals`, `provider`, and `resource`,
		// and labels could be fixed types or configurable names.
		// Thus this code checks which label should be used as type and as a name/id. If the block doesn't have
		// a unique name, then its name/id assigned to an auto-incrementing integer.
		if assertion.SliceContains(blockLabelSyntax["TypeAndName"], blockCategory) {
			blockType = block.Labels()[0]
			blockName = block.Labels()[1]
			properties["__type__"] = blockType
			properties["__name__"] = blockName

		} else if assertion.SliceContains(blockLabelSyntax["TypeOnly"], blockCategory) {
			blockType = block.Labels()[0]
			blockName = strconv.Itoa(i)
			i++
			properties["__type__"] = blockType
			properties["__name__"] = blockName

		} else if assertion.SliceContains(blockLabelSyntax["NameOnly"], blockCategory) {
			// A special handling for module to add its source as a type.
			if blockCategory == "module" {
				blockType = block.GetAttribute("source").Value().AsString()
			} else {
				blockType = blockCategory
			}
			blockName = block.Labels()[0]
			properties["__name__"] = blockName

		} else if assertion.SliceContains(blockLabelSyntax["NoTypeNoName"], blockCategory) {
			blockType = blockCategory
			blockName = strconv.Itoa(i)
			i++
		}

		resource := assertion.Resource{
			ID:         blockName,
			Type:       blockType,
			Category:   blockCategory,
			Properties: properties,
			Filename:   block.Range().Filename,
			LineNumber: block.Range().StartLine,
		}
		resources = append(resources, resource)
	}
	return resources
}

func attributesToMap(block tf12parser.Block) map[string]interface{} {
	propertyMap := make(map[string]interface{})
	allBlocks := block.AllBlocks()
	for _, currentBlock := range allBlocks {
		var toAppend []interface{}
		toAppend = append(toAppend, attributesToMap(*currentBlock))
		if propertyMap[currentBlock.Type()] == nil {
			propertyMap[currentBlock.Type()] = toAppend
		} else {
			v := propertyMap[currentBlock.Type()].([]interface{})
			v = append(v, toAppend[0])
			propertyMap[currentBlock.Type()] = v
		}
	}
	attributes := block.GetAttributes()
	for _, attribute := range attributes {
		value := attribute.Value()
		if value.Type().IsTupleType() {
			innerArray := make([]interface{}, 0)

			iter := value.ElementIterator()
			for iter.Next() {
				_, iterValue := iter.Element()
				if iterValue.CanIterateElements() {
					iterateElements(propertyMap, attribute.Name(), iterValue)
				} else {
					innerArray = append(innerArray, ctyValueToString(iterValue))
				}
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
	// In case the value is nil but the type is not necessarily <nil>, ~~return an empty string~~
	// Update: return an actual string.
	// We cannot evaluate tf generated values in tf12, such as referenced arn, but we still want to be able to check for it
	if value.IsNull() || !value.IsKnown() {
		return "UNDEFINED"
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
			//return ""
		}
	}
}

// PostLoad resolves variable expressions
func (l Terraform12ResourceLoader) PostLoad(inputResources FileResources) ([]assertion.Resource, error) {
	return inputResources.Resources, nil
}
