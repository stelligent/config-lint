resource "aws_redshift_cluster" "cluster" {
  cluster_identifier = "my-redshift-cluster"
  database_name      = "mydb"
  master_username    = "admin"
  master_password    = "foobarbaz"
  node_type          = "dc2.large"
  cluster_type       = "single-node"
}
