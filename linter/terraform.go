package linter

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/parser"
	"github.com/stelligent/config-lint/assertion"
	"io/ioutil"
)

// TerraformLinter implements a Linter for Terraform configuration files
type TerraformLinter struct {
	Filenames   []string
	ValueSource assertion.ValueSource
}

// TerraformResourceLoader converts Terraform configuration files into JSON objects
type TerraformResourceLoader struct{}

func parsePolicy(templateResource interface{}) (map[string]interface{}, error) {
	firstResource := templateResource.([]interface{})[0] // FIXME does this array always have 1 element?
	properties := firstResource.(map[string]interface{})
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

func loadHCL(filename string) ([]interface{}, *ast.File, error) {
	results := make([]interface{}, 0)
	template, err := ioutil.ReadFile(filename)
	if err != nil {
		return results, nil, nil
	}

	root, _ := parser.Parse(template)

	var v interface{}
	err = hcl.Unmarshal([]byte(template), &v)
	if err != nil {
		return results, root, nil
	}
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return results, root, nil
	}
	assertion.Debugf("LoadHCL: %s\n", string(jsonData))

	var hclData interface{}
	err = yaml.Unmarshal(jsonData, &hclData)
	if err != nil {
		return results, root, nil
	}
	m := hclData.(map[string]interface{})
	if m["resource"] != nil {
		results = append(results, m["resource"].([]interface{})...)
	}
	return results, root, nil
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
func (l TerraformResourceLoader) Load(filename string) ([]assertion.Resource, error) {
	resources := make([]assertion.Resource, 0)
	hclResources, root, err := loadHCL(filename)
	if err != nil {
		return resources, err
	}
	for _, resource := range hclResources {
		for resourceType, templateResources := range resource.(map[string]interface{}) {
			if templateResources != nil {
				for _, templateResource := range templateResources.([]interface{}) {
					for resourceID, templateResource := range templateResource.(map[string]interface{}) {
						properties, err := parsePolicy(templateResource)
						if err != nil {
							return resources, err
						}
						lineNumber := getResourceLineNumber(resourceType, resourceID, filename, root)
						tr := assertion.Resource{
							ID:         resourceID,
							Type:       resourceType,
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
	return resources, nil
}

// Validate uses a RuleSet to validate resources in a collection of Terraform configuration files
func (l TerraformLinter) Validate(ruleSet assertion.RuleSet, options Options) (assertion.ValidationReport, error) {
	loader := TerraformResourceLoader{}
	f := FileLinter{Filenames: l.Filenames, ValueSource: l.ValueSource, Loader: loader}
	return f.ValidateFiles(ruleSet, options)
}

// Search applies a JMESPath expression to the resources in a collection of Terraform configuration files
func (l TerraformLinter) Search(ruleSet assertion.RuleSet, searchExpression string) {
	loader := TerraformResourceLoader{}
	f := FileLinter{Filenames: l.Filenames, ValueSource: l.ValueSource, Loader: loader}
	f.SearchFiles(ruleSet, searchExpression)
}
