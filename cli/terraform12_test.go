package main

import (
	"fmt"
	"testing"

	"github.com/stelligent/config-lint/assertion"
	"github.com/stelligent/config-lint/linter"
	"github.com/stretchr/testify/assert"
)

// Run built in rules against Terraform v0.12 parser
func TestTerraform12BuiltInRules(t *testing.T) {

	// Define file to load rules from
	// This file is located under cli/assets/
	ruleSet := loadRules(t, "terraform.yml")

	// Get list of test cases
	testCases := []BuiltInTestCase{
		// AWS
		{"aws/alb/access_logs_enabled.tf", "ALB_ACCESS_LOGS", 0, 3},
		{"aws/alb_listener/https.tf", "ALB_LISTENER_HTTPS", 0, 4},
		{"aws/alb_listener/ssl_policy.tf", "ALB_LISTENER_SSL_POLICY", 0, 6},
		{"aws/ami/ebs_block_device_encrypted.tf", "AMI_VOLUMES_ENCRYPTED", 0, 2},
		{"aws/ami_copy/encrypted.tf", "AMI_COPY_SNAPSHOTS_ENCRYPTED", 0, 2},
		{"aws/batch_job_definition/container_properties_privileged.tf", "BATCH_DEFINITION_PRIVILEGED", 1, 0},
		{"aws/cloudfront_distribution/custom_origin_config.tf", "CLOUDFRONT_DISTRIBUTION_ORIGIN_POLICY", 0, 2},
		{"aws/cloudfront_distribution/logging_config.tf", "CLOUDFRONT_DISTRIBUTION_LOGGING", 0, 1},
		{"aws/cloudfront_distribution/viewer_protocol_policy.tf", "CLOUDFRONT_DISTRIBUTION_PROTOCOL", 0, 2},
		{"aws/cloudtrail/kms_key_id.tf", "CLOUDTRAIL_ENCRYPTION", 1, 0},
		{"aws/codebuild_project/artifact_encryption.tf", "CODEBUILD_PROJECT_ARTIFACT_ENCRYPTION", 0, 3},
		{"aws/codebuild_project/project_encryption.tf", "CODEBUILD_PROJECT_ENCRYPTION", 0, 1},
		{"aws/codepipeline/encryption_key.tf", "CODEPIPELINE_ENCRYPTION", 1, 0},
		{"aws/db_instance/publicly_accessible.tf", "RDS_PUBLIC_AVAILABILITY", 0, 1},
		{"aws/db_instance/storage_encryption.tf", "DB_INSTANCE_ENCRYPTION", 0, 2},
		{"aws/db_instance/storage_encryption.tf", "REPLICA_DB_INSTANCE_ENCRYPTION", 2, 0},
		{"aws/dms_endpoint/kms_key.tf", "AWS_DMS_ENDPOINT_ENCRYPTION", 1, 0},
		{"aws/ebs_volume/encrypted.tf", "EBS_VOLUME_ENCRYPTION", 0, 2},
		{"aws/ecs_task_definition/secrets.tf", "ECS_ENVIRONMENT_SECRETS", 0, 3},
		{"aws/efs_file_system/encrypted.tf", "EFS_ENCRYPTED", 0, 2},
		{"aws/elasticache_replication_group/encryption_at_rest.tf", "ELASTICACHE_ENCRYPTION_REST", 0, 2},
		{"aws/elasticache_replication_group/encryption_in_transit.tf", "ELASTICACHE_ENCRYPTION_TRANSIT", 0, 2},
		{"aws/elb/access_logs_enabled.tf", "ELB_ACCESS_LOGGING", 2, 0},
		{"aws/emr_cluster/logging.tf", "AWS_EMR_CLUSTER_LOGGING", 1, 0},
		{"aws/iam_group_membership/group_and_users.tf", "IAM_USER_GROUP", 0, 4},
		{"aws/iam_policy/policy_action_wildcard.tf", "IAM_POLICY_WILDCARD_ACTION", 0, 1},
		{"aws/iam_policy/policy_notaction.tf", "IAM_POLICY_NOT_ACTION", 1, 0},
		{"aws/iam_policy/policy_notresource.tf", "IAM_POLICY_NOT_RESOURCE", 1, 0},
		{"aws/iam_policy/policy_resource_wildcard.tf", "IAM_POLICY_WILDCARD_RESOURCE", 0, 1},
		{"aws/iam_policy/policy_version.tf", "POLICY_VERSION", 0, 1},
		{"aws/iam_role/assume_role_policy_action_wildcard.tf", "IAM_ROLE_WILDCARD_ACTION", 0, 1},
		{"aws/iam_role/assume_role_policy_notaction.tf", "IAM_ROLE_NOT_ACTION", 1, 0},
		{"aws/iam_role/assume_role_policy_notprincipal.tf", "IAM_ROLE_NOT_PRINCIPAL", 1, 0},
		{"aws/iam_role/assume_role_policy_version.tf", "ASSUME_ROLEPOLICY_VERSION", 0, 1},
		{"aws/iam_role_policy/policy_action_wildcard.tf", "IAM_ROLE_POLICY_WILDCARD_ACTION", 0, 1},
		{"aws/iam_role_policy/policy_notaction.tf", "IAM_ROLE_POLICY_NOT_ACTION", 1, 0},
		{"aws/iam_role_policy/policy_notresource.tf", "IAM_ROLE_POLICY_NOT_RESOURCE", 1, 0},
		{"aws/iam_role_policy/policy_resource_wildcard.tf", "IAM_ROLE_POLICY_WILDCARD_RESOURCE", 0, 1},
		{"aws/iam_role_policy/policy_version.tf", "POLICY_VERSION", 0, 1},
		{"aws/iam_user_policy/resource_exists.tf", "IAM_USER_POLICY", 0, 1},
		{"aws/iam_user_policy_attachment/resource_exists.tf", "IAM_USER_POLICY_ATTACHMENT", 0, 1},
		{"aws/instance/ebs_block_device_encrypted.tf", "EBS_BLOCK_DEVICE_ENCRYPTED", 0, 2},
		{"aws/kinesis_firehose_delivery_stream/encryption.tf", "KINESIS_FIREHOSE_DELIVERY_STREAM_ENCRYPTION", 0, 4},
		{"aws/kinesis_stream/encryption.tf", "KINESIS_STREAM_ENCRYPTION", 0, 2},
		{"aws/kinesis_stream/kms_key.tf", "KINESIS_STREAM_KMS", 1, 0},
		{"aws/kms_key/rotation.tf", "AWS_KMS_KEY_ROTATION", 2, 0},
		{"aws/lambda_function/encryption.tf", "LAMBDA_FUNCTION_ENCRYPTION", 1, 0},
		{"aws/lambda_function/environment_variables_aws_secrets.tf", "LAMBDA_ENVIRONMENT_SECRETS", 0, 3},
		{"aws/lambda_permission/action.tf", "LAMBDA_PERMISSION_INVOKE_ACTION", 1, 0},
		{"aws/lambda_permission/principal_wildcard.tf", "LAMBDA_PERMISSION_WILDCARD_PRINCIPAL", 0, 2},
		{"aws/lb/access_logs_enabled.tf", "ALB_ACCESS_LOGS", 0, 3},
		{"aws/lb_listener/https.tf", "ALB_LISTENER_HTTPS", 0, 4},
		{"aws/lb_listener/ssl_policy.tf", "ALB_LISTENER_SSL_POLICY", 0, 6},
		{"aws/neptune_cluster/encryption.tf", "NEPTUNE_DB_ENCRYPTION", 0, 2},
		{"aws/rds_cluster/storage_encryption.tf", "RDS_CLUSTER_ENCYPTION", 0, 5},
		{"aws/redshift_cluster/encrypted.tf", "REDSHIFT_CLUSTER_ENCRYPTION", 0, 2},
		{"aws/redshift_cluster/enhanced_vpc_routing.tf", "REDSHIFT_CLUSTER_ENHANCED_VPC_ROUTING", 2, 0},
		{"aws/redshift_cluster/kms_key_id.tf", "REDSHIFT_CLUSTER_KMS_KEY_ID", 1, 0},
		{"aws/redshift_cluster/logging.tf", "REDSHIFT_CLUSTER_AUDIT_LOGGING", 2, 0},
		{"aws/redshift_cluster/publicly_accessible.tf", "REDSHIFT_CLUSTER_PUBLICLY_ACCESSIBLE", 0, 2},
		{"aws/redshift_parameter_group/require_ssl.tf", "REDSHIFT_CLUSTER_PARAMETER_GROUP_REQUIRE_SSL", 2, 0},
		{"aws/s3_bucket/acl_not_public.tf", "S3_BUCKET_ACL", 0, 2},
		{"aws/s3_bucket/server_side_encryption_enabled.tf", "S3_BUCKET_ENCRYPTION", 0, 1},
		{"aws/s3_bucket_object/encryption_enabled.tf", "S3_BUCKET_OBJECT_ENCRYPTION", 0, 1},
		{"aws/s3_bucket_policy/policy_statement_action_wildcard.tf", "S3_BUCKET_POLICY_WILDCARD_ACTION", 0, 1},
		{"aws/s3_bucket_policy/policy_statement_notaction.tf", "S3_NOT_ACTION", 1, 0},
		{"aws/s3_bucket_policy/policy_statement_notprincipal.tf", "S3_NOT_PRINCIPAL", 1, 0},
		{"aws/s3_bucket_policy/policy_statement_principal_wildcard.tf", "S3_BUCKET_POLICY_WILDCARD_PRINCIPAL", 0, 1},
		{"aws/s3_bucket_policy/policy_statement_secure_transport.tf", "S3_BUCKET_POLICY_ONLY_HTTPS", 0, 1},
		{"aws/s3_bucket_policy/policy_version.tf", "POLICY_VERSION", 0, 1},
		{"aws/sagemaker_endpoint_configuration/kms_key.tf", "SAGEMAKER_ENDPOINT_ENCRYPTION", 1, 0},
		{"aws/sagemaker_notebook_instance/kms_key.tf", "SAGEMAKER_NOTEBOOK_ENCRYPTION", 1, 0},
		{"aws/security_group/egress_all_protocols.tf", "SG_EGRESS_ALL_PROTOCOLS", 1, 0},
		{"aws/security_group/egress_port_range.tf", "SG_EGRESS_PORT_RANGE", 1, 0},
		{"aws/security_group/ingress_all_protocols.tf", "SG_INGRESS_ALL_PROTOCOLS", 1, 0},
		{"aws/security_group/ingress_port_range.tf", "SG_INGRESS_PORT_RANGE", 1, 0},
		{"aws/security_group/missing_egress.tf", "SG_MISSING_EGRESS", 1, 0},
		{"aws/security_group/non_32_ingress.tf", "SG_NON_32_INGRESS", 2, 0},
		{"aws/security_group/rdp_world_ingress.tf", "SG_RDP_WORLD_INGRESS", 0, 2},
		{"aws/security_group/ssh_world_ingress.tf", "SG_SSH_WORLD_INGRESS", 0, 2},
		{"aws/security_group/world_egress.tf", "SG_WORLD_EGRESS", 2, 0},
		{"aws/security_group/world_ingress.tf", "SG_WORLD_INGRESS", 2, 0},
		{"aws/sns_topic_policy/policy_statement_notaction.tf", "SNS_TOPIC_POLICY_NOT_ACTION", 1, 0},
		{"aws/sns_topic_policy/policy_statement_notprincipal.tf", "SNS_TOPIC_POLICY_NOT_PRINCIPAL", 1, 0},
		{"aws/sns_topic_policy/policy_statement_principal_wildcard-copy.tf", "SNS_TOPIC_POLICY_WILDCARD_PRINCIPAL", 0, 1},
		{"aws/sns_topic_policy/policy_version.tf", "POLICY_VERSION", 0, 1},
		{"aws/sqs_queue/encryption.tf", "SQS_QUEUE_ENCRYPTION", 0, 1},
		{"aws/sqs_queue_policy/policy_statement_action_wildcard.tf", "SQS_QUEUE_POLICY_WILDCARD_ACTION", 0, 1},
		// {"aws/sqs_queue_policy/policy_statement_notaction.tf", "SQS_QUEUE_POLICY_NOT_ACTION", 1, 0},
		// {"aws/sqs_queue_policy/policy_statement_notprincipal.tf", "SQS_QUEUE_POLICY_NOT_PRINCIPAL", 1, 0},
		// {"aws/sqs_queue_policy/policy_statement_principal_wildcard.tf", "SQS_QUEUE_POLICY_WILDCARD_PRINCIPAL", 0, 1},
		// {"aws/sqs_queue_policy/policy_version.tf", "POLICY_VERSION", 0, 1},
		{"aws/subnet/map_public_ip_on_launch.tf", "EC2_SUBNET_MAP_PUBLIC", 1, 0},
		{"aws/waf_web_acl/default_action_type.tf", "WAF_WEB_ACL", 0, 1},
	}

	// Run test cases
	// test files must be included under testdata/builtin/terraform12
	for _, tc := range testCases {
		filenames := []string{"testdata/builtin/terraform12/" + tc.Filename}
		options := linter.Options{
			RuleIDs: []string{tc.RuleID},
		}
		vs := assertion.StandardValueSource{}

		// Defining 'tf12' for the Parser type
		l, err := linter.NewLinter(ruleSet, vs, filenames, "tf12")
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
