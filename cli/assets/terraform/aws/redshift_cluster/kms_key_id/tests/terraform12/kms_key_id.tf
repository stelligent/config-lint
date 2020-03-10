variable "kms_key_arn" {
  default = "arn:aws:kms:us-east-1:1234567890:key/foobar"
}

# Warn
resource "aws_redshift_cluster" "kms_key_id_not_set" {
  cluster_identifier = "my-redshift-cluster"
  database_name      = "mydb"
  master_username    = "admin"
  master_password    = "foobarbaz"
  node_type          = "dc2.large"
  cluster_type       = "single-node"
}

# Pass
resource "aws_redshift_cluster" "kms_key_id_set" {
  cluster_identifier = "my-redshift-cluster"
  database_name      = "mydb"
  master_username    = "admin"
  master_password    = "foobarbaz"
  node_type          = "dc2.large"
  cluster_type       = "single-node"
  encrypted          = true
  kms_key_id         = "arn:aws:kms:us-east-1:1234567890:key/foobar"
}

# Pass
resource "aws_redshift_cluster" "kms_key_id_set_as_variable" {
  cluster_identifier = "my-redshift-cluster"
  database_name      = "mydb"
  master_username    = "admin"
  master_password    = "foobarbaz"
  node_type          = "dc2.large"
  cluster_type       = "single-node"
  encrypted          = true
  kms_key_id         = var.kms_key_arn
}
