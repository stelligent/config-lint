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
	loaded, err := loader.Load("./testdata/resources/uses_variables.tf")
	if err != nil {
		t.Error("Expecting TestTerraformLinter.Load to not return an error")
	}
	resources, err := loader.ReplaceVariables(loaded.Resources, loaded.Variables)
	if err != nil {
		t.Error("Expecting TestTerraformLinter.ReplaceVariables to not return an error")
	}
	if len(resources) != 1 {
		t.Errorf("Expecting to load 1 resources, not %d", len(loaded.Resources))
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

func TestTerraformVariablesInDifferentFile(t *testing.T) {
	options := Options{
		Tags:    []string{},
		RuleIDs: []string{},
	}
	filenames := []string{
		"./testdata/resources/defines_variables.tf",
		"./testdata/resources/reference_variables.tf",
	}
	linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: TerraformResourceLoader{}}
	ruleSet := loadRulesForTest("./testdata/rules/terraform_instance.yml", t)
	report, err := linter.Validate(ruleSet, options)
	if err != nil {
		t.Error("Expecting TestTerraformVariablesInDifferentFile to not return an error")
	}
	if len(report.ResourcesScanned) != 1 {
		t.Errorf("TestTerraformVariablesInDifferentFile scanned %d resources, expecting 1", len(report.ResourcesScanned))
	}
	if len(report.FilesScanned) != 2 {
		t.Errorf("TestTerraformVariablesInDifferentFile scanned %d files, expecting 2", len(report.FilesScanned))
	}
	if len(report.Violations) != 0 {
		t.Errorf("TestTerraformVariablesInDifferentFile expecting no violations, found %v", report.Violations)
	}
}

type TestingValueSource struct{}

func (s TestingValueSource) GetValue(a assertion.Expression) (string, error) {
	if a.ValueFrom.URL != "" {
		return "TEST", nil
	}
	return a.Value, nil
}

func TestTerraformPolicies(t *testing.T) {
	options := Options{
		Tags:    []string{},
		RuleIDs: []string{},
	}
	filenames := []string{"./testdata/resources/terraform_policy.tf"}
	linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: TerraformResourceLoader{}}
	ruleSet := loadRulesForTest("./testdata/rules/terraform_policy.yml", t)
	report, err := linter.Validate(ruleSet, options)
	if err != nil {
		t.Error("Expecting TestTerraformPolicies to not return an error")
	}
	if len(report.ResourcesScanned) != 1 {
		t.Errorf("TestTerraformPolicies scanned %d resources, expecting 1", len(report.ResourcesScanned))
	}
	if len(report.FilesScanned) != 1 {
		t.Errorf("TestTerraformPolicies scanned %d files, expecting 1", len(report.FilesScanned))
	}
	if len(report.Violations) != 1 {
		t.Errorf("TestTerraformPolicies returned %d violations, expecting 1", len(report.Violations))
	}
}

func TestTerraformPoliciesWithVariables(t *testing.T) {
	options := Options{
		Tags:    []string{},
		RuleIDs: []string{},
	}
	filenames := []string{"./testdata/resources/policy_with_variables.tf"}
	linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: TerraformResourceLoader{}}
	ruleSet := loadRulesForTest("./testdata/rules/policy_variable.yml", t)
	report, err := linter.Validate(ruleSet, options)
	if err != nil {
		t.Error("Expecting TestTerraformPoliciesWithVariables to not return an error:" + err.Error())
	}
	if len(report.Violations) != 0 {
		t.Errorf("TestTerraformPoliciesWithVariables returned %d violations, expecting 0", len(report.Violations))
	}
}
