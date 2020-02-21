# Pass
resource "aws_sqs_queue" "encryption_enabled" {
  name              = "queue-with-encryption"
  kms_master_key_id = "alias/foo/bar"
}

# Pass
resource "aws_sqs_queue" "encryption_enabled_with_reuse_set" {
  name                              = "queue-with-encryption"
  kms_master_key_id                 = "alias/foo/bar"
  kms_data_key_reuse_period_seconds = 60
}

# Fail
resource "aws_sqs_queue" "encryption_disabled" {
  name = "queue-without-encryption"
}
