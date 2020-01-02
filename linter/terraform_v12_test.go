package linter

import (
	"github.com/stelligent/config-lint/assertion"
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestTerraformV12Linter(t *testing.T) {
//	options := Options{
//		Tags:    []string{},
//		RuleIDs: []string{},
//	}
//	filenames := []string{"./testdata/resources/terraform_instance.tf"}
//	linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: Terraform12ResourceLoader{}}
//	ruleSet := loadRulesForTest("./testdata/rules/terraform_instance.yml", t)
//	report, err := linter.Validate(ruleSet, options)
//	assert.Nil(t, err, "Expecting Validate to run without error")
//	assert.Equal(t, len(report.ResourcesScanned), 1, "Unexpected number of resources scanned")
//	assert.Equal(t, len(report.FilesScanned), 1, "Unexpected number of files scanned")
//	assertViolationsCount("TestTerraformLinter ", 0, report.Violations, t)
//}

func loadResources12ToTest(t *testing.T, filename string) []assertion.Resource {
	loader := Terraform12ResourceLoader{}
	loaded, err := loader.Load(filename)
	assert.Nil(t, err, "Expecting Load to run without error")
	resources, err := loader.PostLoad(loaded)
	assert.Nil(t, err, "Expecting PostLoad to run without error")
	return resources
}
//
//func getResourceTags(r assertion.Resource) map[string]interface{} {
//	properties := r.Properties.(map[string]interface{})
//	tags := properties["tags"].([]interface{})
//	return tags[0].(map[string]interface{})
//}

func TestSingleResourceType(t *testing.T) {
	resources := loadResources12ToTest(t, "./testdata/resources/uses_variables.tf")
	assert.Equal(t, 1, len(resources), "Expecting 1 resource")
	assert.Equal(t, "aws_instance", resources[0].Type)
	assert.Equal(t, "first", resources[0].ID)
}

func TestTerraform12Variable(t *testing.T) {
	loadResources12ToTest(t, "./testdata/resources/uses_variables.tf")
	resources := loadResources12ToTest(t, "./testdata/resources/uses_variables.tf")
	assert.Equal(t, 1, len(resources), "Expecting 1 resource")
	properties := resources[0].Properties.(map[string]interface{})
	assert.Equal(t, "ami-f2d3638a", properties["ami"], "Unexpected value for simple variable")
}

func TestTerraform12VariableWithNoDefault(t *testing.T) {
	resources := loadResources12ToTest(t, "./testdata/resources/uses_variables.tf")
	assert.Equal(t, len(resources), 1, "Expecting 1 resource")
	tags := getResourceTags(resources[0])
	assert.Equal(t, "", tags["department"], "Unexpected value for variable with no default")
}
//
//func TestTerraform12FunctionCall(t *testing.T) {
//	resources := loadResources12ToTest(t, "./testdata/resources/uses_variables.tf")
//	assert.Equal(t, len(resources), 1, "Expecting 1 resource")
//	tags := getResourceTags(resources[0])
//	assert.Equal(t, "test", tags["environment"], "Unexpected value for lookup function")
//}

func TestTerraform12ListVariable(t *testing.T) {
	resources := loadResources12ToTest(t, "./testdata/resources/uses_variables.tf")
	assert.Equal(t, len(resources), 1, "Expecting 1 resource")
	tags := getResourceTags(resources[0])
	assert.Equal(t, tags["comment"], "bar", "Unexpected value for list variable")
}

func TestTerraform12LocalVariable(t *testing.T) {
	resources := loadResources12ToTest(t, "./testdata/resources/uses_local_variables.tf")
	assert.Equal(t, len(resources), 1, "Expecting 1 resource")
	properties := resources[0].Properties.(map[string]interface{})
	assert.Equal(t, "myprojectbucket", properties["name"], "Unexpected value for name attribute")
}
//
//func TestTerraform12VariablesFromEnvironment(t *testing.T) {
//	os.Setenv("TF_VAR_instance_type", "c4.large")
//	resources := loadResources12ToTest(t, "./testdata/resources/uses_variables.tf")
//	assert.Equal(t, len(resources), 1, "Unexpected number of resources found")
//	properties := resources[0].Properties.(map[string]interface{})
//	assert.Equal(t, properties["instance_type"], "c4.large", "Unexpected value for instance_type")
//	os.Setenv("TF_VAR_instance_type", "")
//}
//
//func TestTerraform12FileFunction(t *testing.T) {
//	resources := loadResources12ToTest(t, "./testdata/resources/reference_file.tf")
//	assert.Equal(t, len(resources), 1, "Unexpected number of resources found")
//	properties := resources[0].Properties.(map[string]interface{})
//	assert.Equal(t, properties["bucket"], "example", "Unexpected value for bucket property")
//}

func TestTerraform12VariablesInDifferentFile(t *testing.T) {
	options := Options{
		Tags:    []string{},
		RuleIDs: []string{},
	}
	filenames := []string{
		"./testdata/resources/defines_variables.tf",
		"./testdata/resources/reference_variables.tf",
	}
	linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: Terraform12ResourceLoader{}}
	ruleSet := loadRulesForTest("./testdata/rules/terraform_instance.yml", t)
	report, err := linter.Validate(ruleSet, options)
	assert.Nil(t, err, "Expecting Validate to run without error")
	assert.Equal(t, 1, len(report.ResourcesScanned), "Unexpected number of resources")
	assert.Equal(t, 2, len(report.FilesScanned), "Unexpected number of files scanned")
	assertViolationsCount("TestTerraformVariablesInDifferentFile ", 0, report.Violations, t)
}

//type TestingValueSource struct{}

//func (s TestingValueSource) GetValue(a assertion.Expression) (string, error) {
//	if a.ValueFrom.URL != "" {
//		return "TEST", nil
//	}
//	return a.Value, nil
//}

func TestTerraform12DataLoader(t *testing.T) {
	loader := TerraformResourceLoader{}
	loaded, err := loader.Load("./testdata/resources/terraform_data.tf")
	assert.Nil(t, err, "Expecting Load to run without error")
	assert.Equal(t, len(loaded.Resources), 1, "Unexpected number of resources")
}

//type terraformLinterTestCase struct {
//	ConfigurationFilename   string
//	RulesFilename           string
//	ExpectedViolationCount  int
//	ExpectedViolationRuleID string
//}
//
func TestTerraform12LinterCases(t *testing.T) {
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
		//"DataObject": {
		//	"./testdata/resources/terraform_data.tf",
		//	"./testdata/rules/terraform_data.yml",
		//	1,
		//	"DATA_NOT_CONTAINS",
		//},
		//"PoliciesWithVariables": {
		//	"./testdata/resources/policy_with_variables.tf",
		//	"./testdata/rules/policy_variable.yml",
		//	0,
		//	"",
		//},
		//"HereDocWithExpression": {
		//	"./testdata/resources/policy_with_expression.tf",
		//	"./testdata/rules/policy_variable.yml",
		//	0,
		//	"",
		//},
		//"Policies": {
		//	"./testdata/resources/terraform_policy.tf",
		//	"./testdata/rules/terraform_policy.yml",
		//	1,
		//	"TEST_POLICY",
		//},
		//"PolicyInvalidJSON": {
		//	"./testdata/resources/terraform_policy_invalid_json.tf",
		//	"./testdata/rules/terraform_policy.yml",
		//	0,
		//	"",
		//},
		//"PolicyEmpty": {
		//	"./testdata/resources/terraform_policy_empty.tf",
		//	"./testdata/rules/terraform_policy.yml",
		//	0,
		//	"",
		//},
		//"Module": {
		//	"./testdata/resources/terraform_module.tf",
		//	"./testdata/rules/terraform_module.yml",
		//	1,
		//	"MODULE_DESCRIPTION",
		//},
		//"BatchPrivileged": {
		//	"./testdata/resources/batch_privileged.tf",
		//	"./testdata/rules/batch_definition.yml",
		//	1,
		//	"BATCH_DEFINITION_PRIVILEGED",
		//},
		//"CloudfrontAccessLogs": {
		//	"./testdata/resources/cloudfront_access_logs.tf",
		//	"./testdata/rules/cloudfront_access_logs.yml",
		//	0,
		//	"",
		//},
		//"PublicEC2": {
		//	"./testdata/resources/ec2_public.tf",
		//	"./testdata/rules/ec2_public.yml",
		//	0,
		//	"",
		//},
		//"ElastiCacheRest": {
		//	"./testdata/resources/elasticache_encryption_rest.tf",
		//	"./testdata/rules/elasticache_encryption_rest.yml",
		//	1,
		//	"ELASTICACHE_ENCRYPTION_REST",
		//},
		//"ElastiCacheTransit": {
		//	"./testdata/resources/elasticache_encryption_transit.tf",
		//	"./testdata/rules/elasticache_encryption_transit.yml",
		//	1,
		//	"ELASTICACHE_ENCRYPTION_TRANSIT",
		//},
		//"NeptuneClusterEncryption": {
		//	"./testdata/resources/neptune_db_encryption.tf",
		//	"./testdata/rules/neptune_db_encryption.yml",
		//	1,
		//	"NEPTUNE_DB_ENCRYPTION",
		//},
		//"RdsPublic": {
		//	"./testdata/resources/rds_publicly_available.tf",
		//	"./testdata/rules/rds_publicly_available.yml",
		//	0,
		//	"",
		//},
		//"KinesisKms": {
		//	"./testdata/resources/kinesis_kms_stream.tf",
		//	"./testdata/rules/kinesis_kms_stream.yml",
		//	1,
		//	"KINESIS_STREAM_KMS",
		//},
		//"DmsEncryption": {
		//	"./testdata/resources/dms_endpoint_encryption.tf",
		//	"./testdata/rules/dms_endpoint_encryption.yml",
		//	0,
		//	"",
		//},
		//"EmrClusterLogs": {
		//	"./testdata/resources/emr_cluster_logs.tf",
		//	"./testdata/rules/emr_cluster_logs.yml",
		//	1,
		//	"AWS_EMR_CLUSTER_LOGGING",
		//},
		//"KmsKeyRotation": {
		//	"./testdata/resources/kms_key_rotation.tf",
		//	"./testdata/rules/kms_key_rotation.yml",
		//	1,
		//	"AWS_KMS_KEY_ROTATION",
		//},
		//"SagemakerEndpoint": {
		//	"./testdata/resources/sagemaker_endpoint_encryption.tf",
		//	"./testdata/rules/sagemaker_endpoint_encryption.yml",
		//	1,
		//	"SAGEMAKER_ENDPOINT_ENCRYPTION",
		//},
		//"SagemakerNotebook": {
		//	"./testdata/resources/sagemaker_notebook_encryption.tf",
		//	"./testdata/rules/sagemaker_notebook_encryption.yml",
		//	1,
		//	"SAGEMAKER_NOTEBOOK_ENCRYPTION",
		//},
	}
	for name, tc := range testCases {
		options := Options{
			Tags:    []string{},
			RuleIDs: []string{},
		}
		filenames := []string{tc.ConfigurationFilename}
		linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: Terraform12ResourceLoader{}}
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
