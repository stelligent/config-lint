# Test that sqs queue has kms_master_key_id and kms_data_key_reuse_period_seconds defined
# https://www.terraform.io/docs/providers/aws/r/sqs_queue.html#kms_master_key_id
# https://www.terraform.io/docs/providers/aws/r/sqs_queue.html#kms_data_key_reuse_period_seconds

provider "aws" {
  region = "us-east-1"
}

# Pass: both kms_master_key_id and kms_data_key_reuse_period_seconds are used
resource "aws_sqs_queue" "encryption_enabled_with_reuse_set" {
  name                              = "queue-with-encryption"
  kms_master_key_id                 = "alias/foo/bar"
  kms_data_key_reuse_period_seconds = 60
}

# Fail: neither kms_master_key_id or kms_data_key_reuse_period_seconds are used
resource "aws_sqs_queue" "encryption_disabled" {
  name = "queue-without-encryption"
}

# Fail: kms key defined without kms_data_key_reuse_period_seconds
resource "aws_sqs_queue" "encryption_without_reuse_period" {
  name              = "queue-without-encryption"
  kms_master_key_id = "alias/foo/bar"
}
