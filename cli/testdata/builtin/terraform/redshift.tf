# Fail
resource "aws_redshift_cluster" "encrypted_not_set" {
  cluster_identifier = "my-redshift-cluster"
  database_name      = "mydb"
  master_username    = "admin"
  master_password    = "foobarbaz"
  node_type          = "dc2.large"
  cluster_type       = "single-node"
}

# Fail
resource "aws_redshift_cluster" "encrypted_set_to_false" {
  cluster_identifier = "my-redshift-cluster"
  database_name      = "mydb"
  master_username    = "admin"
  master_password    = "foobarbaz"
  node_type          = "dc2.large"
  cluster_type       = "single-node"
  encrypted          = false
}

# Pass
resource "aws_redshift_cluster" "encrypted_set_to_true" {
  cluster_identifier = "my-redshift-cluster"
  database_name      = "mydb"
  master_username    = "admin"
  master_password    = "foobarbaz"
  node_type          = "dc2.large"
  cluster_type       = "single-node"
  encrypted          = true
}
