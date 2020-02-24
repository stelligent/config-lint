# Pass
resource "aws_rds_cluster" "storage_encrypted_set_to_true" {
  engine            = "aurora-mysql"
  master_username   = "foo"
  master_password   = "bar"
  storage_encrypted = true
}

# Fail
resource "aws_rds_cluster" "storage_encrypted_set_to_false" {
  engine            = "aurora-mysql"
  master_username   = "foo"
  master_password   = "bar"
  storage_encrypted = false
}

# Fail
resource "aws_rds_cluster" "storage_encrypted_not_set" {
  engine          = "aurora-mysql"
  master_username = "foo"
  master_password = "bar"
}

# Pass
resource "aws_rds_cluster" "serverless_engine_mode_encrypted_by_default" {
  engine          = "aurora-mysql"
  engine_mode     = "serverless"
  master_username = "foo"
  master_password = "bar"
}

# Pass
resource "aws_rds_cluster" "serverless_engine_mode_with_storage_encrypted_set_to_true" {
  engine            = "aurora-mysql"
  engine_mode       = "serverless"
  master_username   = "foo"
  master_password   = "bar"
  storage_encrypted = true
}

# Fail
resource "aws_rds_cluster" "serverless_engine_mode_with_storage_encrypted_set_to_false" {
  engine            = "aurora-mysql"
  engine_mode       = "serverless"
  master_username   = "foo"
  master_password   = "bar"
  storage_encrypted = false
}

# Fail
resource "aws_rds_cluster" "provisioned_engine_mode_unencrypted_by_default" {
  engine          = "aurora-mysql"
  engine_mode     = "provisioned"
  master_username = "foo"
  master_password = "bar"
}

# Pass
resource "aws_rds_cluster" "provisioned_engine_mode_with_storage_encrypted_set_to_true" {
  engine            = "aurora-mysql"
  engine_mode       = "provisioned"
  master_username   = "foo"
  master_password   = "bar"
  storage_encrypted = true
}

# Fail
resource "aws_rds_cluster" "provisioned_engine_mode_with_storage_encrypted_set_to_false" {
  engine            = "aurora-mysql"
  engine_mode       = "provisioned"
  master_username   = "foo"
  master_password   = "bar"
  storage_encrypted = false
}
