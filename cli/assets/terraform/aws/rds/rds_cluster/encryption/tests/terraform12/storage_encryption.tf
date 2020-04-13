# Test that an RDS cluster has encryption endabled and a KMS key specified
# https://www.terraform.io/docs/providers/aws/r/rds_cluster.html#storage_encrypted
# https://www.terraform.io/docs/providers/aws/r/rds_cluster.html#kms_key_id

provider "aws" {
  region = "us-east-1"
}

# PASS: storage_encrypted enabled and kms_key_id specified
resource "aws_rds_cluster" "storage_encrypted_set_to_true" {
  engine            = "aurora-mysql"
  master_username   = "foo"
  master_password   = "bar"
  storage_encrypted = true
  kms_key_id        = "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"
}

# FAIL: storage_encrypted disabled
resource "aws_rds_cluster" "storage_encrypted_set_to_false" {
  engine            = "aurora-mysql"
  master_username   = "foo"
  master_password   = "bar"
  storage_encrypted = false
}

# FAIL: storage_encrypted not defined
resource "aws_rds_cluster" "storage_encrypted_not_set" {
  engine          = "aurora-mysql"
  master_username = "foo"
  master_password = "bar"
}

# PASS: serverless mode uses default encryption
resource "aws_rds_cluster" "serverless_engine_mode_encrypted_by_default" {
  engine          = "aurora-mysql"
  engine_mode     = "serverless"
  master_username = "foo"
  master_password = "bar"
}

# PASS: serverless mode with encryption enabled and a kms key
resource "aws_rds_cluster" "serverless_engine_mode_with_storage_encrypted_set_to_true" {
  engine            = "aurora-mysql"
  engine_mode       = "serverless"
  master_username   = "foo"
  master_password   = "bar"
  storage_encrypted = true
  kms_key_id        = "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"
}

# WARN: serverless mode with encryption enabled but no kms key
resource "aws_rds_cluster" "serverless_engine_mode_with_storage_encrypted_set_to_true_no_kms" {
  engine            = "aurora-mysql"
  engine_mode       = "serverless"
  master_username   = "foo"
  master_password   = "bar"
  storage_encrypted = true
}

# FAIL: encryption is disabled
resource "aws_rds_cluster" "serverless_engine_mode_with_storage_encrypted_set_to_false" {
  engine            = "aurora-mysql"
  engine_mode       = "serverless"
  master_username   = "foo"
  master_password   = "bar"
  storage_encrypted = false
}

# FAIL: Encryption not enabled
resource "aws_rds_cluster" "provisioned_engine_mode_unencrypted_by_default" {
  engine          = "aurora-mysql"
  engine_mode     = "provisioned"
  master_username = "foo"
  master_password = "bar"
}

# PASS: Encryption enabled with kms key
resource "aws_rds_cluster" "provisioned_engine_mode_with_storage_encrypted_set_to_true" {
  engine            = "aurora-mysql"
  engine_mode       = "provisioned"
  master_username   = "foo"
  master_password   = "bar"
  storage_encrypted = true
  kms_key_id        = "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"
}

# WARN: Encryption enabled without kms key
resource "aws_rds_cluster" "provisioned_engine_mode_with_storage_encrypted_set_to_true_no_kms" {
  engine            = "aurora-mysql"
  engine_mode       = "provisioned"
  master_username   = "foo"
  master_password   = "bar"
  storage_encrypted = true
}

# FAIL: Encrypytion disabled
resource "aws_rds_cluster" "provisioned_engine_mode_with_storage_encrypted_set_to_false" {
  engine            = "aurora-mysql"
  engine_mode       = "provisioned"
  master_username   = "foo"
  master_password   = "bar"
  storage_encrypted = false
}
