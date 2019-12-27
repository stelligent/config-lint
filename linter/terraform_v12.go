package linter

import (
	"github.com/hashicorp/hcl/v2/hclparse"
	//"github.com/ghodss/yaml"
	"github.com/stelligent/config-lint/assertion"

	"github.com/hashicorp/hcl/v2"
	//hclsyntax "github.com/hashicorp/hcl/v2/hclsyntax"
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
		loaded.Resources = append(loaded.Resources, element.(assertion.Resource))
	}

	assertion.DebugJSON("loaded.Resources", loaded.Resources)

	return loaded, nil
}

func loadHCLv2(filename string) (Terraform12LoadResult, error) {
	result := Terraform12LoadResult{
		Resources: []interface{}{},
		Data:      []interface{}{},
		Providers: []interface{}{},
		Modules:   []interface{}{},
		Variables: []Variable{},
	}

	var file *hcl.File
	parser := hclparse.NewParser()
	file, _ = parser.ParseHCLFile(filename)
	content, _ := file.Body.Content(terraformSchema)
	resourceBlocks := content.Blocks.OfType("resource")
	//TODO: Only getting the first resource Block here, almost certainly need to recurse
	resource := assertion.Resource{
		ID:         resourceBlocks[0].Labels[0],
		Type:       resourceBlocks[0].Type,
		Category:   "",
		Filename:   filename,
		LineNumber: 0,
	}
	Variable{
		Name:  "",
		Value: nil,
	}
	props := make(map[string]interface{})
	resource.Properties = props
	props["ami"] = "ami-f2d3638a"
	result.Resources = append(result.Resources, resource)

	assertion.Debugf("LoadHCL Variables: %v\n", result.Variables)
	return result, nil
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
