package linter

import (
	"github.com/stelligent/config-lint/assertion"
	"github.com/stretchr/testify/assert"
	"os"
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
	assert.Nil(t, err, "Expecting Validate to run without error")
	assert.Equal(t, len(report.ResourcesScanned), 1, "Unexpected number of resources scanned")
	assert.Equal(t, len(report.FilesScanned), 1, "Unexpected number of files scanned")
	assertViolationsCount("TestTerraformLinter ", 0, report.Violations, t)
}

func loadResourcesToTest(t *testing.T, filename string) []assertion.Resource {
	loader := TerraformResourceLoader{}
	loaded, err := loader.Load(filename)
	assert.Nil(t, err, "Expecting Load to run without error")
	resources, err := loader.PostLoad(loaded)
	assert.Nil(t, err, "Expecting PostLoad to run without error")
	return resources
}

func getResourceTags(r assertion.Resource) map[string]interface{} {
	properties := r.Properties.(map[string]interface{})
	tags := properties["tags"].([]interface{})
	return tags[0].(map[string]interface{})
}

func TestTerraformVariable(t *testing.T) {
	resources := loadResourcesToTest(t, "./testdata/resources/uses_variables.tf")
	assert.Equal(t, len(resources), 1, "Expecting 1 resource")
	properties := resources[0].Properties.(map[string]interface{})
	assert.Equal(t, properties["ami"], "ami-f2d3638a", "Unexpected value for simple variable")
}

func TestTerraformVariableWithNoDefault(t *testing.T) {
	resources := loadResourcesToTest(t, "./testdata/resources/uses_variables.tf")
	assert.Equal(t, len(resources), 1, "Expecting 1 resource")
	tags := getResourceTags(resources[0])
	assert.Equal(t, tags["department"], "", "Unexpected value for variable with no default")
}

func TestTerraformFunctionCall(t *testing.T) {
	resources := loadResourcesToTest(t, "./testdata/resources/uses_variables.tf")
	assert.Equal(t, len(resources), 1, "Expecting 1 resource")
	tags := getResourceTags(resources[0])
	assert.Equal(t, tags["environment"], "test", "Unexpected value for lookup function")
}

func TestTerraformListVariable(t *testing.T) {
	resources := loadResourcesToTest(t, "./testdata/resources/uses_variables.tf")
	assert.Equal(t, len(resources), 1, "Expecting 1 resource")
	tags := getResourceTags(resources[0])
	assert.Equal(t, tags["comment"], "bar", "Unexpected value for list variable")
}

func TestTerraformLocalVariable(t *testing.T) {
	resources := loadResourcesToTest(t, "./testdata/resources/uses_local_variables.tf")
	assert.Equal(t, len(resources), 1, "Expecting 1 resource")
	properties := resources[0].Properties.(map[string]interface{})
	assert.Equal(t, "myprojectbucket", properties["name"], "Unexpected value for name attribute")
}

func TestTerraformVariablesFromEnvironment(t *testing.T) {
	os.Setenv("TF_VAR_instance_type", "c4.large")
	resources := loadResourcesToTest(t, "./testdata/resources/uses_variables.tf")
	assert.Equal(t, len(resources), 1, "Unexpected number of resources found")
	properties := resources[0].Properties.(map[string]interface{})
	assert.Equal(t, properties["instance_type"], "c4.large", "Unexpected value for instance_type")
	os.Setenv("TF_VAR_instance_type", "")
}

func TestTerraformFileFunction(t *testing.T) {
	resources := loadResourcesToTest(t, "./testdata/resources/reference_file.tf")
	assert.Equal(t, len(resources), 1, "Unexpected number of resources found")
	properties := resources[0].Properties.(map[string]interface{})
	assert.Equal(t, properties["bucket"], "example", "Unexpected value for bucket property")
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
	assert.Nil(t, err, "Expecting Validate to run without error")
	assert.Equal(t, len(report.ResourcesScanned), 1, "Unexpected number of resources")
	assert.Equal(t, len(report.FilesScanned), 2, "Unexpected number of files scanned")
	assertViolationsCount("TestTerraformVariablesInDifferentFile ", 0, report.Violations, t)
}

type TestingValueSource struct{}

func (s TestingValueSource) GetValue(a assertion.Expression) (string, error) {
	if a.ValueFrom.URL != "" {
		return "TEST", nil
	}
	return a.Value, nil
}

func TestTerraformDataLoader(t *testing.T) {
	loader := TerraformResourceLoader{}
	loaded, err := loader.Load("./testdata/resources/terraform_data.tf")
	assert.Nil(t, err, "Expecting Load to run without error")
	assert.Equal(t, len(loaded.Resources), 1, "Unexpected number of resources")
}

type terraformLinterTestCase struct {
	ConfigurationFilename   string
	RulesFilename           string
	ExpectedViolationCount  int
	ExpectedViolationRuleID string
}

func TestTerraformLinterCases(t *testing.T) {
	testCases := map[string]terraformLinterTestCase{
		"ParseError": {
			"./testdata/resources/terraform_syntax_error.tf",
			"./testdata/rules/terraform_provider.yml",
			1,
			"FILE_LOAD",
		},
		"Provider": {
			"./testdata/resources/terraform_provider.tf",
			"./testdata/rules/terraform_provider.yml",
			1,
			"AWS_PROVIDER",
		},
		"DataObject": {
			"./testdata/resources/terraform_data.tf",
			"./testdata/rules/terraform_data.yml",
			1,
			"DATA_NOT_CONTAINS",
		},
		"PoliciesWithVariables": {
			"./testdata/resources/policy_with_variables.tf",
			"./testdata/rules/policy_variable.yml",
			0,
			"",
		},
		"HereDocWithExpression": {
			"./testdata/resources/policy_with_expression.tf",
			"./testdata/rules/policy_variable.yml",
			0,
			"",
		},
		"Policies": {
			"./testdata/resources/terraform_policy.tf",
			"./testdata/rules/terraform_policy.yml",
			1,
			"TEST_POLICY",
		},
		"PolicyInvalidJSON": {
			"./testdata/resources/terraform_policy_invalid_json.tf",
			"./testdata/rules/terraform_policy.yml",
			0,
			"",
		},
		"PolicyEmpty": {
			"./testdata/resources/terraform_policy_empty.tf",
			"./testdata/rules/terraform_policy.yml",
			0,
			"",
		},
		"Module": {
			"./testdata/resources/terraform_module.tf",
			"./testdata/rules/terraform_module.yml",
			1,
			"MODULE_DESCRIPTION",
		},
		"BatchPrivileged": {
			"./testdata/resources/batch_privileged.tf",
			"./testdata/rules/batch_definition.yml",
			1,
			"BATCH_DEFINITION_PRIVILEGED",
		},
		"CloudfrontAccessLogs": {
			"./testdata/resources/cloudfront_access_logs.tf",
			"./testdata/rules/cloudfront_access_logs.yml",
			0,
			"",
		},
		"PublicEC2": {
			"./testdata/resources/ec2_public.tf",
			"./testdata/rules/ec2_public.yml",
			0,
			"",
		},
		"ElastiCacheRest": {
			"./testdata/resources/elasticache_encryption_rest.tf",
			"./testdata/rules/elasticache_encryption_rest.yml",
			1,
			"ELASTICACHE_ENCRYPTION_REST",
		},
		"ElastiCacheTransit": {
			"./testdata/resources/elasticache_encryption_transit.tf",
			"./testdata/rules/elasticache_encryption_transit.yml",
			1,
			"ELASTICACHE_ENCRYPTION_TRANSIT",
		},
		"NeptuneClusterEncryption": {
			"./testdata/resources/neptune_db_encryption.tf",
			"./testdata/rules/neptune_db_encryption.yml",
			1,
			"NEPTUNE_DB_ENCRYPTION",
		},
		"RdsPublic": {
			"./testdata/resources/rds_publicly_available.tf",
			"./testdata/rules/rds_publicly_available.yml",
			0,
			"",
		},
	}
	for name, tc := range testCases {
		options := Options{
			Tags:    []string{},
			RuleIDs: []string{},
		}
		filenames := []string{tc.ConfigurationFilename}
		linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: TerraformResourceLoader{}}
		ruleSet := loadRulesForTest(tc.RulesFilename, t)
		report, err := linter.Validate(ruleSet, options)
		if err != nil {
			t.Errorf("Expecting %s to return without an error: %s", name, err.Error())
		}
		if len(report.FilesScanned) != 1 {
			t.Errorf("TestTerraformLinterCases scanned %d files, expecting 1", len(report.FilesScanned))
		}
		if len(report.Violations) != tc.ExpectedViolationCount {
			t.Errorf("%s returned %d violations, expecting %d", name, len(report.Violations), tc.ExpectedViolationCount)
			t.Errorf("Violations: %v", report.Violations)
		}
		if tc.ExpectedViolationRuleID != "" {
			assertViolationByRuleID(name, tc.ExpectedViolationRuleID, report.Violations, t)
		}
	}
}
