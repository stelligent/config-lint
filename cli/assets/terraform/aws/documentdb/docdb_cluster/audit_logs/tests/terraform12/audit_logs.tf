# Test that DocumentDB audit logging is enabled
# https://www.terraform.io/docs/providers/aws/r/docdb_cluster.html#enabled_cloudwatch_logs_exports

provider "aws" {
  region = "us-east-1"
}

# PASS: enabled_cloudwatch_logs_exports defined with audit logging
resource "aws_docdb_cluster" "docdb_audit_logging_enabled" {
  cluster_identifier      = "my-docdb-cluster"
  engine                  = "docdb"
  master_username         = "foo"
  master_password         = "mustbeeightchars"
  backup_retention_period = 5
  preferred_backup_window = "07:00-09:00"
  skip_final_snapshot     = true

  enabled_cloudwatch_logs_exports = ["audit"]
}

# FAIL: enabled_cloudwatch_logs_exports not defined
resource "aws_docdb_cluster" "docdb_cloudwatch_logs_exports_not_defined" {
  cluster_identifier      = "my-docdb-cluster"
  engine                  = "docdb"
  master_username         = "foo"
  master_password         = "mustbeeightchars"
  backup_retention_period = 5
  preferred_backup_window = "07:00-09:00"
  skip_final_snapshot     = true
}

# FAIL: enabled_cloudwatch_logs_exports defined without audit logging
resource "aws_docdb_cluster" "docdb_cloudwatch_logs_exports_without_audit" {
  cluster_identifier      = "my-docdb-cluster"
  engine                  = "docdb"
  master_username         = "foo"
  master_password         = "mustbeeightchars"
  backup_retention_period = 5
  preferred_backup_window = "07:00-09:00"
  skip_final_snapshot     = true

  enabled_cloudwatch_logs_exports = ["profiler"]
}
