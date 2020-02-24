package main

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stelligent/config-lint/assertion"
	"github.com/stelligent/config-lint/linter"
	"github.com/stretchr/testify/assert"
)

func loadRules(t *testing.T, filename string) assertion.RuleSet {
	r, err := loadBuiltInRuleSet(filename)
	assert.Nil(t, err, "Cannot load ruleset file")
	return r
}

type BuiltInTestCase struct {
	TerraformVersion string // "tf", "tf12", or "both"
	Filename         string
	RuleID           string
	WarningCount     int
	FailureCount     int
}

func numberOfWarnings(violations []assertion.Violation) int {
	n := 0
	for _, v := range violations {
		if v.Status == "WARNING" {
			n += 1
		}
	}
	return n
}
func numberOfFailures(violations []assertion.Violation) int {
	n := 0
	for _, v := range violations {
		if v.Status == "FAILURE" {
			n += 1
		}
	}
	return n
}

// String build message for violations. Debug helper
func getViolationsString(severity string, violations []assertion.Violation) string {
	var violationsReported string
	for count, v := range violations {
		if v.Status == severity {
			violationsReported += strconv.Itoa(count+1) + ". Violation:"
			violationsReported += "\n\tRule Message: " + v.RuleMessage
			violationsReported += "\n\tRule Id: " + v.RuleID
			violationsReported += "\n\tResource ID: " + v.ResourceID
			violationsReported += "\n\tResource Type: " + v.ResourceType
			violationsReported += "\n\tCategory: " + v.Category
			violationsReported += "\n\tStatus: " + v.Status
			violationsReported += "\n\tAssertion Message: " + v.AssertionMessage
			violationsReported += "\n\tFilename: " + v.Filename
			violationsReported += "\n\tLine Number: " + strconv.Itoa(v.LineNumber)
			violationsReported += "\n\tCreated At: " + v.CreatedAt + "\n"
		}
	}
	return violationsReported
}

func TestTerraformBuiltInRules(t *testing.T) {
	RunTestTerraformBuiltInRules(t, "tf")
}

func TestTerraform12BuiltInRules(t *testing.T) {
	RunTestTerraformBuiltInRules(t, "tf12")
}

// Run built in rules against Terraform
func RunTestTerraformBuiltInRules(t *testing.T, terraformVersion string) {

	// Define a version of terraform to run against.
	// Defaults to terraform 11
	// These files are under cli/assets/
	ruleSet := loadRules(t, "terraform.yml")
	if terraformVersion == "tf12" {
		ruleSet = loadRules(t, "terraform12.yml")
	}

	// Get list of test cases
	// Test cases can be applied to tf11, tf12, or both
	testCases := []BuiltInTestCase{
		// AWS
		{"tf", "aws/security_group/world_ingress.tf", "SG_WORLD_INGRESS", 2, 0},
		{"tf", "aws/security_group/world_egress.tf", "SG_WORLD_EGRESS", 2, 0},
		{"tf", "aws/security_group/ssh_world_ingress.tf", "SG_SSH_WORLD_INGRESS", 0, 2},
		{"tf", "aws/security_group/rdp_world_ingress.tf", "SG_RDP_WORLD_INGRESS", 0, 2},
		{"tf", "aws/security_group/non_32_ingress.tf", "SG_NON_32_INGRESS", 2, 0},
		{"both", "aws/security_group/ingress_port_range.tf", "SG_INGRESS_PORT_RANGE", 1, 0},
		{"both", "aws/security_group/egress_port_range.tf", "SG_EGRESS_PORT_RANGE", 1, 0},
		{"both", "aws/security_group/missing_egress.tf", "SG_MISSING_EGRESS", 1, 0},
		{"both", "aws/security_group/ingress_all_protocols.tf", "SG_INGRESS_ALL_PROTOCOLS", 1, 0},
		{"both", "aws/security_group/egress_all_protocols.tf", "SG_EGRESS_ALL_PROTOCOLS", 1, 0},
		{"both", "aws/cloudfront_distribution/logging_config.tf", "CLOUDFRONT_DISTRIBUTION_LOGGING", 0, 1},
		{"tf", "aws/cloudfront_distribution/custom_origin_config.tf", "CLOUDFRONT_DISTRIBUTION_ORIGIN_POLICY", 0, 2},
		{"tf", "aws/cloudfront_distribution/viewer_protocol_policy.tf", "CLOUDFRONT_DISTRIBUTION_PROTOCOL", 0, 2},
		{"both", "aws/iam_role/assume_role_policy_notaction.tf", "IAM_ROLE_NOT_ACTION", 1, 0},
		{"tf", "aws/iam_role/assume_role_policy_notprincipal.tf", "IAM_ROLE_NOT_PRINCIPAL", 1, 0},
		{"both", "aws/iam_role/assume_role_policy_action_wildcard.tf", "IAM_ROLE_WILDCARD_ACTION", 0, 1},
		{"both", "aws/iam_role_policy/policy_notaction.tf", "IAM_ROLE_POLICY_NOT_ACTION", 1, 0},
		{"both", "aws/iam_role_policy/policy_notresource.tf", "IAM_ROLE_POLICY_NOT_RESOURCE", 1, 0},
		{"both", "aws/iam_role_policy/policy_action_wildcard.tf", "IAM_ROLE_POLICY_WILDCARD_ACTION", 0, 1},
		{"tf", "aws/iam_role_policy/policy_resource_wildcard.tf", "IAM_ROLE_POLICY_WILDCARD_RESOURCE", 0, 1},
		{"both", "aws/iam_policy/policy_notaction.tf", "IAM_POLICY_NOT_ACTION", 1, 0},
		{"both", "aws/iam_policy/policy_notresource.tf", "IAM_POLICY_NOT_RESOURCE", 1, 0},
		{"both", "aws/iam_policy/policy_action_wildcard.tf", "IAM_POLICY_WILDCARD_ACTION", 0, 1},
		{"tf", "aws/iam_policy/policy_resource_wildcard.tf", "IAM_POLICY_WILDCARD_RESOURCE", 0, 1},
		{"tf", "aws/iam_user_policy/resource_exists.tf", "IAM_USER_POLICY", 0, 1},
		{"both", "aws/iam_user_policy_attachment/resource_exists.tf", "IAM_USER_POLICY_ATTACHMENT", 0, 1},
		{"both", "aws/iam_group_membership/group_and_users.tf", "IAM_USER_GROUP", 0, 4},
		{"tf", "aws/iam_policy/policy_version.tf", "POLICY_VERSION", 0, 1},
		{"tf", "aws/iam_role_policy/policy_version.tf", "POLICY_VERSION", 0, 1},
		{"both", "aws/iam_role/assume_role_policy_version.tf", "ASSUME_ROLEPOLICY_VERSION", 0, 1},
		{"tf", "aws/elb/access_logs_enabled.tf", "ELB_ACCESS_LOGGING", 2, 0},
		{"both", "aws/s3_bucket/acl_not_public.tf", "S3_BUCKET_ACL", 0, 2},
		{"both", "aws/s3_bucket_policy/policy_statement_notaction.tf", "S3_NOT_ACTION", 1, 0},
		{"both", "aws/s3_bucket_policy/policy_statement_notprincipal.tf", "S3_NOT_PRINCIPAL", 1, 0},
		{"tf", "aws/s3_bucket_policy/policy_version.tf", "POLICY_VERSION", 0, 1},
		{"tf", "aws/s3_bucket_policy/policy_statement_principal_wildcard.tf", "S3_BUCKET_POLICY_WILDCARD_PRINCIPAL", 0, 1},
		{"tf", "aws/s3_bucket_policy/policy_statement_action_wildcard.tf", "S3_BUCKET_POLICY_WILDCARD_ACTION", 0, 1},
		{"tf", "aws/s3_bucket_policy/policy_statement_secure_transport.tf", "S3_BUCKET_POLICY_ONLY_HTTPS", 0, 1},
		{"both", "aws/s3_bucket/server_side_encryption_enabled.tf", "S3_BUCKET_ENCRYPTION", 0, 1},
		{"tf", "aws/s3_bucket_object/encryption_enabled.tf", "S3_BUCKET_OBJECT_ENCRYPTION", 0, 1},
		{"tf", "aws/sns_topic_policy/policy_version.tf", "POLICY_VERSION", 0, 1},
		{"tf", "aws/sns_topic_policy/policy_statement_principal_wildcard-copy.tf", "SNS_TOPIC_POLICY_WILDCARD_PRINCIPAL", 0, 1},
		{"both", "aws/sns_topic_policy/policy_statement_notaction.tf", "SNS_TOPIC_POLICY_NOT_ACTION", 1, 0},
		{"both", "aws/sns_topic_policy/policy_statement_notprincipal.tf", "SNS_TOPIC_POLICY_NOT_PRINCIPAL", 1, 0},
		{"tf", "aws/sqs_queue_policy/policy_version.tf", "POLICY_VERSION", 0, 1},
		{"tf", "aws/sqs_queue_policy/policy_statement_principal_wildcard.tf", "SQS_QUEUE_POLICY_WILDCARD_PRINCIPAL", 0, 1},
		{"tf", "aws/sqs_queue_policy/policy_statement_action_wildcard.tf", "SQS_QUEUE_POLICY_WILDCARD_ACTION", 0, 1},
		{"tf", "aws/sqs_queue_policy/policy_statement_notaction.tf", "SQS_QUEUE_POLICY_NOT_ACTION", 1, 0},
		{"tf", "aws/sqs_queue_policy/policy_statement_notprincipal.tf", "SQS_QUEUE_POLICY_NOT_PRINCIPAL", 1, 0},
		{"tf", "aws/sqs_queue/encryption.tf", "SQS_QUEUE_ENCRYPTION", 0, 1},
		{"both", "aws/lambda_permission/action.tf", "LAMBDA_PERMISSION_INVOKE_ACTION", 1, 0},
		{"tf", "aws/lambda_permission/principal_wildcard.tf", "LAMBDA_PERMISSION_WILDCARD_PRINCIPAL", 0, 2},
		{"tf", "aws/lambda_function/encryption.tf", "LAMBDA_FUNCTION_ENCRYPTION", 1, 0},
		{"tf", "aws/lambda_function/environment_variables_aws_secrets.tf", "LAMBDA_ENVIRONMENT_SECRETS", 0, 3},
		{"both", "aws/waf_web_acl/default_action_type.tf", "WAF_WEB_ACL", 0, 1},
		{"both", "aws/alb_listener/https.tf", "ALB_LISTENER_HTTPS", 0, 4},
		{"tf", "aws/alb_listener/ssl_policy.tf", "ALB_LISTENER_SSL_POLICY", 0, 6},
		{"tf", "aws/alb/access_logs_enabled.tf", "ALB_ACCESS_LOGS", 0, 3},
		{"both", "aws/lb_listener/https.tf", "ALB_LISTENER_HTTPS", 0, 4},
		{"tf", "aws/lb_listener/ssl_policy.tf", "ALB_LISTENER_SSL_POLICY", 0, 6},
		{"tf", "aws/lb/access_logs_enabled.tf", "ALB_ACCESS_LOGS", 0, 3},
		{"tf", "aws/ami/ebs_block_device_encrypted.tf", "AMI_VOLUMES_ENCRYPTED", 0, 2},
		{"tf", "aws/ami_copy/encrypted.tf", "AMI_COPY_SNAPSHOTS_ENCRYPTED", 0, 2},
		{"tf", "aws/instance/ebs_block_device_encrypted.tf", "EBS_BLOCK_DEVICE_ENCRYPTED", 0, 2},
		{"tf", "aws/ebs_volume/encrypted.tf", "EBS_VOLUME_ENCRYPTION", 0, 2},
		{"tf", "aws/subnet/map_public_ip_on_launch.tf", "EC2_SUBNET_MAP_PUBLIC", 1, 0},
		{"tf", "aws/cloudtrail/kms_key_id.tf", "CLOUDTRAIL_ENCRYPTION", 1, 0},
		{"both", "aws/codebuild_project/project_encryption.tf", "CODEBUILD_PROJECT_ENCRYPTION", 0, 1},
		{"both", "aws/codebuild_project/artifact_encryption.tf", "CODEBUILD_PROJECT_ARTIFACT_ENCRYPTION", 0, 3},
		{"tf", "aws/codepipeline/encryption_key.tf", "CODEPIPELINE_ENCRYPTION", 1, 0},
		{"tf", "aws/db_instance/storage_encryption.tf", "DB_INSTANCE_ENCRYPTION", 0, 2},
		{"tf", "aws/db_instance/storage_encryption.tf", "REPLICA_DB_INSTANCE_ENCRYPTION", 1, 0},
		{"tf", "aws/db_instance/publicly_accessible.tf", "RDS_PUBLIC_AVAILABILITY", 0, 1},
		{"tf", "aws/rds_cluster/storage_encryption.tf", "RDS_CLUSTER_ENCYPTION", 0, 5},
		{"both", "aws/efs.tf", "EFS_ENCRYPTED", 0, 1},
		{"both", "aws/kinesis.tf", "KINESIS_FIREHOSE_DELIVERY_STREAM_ENCRYPTION", 0, 1},
		{"tf", "aws/redshift/cluster/encrypted.tf", "REDSHIFT_CLUSTER_ENCRYPTION", 0, 2},
		{"tf", "aws/redshift/cluster/enhanced_vpc_routing.tf", "REDSHIFT_CLUSTER_ENHANCED_VPC_ROUTING", 2, 0},
		{"tf", "aws/redshift/cluster/kms_key_id.tf", "REDSHIFT_CLUSTER_KMS_KEY_ID", 1, 0},
		{"tf", "aws/redshift/cluster/logging.tf", "REDSHIFT_CLUSTER_AUDIT_LOGGING", 2, 0},
		{"tf", "aws/redshift/cluster/publicly_accessible.tf", "REDSHIFT_CLUSTER_PUBLICLY_ACCESSIBLE", 0, 2},
		{"tf", "aws/redshift/parameter_group/require_ssl.tf", "REDSHIFT_CLUSTER_PARAMETER_GROUP_REQUIRE_SSL", 2, 0},
		{"both", "aws/ecs.tf", "ECS_ENVIRONMENT_SECRETS", 0, 1},
	}

	// Run test cases
	// test files must be included under testdata/builtin/terraform
	for _, tc := range testCases {
		if tc.TerraformVersion == "both" || tc.TerraformVersion == terraformVersion {
			filenames := []string{"testdata/builtin/terraform/" + tc.Filename}
			options := linter.Options{
				RuleIDs: []string{tc.RuleID},
			}
			vs := assertion.StandardValueSource{}
			l, err := linter.NewLinter(ruleSet, vs, filenames, "")
			report, err := l.Validate(ruleSet, options)
			assert.Nil(t, err, "Validate failed for file")

			warningViolationsReported := getViolationsString("WARNING", report.Violations)
			warningMessage := fmt.Sprintf("Expecting %d warnings for RuleID %s in File %s:\n %s", tc.WarningCount, tc.RuleID, tc.Filename, warningViolationsReported)
			assert.Equal(t, tc.WarningCount, numberOfWarnings(report.Violations), warningMessage)
			failureViolationsReported := getViolationsString("FAILURE", report.Violations)
			failureMessage := fmt.Sprintf("Expecting %d failures for RuleID %s in File %s:\n %s", tc.FailureCount, tc.RuleID, tc.Filename, failureViolationsReported)
			assert.Equal(t, tc.FailureCount, numberOfFailures(report.Violations), failureMessage)
		}
	}
}
