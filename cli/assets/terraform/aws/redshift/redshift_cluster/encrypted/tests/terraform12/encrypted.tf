# Test that a Redshift cluster is encrypted with a KMS key
# https://www.terraform.io/docs/providers/aws/r/redshift_cluster.html#encrypted
# https://www.terraform.io/docs/providers/aws/r/redshift_cluster.html#kms_key_id

provider "aws" {
  region = "us-east-1"
}

# FAIL: Encryption disabled
# WARN: No KMS key
resource "aws_redshift_cluster" "encrypted_not_set" {
  cluster_identifier = "my-redshift-cluster"
  database_name      = "mydb"
  master_username    = "admin"
  master_password    = "F0obarbaz"
  node_type          = "dc2.large"
  cluster_type       = "single-node"
}

# FAIL: Encryption disabled
# WARN: No KMS key
resource "aws_redshift_cluster" "encrypted_set_to_false" {
  cluster_identifier = "my-redshift-cluster"
  database_name      = "mydb"
  master_username    = "admin"
  master_password    = "F0obarbaz"
  node_type          = "dc2.large"
  cluster_type       = "single-node"
  encrypted          = false
}

# PASS: Cluster encrypted with KMS key
resource "aws_redshift_cluster" "encrypted_set_to_true" {
  cluster_identifier = "my-redshift-cluster"
  database_name      = "mydb"
  master_username    = "admin"
  master_password    = "F0obarbaz"
  node_type          = "dc2.large"
  cluster_type       = "single-node"
  encrypted          = true
  kms_key_id         = "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"
}

# WARN: Cluster encrypted without KMS key
resource "aws_redshift_cluster" "encrypted_set_to_true_no_kms" {
  cluster_identifier = "my-redshift-cluster"
  database_name      = "mydb"
  master_username    = "admin"
  master_password    = "F0obarbaz"
  node_type          = "dc2.large"
  cluster_type       = "single-node"
  encrypted          = true
}
