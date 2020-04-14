# Test that Neptune cluster has encryption enablded and KMS
# https://www.terraform.io/docs/providers/aws/r/neptune_cluster.html#storage_encrypted
# https://www.terraform.io/docs/providers/aws/r/neptune_cluster.html#kms_key_arn

provider "aws" {
  region = "us-east-1"
}

# PASS: Encryption is enabled with KMS
resource "aws_neptune_cluster" "storage_encrypted_set_to_true" {
  storage_encrypted = true
  kms_key_arn       = "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"
}

# WARN: Encryption is enabled without KMS
resource "aws_neptune_cluster" "storage_encrypted_set_to_true_no_kms" {
  storage_encrypted = true
}

# FAIL: Encryption is disabled
# WARN: KMS key is not specified
resource "aws_neptune_cluster" "storage_encrypted_set_to_false" {
  storage_encrypted = false
}

# FAIL: Encryption is not enabled
# WARN: KMS key is not specified
resource "aws_neptune_cluster" "storage_encrypted_not_set" {
}
