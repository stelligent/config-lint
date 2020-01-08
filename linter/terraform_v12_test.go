package linter

import (
	"fmt"
	"github.com/stelligent/config-lint/assertion"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestTerraformV12Linter(t *testing.T) {
	options := Options{
		Tags:    []string{},
		RuleIDs: []string{},
	}
	filenames := []string{"./testdata/resources/terraform_instance.tf"}
	linter := FileLinter{Filenames: filenames, ValueSource: TestingValueSource{}, Loader: Terraform12ResourceLoader{}}
	ruleSet := loadRulesForTest("./testdata/rules/terraform_instance.yml", t)
	report, err := linter.Validate(ruleSet, options)
	assert.Nil(t, err, "Expecting Validate to run without error")
	assert.Equal(t, len(report.ResourcesScanned), 1, "Unexpected number of resources scanned")
	assert.Equal(t, len(report.FilesScanned), 1, "Unexpected number of files scanned")
	assertViolationsCount("TestTerraformLinter ", 0, report.Violations, t)
}

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

//The idea of this test is to confirm a particular difference between the original parser and the new
//I know it's not clear. - MN
func TestTupleType(t *testing.T) {
	resources := loadResources12ToTest(t, "./testdata/resources/multi_level.tf")
	assert.Equal(t, 1, len(resources), "Expecting 1 resource")
	statement := resources[0].Properties.(map[string]interface{})["statement"]
	principals := statement.([]interface{})[0].(map[string]interface{})["principals"]
	identifiers := principals.([]interface{})[0].(map[string]interface{})["identifiers"]
	value, ok := identifiers.([]interface{})
	assert.True(t, ok)
	_, ok = value[0].(string)
	assert.True(t, ok)
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

func TestTerraform12FunctionCall(t *testing.T) {
	resources := loadResources12ToTest(t, "./testdata/resources/uses_variables.tf")
	assert.Equal(t, len(resources), 1, "Expecting 1 resource")
	tags := getResourceTags(resources[0])
	assert.Equal(t, "test", tags["environment"], "Unexpected value for lookup function")
}

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

func TestTerraform12VariablesFromEnvironment(t *testing.T) {
	os.Setenv("TF_VAR_instance_type", "c4.large")
	resources := loadResources12ToTest(t, "./testdata/resources/uses_variables.tf")
	assert.Equal(t, len(resources), 1, "Unexpected number of resources found")
	properties := resources[0].Properties.(map[string]interface{})
	assert.Equal(t, properties["instance_type"], "c4.large", "Unexpected value for instance_type")
	os.Setenv("TF_VAR_instance_type", "")
}


func TestTerraform12FileFunction(t *testing.T) {
	resources := loadResources12ToTest(t, "./testdata/resources/reference_file.tf")
	assert.Equal(t, len(resources), 1, "Unexpected number of resources found")
	properties := resources[0].Properties.(map[string]interface{})
	assert.Equal(t, properties["bucket"], "example", "Unexpected value for bucket property")
}

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

func TestTerraform12ResourceLineNumber(t *testing.T) {
	resources := loadResources12ToTest(t, "./testdata/resources/uses_variables.tf")
	assert.Equal(t, 30, resources[0].LineNumber)
}

func TestTerraform12ResourceFileName(t *testing.T) {
	resources := loadResources12ToTest(t, "./testdata/resources/uses_variables.tf")
	assert.Equal(t, "./testdata/resources/uses_variables.tf", resources[0].Filename)
}

func TestTerraform12DataLineNumber(t *testing.T) {
	resources := loadResources12ToTest(t, "./testdata/resources/terraform_data.tf")
	assert.Equal(t, 1, resources[0].LineNumber)
}

func TestTerraform12DataFileName(t *testing.T) {
	resources := loadResources12ToTest(t, "./testdata/resources/terraform_data.tf")
	assert.Equal(t, "./testdata/resources/terraform_data.tf", resources[0].Filename)
}

func TestTerraform12ProviderLineNumber(t *testing.T) {
	resources := loadResources12ToTest(t, "./testdata/resources/terraform_provider.tf")
	assert.Equal(t, 1, resources[0].LineNumber)
}

func TestTerraform12ProviderFileName(t *testing.T) {
	resources := loadResources12ToTest(t, "./testdata/resources/terraform_provider.tf")
	assert.Equal(t, "./testdata/resources/terraform_provider.tf", resources[0].Filename)
}

func TestTerraform12ModuleLineNumber(t *testing.T) {
	resources := loadResources12ToTest(t, "./testdata/resources/terraform_module.tf")
	assert.Equal(t, 1, resources[0].LineNumber)
}

func TestTerraform12ModuleFileName(t *testing.T) {
	resources := loadResources12ToTest(t, "./testdata/resources/terraform_module.tf")
	assert.Equal(t, "./testdata/resources/terraform_module.tf", resources[0].Filename)
}

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
		"BatchPrivileged": {
			"./testdata/resources/batch_privileged.tf",
			"./testdata/rules/batch_definition.yml",
			1,
			"BATCH_DEFINITION_PRIVILEGED",
		},
		"PublicEC2": {
			"./testdata/resources/ec2_public.tf",
			"./testdata/rules/ec2_public.yml",
			0,
			"",
		},
		"CloudfrontAccessLogs": {
			"./testdata/resources/cloudfront_access_logs.tf",
			"./testdata/rules/cloudfront_access_logs.yml",
			0,
			"",
		},
		"Module": {
			"./testdata/resources/terraform_module.tf",
			"./testdata/rules/terraform_module.yml",
			1,
			"MODULE_DESCRIPTION",
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
		"KinesisKms": {
			"./testdata/resources/kinesis_kms_stream.tf",
			"./testdata/rules/kinesis_kms_stream.yml",
			1,
			"KINESIS_STREAM_KMS",
		},
		"DmsEncryption": {
			"./testdata/resources/dms_endpoint_encryption.tf",
			"./testdata/rules/dms_endpoint_encryption.yml",
			0,
			"",
		},
		"EmrClusterLogs": {
			"./testdata/resources/emr_cluster_logs.tf",
			"./testdata/rules/emr_cluster_logs.yml",
			1,
			"AWS_EMR_CLUSTER_LOGGING",
		},
		"KmsKeyRotation": {
			"./testdata/resources/kms_key_rotation.tf",
			"./testdata/rules/kms_key_rotation.yml",
			1,
			"AWS_KMS_KEY_ROTATION",
		},
		"SagemakerEndpoint": {
			"./testdata/resources/sagemaker_endpoint_encryption.tf",
			"./testdata/rules/sagemaker_endpoint_encryption.yml",
			1,
			"SAGEMAKER_ENDPOINT_ENCRYPTION",
		},
		"SagemakerNotebook": {
			"./testdata/resources/sagemaker_notebook_encryption.tf",
			"./testdata/rules/sagemaker_notebook_encryption.yml",
			1,
			"SAGEMAKER_NOTEBOOK_ENCRYPTION",
		},
		"TF12Variables": {
			"./testdata/resources/uses_tf12_variables.tf",
			"./testdata/rules/terraform_v12_variables.yml",
			0,
			"",
		},
		"TF12ForLoop": {
			"./testdata/resources/tf12_for_loop.tf",
			"./testdata/rules/tf12_for_loop.yml",
			0,
			"",
		},
		"TF12NullValue": {
			"./testdata/resources/nullable_value.tf",
			"./testdata/rules/nullable_value.yml",
			0,
			"",
		},
		//"TF12DynamicBlock": {
		//	"./testdata/resources/dynamic_block.tf",
		//	"./testdata/rules/dynamic_block.yml",
		//	1,
		//	"NO_SSH_ACCESS",
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

func TestTerraform12FileFunctionMultiLineContent(t *testing.T) {
	resources := loadResources12ToTest(t, "./testdata/resources/reference_file_multi_line.tf")
	assert.Equal(t, len(resources), 2, "Unexpected number of resources found")
	properties_1 := resources[0].Properties.(map[string]interface{})
	assert.Equal(t, properties_1["test_value"], "multi\nline\nexample", "Unexpected value for bucket property")
	properties_2 := resources[1].Properties.(map[string]interface{})
	assert.Equal(t, properties_2["test_value2"], properties_1["test_value"], "Unexpected value for bucket property")
}

func TestTerraform12FileFunctionResourceFileAbsolutePath(t *testing.T) {
	absolutePath, _ := filepath.Abs("./testdata/resources/reference_file.tf")
	resources := loadResources12ToTest(t, absolutePath)
	assert.Equal(t, len(resources), 1, "Unexpected number of resources found")
	properties := resources[0].Properties.(map[string]interface{})
	assert.Equal(t, properties["bucket"], "example", "Unexpected value for bucket property")
}

func TestTerraform12FileFunctionTemplateFileFunction(t *testing.T) {
	resources := loadResources12ToTest(t, "./testdata/resources/template_file_function_basic.tf")
	assert.Equal(t, len(resources), 1, "Unexpected number of resources found")
	properties := resources[0].Properties.(map[string]interface{})
	assert.Equal(t, properties["bucket"], "bucket-foo-example-bar", "Unexpected value for bucket property")
}

func TestTerraform12FileFunctionTemplateFileForLoop(t *testing.T) {
	resources := loadResources12ToTest(t, "./testdata/resources/template_file_function_for_loop.tf")
	assert.Equal(t, len(resources), 1, "Unexpected number of resources found")
	properties := resources[0].Properties.(map[string]interface{})
	assert.Equal(t, properties["test_value"], "testing:foo\ntesting:bar", "Unexpected value for bucket property")
}

func TestTerraform12FileFunctionTemplateFileConditional(t *testing.T) {
	resources := loadResources12ToTest(t, "./testdata/resources/template_file_function_conditional.tf")
	assert.Equal(t, len(resources), 2, "Unexpected number of resources found")
	properties := resources[0].Properties.(map[string]interface{})
	assert.Equal(t, properties["test_value"], "Foo", "Unexpected value for bucket property")
	properties2 := resources[1].Properties.(map[string]interface{})
	assert.Equal(t, properties2["test_value2"], "Bar", "Unexpected value for bucket property")
}

//func TestTerraform12FileFunctionReferenceAndResourceFileSameDir(t *testing.T) {
//	resources := loadResources12ToTest(t, "./testdata/data/reference_relative.tf")
//	assert.Equal(t, len(resources), 1, "Unexpected number of resources found")
//	properties := resources[0].Properties.(map[string]interface{})
//	assert.Equal(t, properties["bucket"], "example", "Unexpected value for bucket property")
//}

func TestTerraform12FileFunctionReferenceFileAbsoultePath(t *testing.T) {
	path, _ := os.Getwd()
	var err error
	var tempResourceFile *os.File
	var tempReferenceFile *os.File
	var tempResourceDir string
	var tempReferenceDir string

	tempResourceDir, err = ioutil.TempDir(path, "tf_resource")
	if err != nil{log.Fatal(err)}
	tempReferenceDir, err = ioutil.TempDir(tempResourceDir, "tf_reference")
	if err != nil{log.Fatal(err)}
	tempResourceFile, err = ioutil.TempFile(tempResourceDir, "test_resource.tf")
	if err != nil{log.Fatal(err)}
	tempReferenceFile, err = ioutil.TempFile(tempReferenceDir, "test_reference.txt")
	if err != nil{log.Fatal(err)}

	// tempReferenceFile.Name() is returned as the Absolute Path of the temp reference file
	tf12ResourceContent := fmt.Sprintf(`resource "aws_s3_bucket" "a_bucket" {
 bucket = "${file("%v")}"
}
`, tempReferenceFile.Name())

	tf12ReferenceContent := (`example
`)
	err = ioutil.WriteFile(tempResourceFile.Name(), []byte(tf12ResourceContent), 0644)
	if err != nil{log.Fatal(err)}
	err = ioutil.WriteFile(tempReferenceFile.Name(), []byte(tf12ReferenceContent), 0644)
	if err != nil{log.Fatal(err)}

	resources := loadResources12ToTest(t, tempResourceFile.Name())
	assert.Equal(t, len(resources), 1, "Unexpected number of resources found")
	properties := resources[0].Properties.(map[string]interface{})
	assert.Equal(t, properties["bucket"], "example", "Unexpected value for bucket property")
	os.RemoveAll(tempResourceDir)
	os.RemoveAll(tempReferenceDir)
}
