package linter

import (
	"github.com/stelligent/config-lint/assertion"
	"testing"
)

func TestTerraformLinter(t *testing.T) {
	options := Options{
		Tags:    []string{},
		RuleIDs: []string{},
	}
	filenames := []string{"./testdata/resources/terraform_instance.tf"}
	linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: TerraformResourceLoader{}}
	ruleSet := loadRulesForTest("./testdata/rules/terraform_instance.yml", t)
	report, err := linter.Validate(ruleSet, options)
	if err != nil {
		t.Error("Expecting TestTerraformLinter to not return an error")
	}
	if len(report.ResourcesScanned) != 1 {
		t.Errorf("TestTerraformLinter scanned %d resources, expecting 1", len(report.ResourcesScanned))
	}
	if len(report.FilesScanned) != 1 {
		t.Errorf("TestTerraformLinter scanned %d files, expecting 1", len(report.FilesScanned))
	}
	if len(report.Violations) != 0 {
		t.Errorf("TestTerraformLinter returned %d violations, expecting 0", len(report.Violations))
	}
}

func TestTerraformVariables(t *testing.T) {
	loader := TerraformResourceLoader{}
	resources, err := loader.Load("./testdata/resources/uses_variables.tf")
	if err != nil {
		t.Error("Expecting TestTerraformLinter to not return an error")
	}
	if len(resources) != 1 {
		t.Errorf("Expecting to load 1 resources, not %d", len(resources))
	}
	properties := resources[0].Properties.(map[string]interface{})
	if properties["ami"] != "ami-f2d3638a" {
		t.Errorf("Unexpected value for variable: %s", properties["ami"])
	}
	// this test covers string, map, and slice cases
	tags := properties["tags"].([]interface{})
	tag := tags[0].(map[string]interface{})
	project := tag["project"].(string)
	if project != "demo" {
		t.Errorf("Expected project tag to be 'demo', got: %s", project)
	}
}

type TestingValueSource struct{}

func (s TestingValueSource) GetValue(a assertion.Expression) (string, error) {
	if a.ValueFrom.URL != "" {
		return "TEST", nil
	}
	return a.Value, nil
}
