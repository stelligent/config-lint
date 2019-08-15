package main

import (
	"fmt"
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
	Filename     string
	RuleID       string
	WarningCount int
	FailureCount int
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

func TestTerraformBuiltInRules(t *testing.T) {
	ruleSet := loadRules(t, "terraform.yml")
	testCases := []BuiltInTestCase{
		{"security-groups.tf", "SG_WORLD_INGRESS", 1, 0},
		{"security-groups.tf", "SG_WORLD_EGRESS", 2, 0},
		{"security-groups.tf", "SG_SSH_WORLD_INGRESS", 0, 1},
		{"security-groups.tf", "SG_RD_WORLD_INGRESS", 0, 0},
		{"security-groups.tf", "SG_NON_32_INGRESS", 2, 0},
		{"security-groups.tf", "SG_INGRESS_PORT_RANGE", 0, 0},
		{"security-groups.tf", "SG_EGRESS_PORT_RANGE", 0, 0},
		{"security-groups.tf", "SG_MISSING_EGRESS", 0, 0},
		{"security-groups.tf", "SG_INGRESS_ALL_PROTOCOLS", 1, 0},
		{"security-groups.tf", "SG_EGRESS_ALL_PROTOCOLS", 3, 0},
		{"cloudfront.tf", "CLOUDFRONT_DISTRIBUTION_LOGGING", 0, 1},
		{"cloudfront.tf", "CLOUDFRONT_DISTRIBUTION_ORIGIN_POLICY", 0, 0},
		{"cloudfront.tf", "CLOUDFRONT_DISTRIBUTION_DISTRIBUTION_PROTOCOl", 0, 0},
		{"iam.tf", "IAM_ROLE_POLICY_NOT_ACTION", 0, 0},
		{"iam.tf", "IAM_ROLE_POLICY_NOT_RESOURCE", 0, 0},
		{"iam.tf", "IAM_ROLE_POLICY_WILDCARD_ACTION", 0, 0},
		{"iam.tf", "IAM_ROLE_POLICY_WILDCARD_RESOURCE", 0, 0},
		{"iam.tf", "IAM_POLICY_NOT_ACTION", 0, 0},
		{"iam.tf", "IAM_POLICY_NOT_RESOURCE", 0, 0},
		{"iam.tf", "IAM_POLICY_WILDCARD_ACTION", 0, 1},
		{"iam.tf", "IAM_POLICY_WILDCARD_RESOURCE", 1, 0},
		{"iam.tf", "IAM_USER", 0, 0},
		{"iam.tf", "IAM_USER_POLICY_ATTACHMENT", 0, 1},
		{"iam.tf", "IAM_USER_GROUP", 0, 0},
		{"iam.tf", "POLICY_VERSION", 0, 1},
		{"iam.tf", "ASSUME_ROLEPOLICY_VERSION", 0, 1},
		{"elb.tf", "ELB_ACCESS_LOGGING", 1, 0},
		{"s3.tf", "S3_BUCKET_ACL", 0, 0},
		{"s3.tf", "S3_NOT_ACTION", 0, 0},
		{"s3.tf", "S3_NOT_PRINCIPAL", 0, 0},
		{"s3.tf", "S3_BUCKET_POLICY_WILDCARD_PRINCIPAL", 1, 0},
		{"s3.tf", "S3_BUCKET_POLICY_WILDCARD_ACTION", 1, 0},
		{"s3.tf", "S3_BUCKET_ENCRYPTION", 0, 1},
		{"s3.tf", "S3_BUCKET_OBJECT_ENCRYPTION", 0, 1},
		{"sns.tf", "SNS_TOPIC_POLICY_WILDCARD_PRINCIPAL", 1, 0},
		{"sns.tf", "SNS_TOPIC_POLICY_NOT_ACTION", 0, 0},
		{"sns.tf", "SNS_TOPIC_POLICY_NOT_PRINCIPAL", 0, 0},
		{"sqs.tf", "SQS_QUEUE_POLICY_WILDCARD_PRINCIPAL", 1, 0},
		{"sqs.tf", "SQS_QUEUE_POLICY_WILDCARD_ACTION", 0, 0},
		{"sqs.tf", "SQS_QUEUE_POLICY_NOT_ACTION", 0, 0},
		{"sqs.tf", "SQS_QUEUE_POLICY_NOT_PRINCIPAL", 0, 0},
		{"sqs.tf", "SQS_QUEUE_ENCRYPTION", 0, 1},
		{"lambda.tf", "LAMBDA_PERMISSION_INVOKE_ACTION", 0, 0},
		{"lambda.tf", "LAMBDA_PERMISSION_WILDCARD_PRINCIPAL", 0, 0},
		{"lambda.tf", "LAMBDA_FUNCTION_ENCRYPTION", 1, 0},
		{"lambda.tf", "LAMBDA_ENVIRONMENT_SECRETS", 0, 1},
		{"waf.tf", "WAF_WEB_ACL", 0, 0},
		{"alb.tf", "ALB_LISTENER_HTTPS", 0, 3},
		{"alb.tf", "ALB_ACCESS_LOGS", 0, 0},
		{"ami.tf", "AMI_VOLUMES_ENCRYPTED", 0, 1},
		{"ami.tf", "AMI_COPY_SNAPSHOTS_ENCRYPTED", 0, 1},
		{"ec2.tf", "EBS_BLOCK_DEVICE_ENCRYPTED", 0, 0},
		{"ec2.tf", "EBS_VOLUME_ENCRYPTION", 0, 2},
		{"cloudtrail.tf", "CLOUDTRAIL_ENCRYPTION", 0, 1},
		{"codebuild.tf", "CODEBUILD_PROJECT_ENCRYPTION", 0, 1},
		{"codepipeline.tf", "CODEPIPELINE_ENCRYPTION", 0, 0},
		{"db.tf", "DB_INSTANCE_ENCRYPTION", 0, 1},
		{"db.tf", "RDS_CLUSTER_ENCYPTION", 0, 2},
		{"efs.tf", "EFS_ENCRYPTED", 0, 1},
		{"kinesis.tf", "KINESIS_FIREHOSE_DELIVERY_STREAM_ENCRYPTION", 0, 1},
		{"redshift.tf", "REDSHIFT_CLUSTER_ENCRYPTION", 0, 1},
		{"ecs.tf", "ECS_ENVIRONMENT_SECRETS", 0, 1},
	}
	for _, tc := range testCases {

		filenames := []string{"testdata/builtin/terraform/" + tc.Filename}
		options := linter.Options{
			RuleIDs: []string{tc.RuleID},
		}
		vs := assertion.StandardValueSource{}
		l, err := linter.NewLinter(ruleSet, vs, filenames)
		report, err := l.Validate(ruleSet, options)
		assert.Nil(t, err, "Validate failed for file")
		warningMessage := fmt.Sprintf("Expecting %d warnings for RuleID %s in File %s", tc.WarningCount, tc.RuleID, tc.Filename)
		assert.Equal(t, tc.WarningCount, numberOfWarnings(report.Violations), warningMessage)
		failureMessage := fmt.Sprintf("Expecting %d failures for RuleID %s in File %s", tc.FailureCount, tc.RuleID, tc.Filename)
		assert.Equal(t, tc.FailureCount, numberOfFailures(report.Violations), failureMessage)
	}
}
