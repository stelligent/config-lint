# Test that a DocumentDB cluster has encryption enabled
# https://www.terraform.io/docs/providers/aws/r/docdb_cluster.html#storage_encrypted

provider "aws" {
  region = "us-east-1"
}

# PASS: storage_encrypted defined and set to true
resource "aws_docdb_cluster" "docdb_storage_encrypted_true" {
  cluster_identifier      = "my-docdb-cluster"
  engine                  = "docdb"
  master_username         = "foo"
  master_password         = "mustbeeightchars"
  backup_retention_period = 5
  preferred_backup_window = "07:00-09:00"
  skip_final_snapshot     = true
  storage_encrypted       = true
}

# FAIL: storage_encrypted not defined
resource "aws_docdb_cluster" "docdb_storage_encrypted_not_defined" {
  cluster_identifier      = "my-docdb-cluster"
  engine                  = "docdb"
  master_username         = "foo"
  master_password         = "mustbeeightchars"
  backup_retention_period = 5
  preferred_backup_window = "07:00-09:00"
  skip_final_snapshot     = true
}

# FAIL: storage_encrypted set to false
resource "aws_docdb_cluster" "docdb_storage_encrypted_false" {
  cluster_identifier      = "my-docdb-cluster"
  engine                  = "docdb"
  master_username         = "foo"
  master_password         = "mustbeeightchars"
  backup_retention_period = 5
  preferred_backup_window = "07:00-09:00"
  skip_final_snapshot     = true
  storage_encrypted       = false
}
