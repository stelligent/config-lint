variable "kms_key_arn" {
  default = "arn:aws:kms:us-east-1:1234567890:key/foobar"
}

# Warn
resource "aws_redshift_cluster" "enhanced_vpc_routing_not_set" {
  cluster_identifier = "my-redshift-cluster"
  database_name      = "mydb"
  master_username    = "admin"
  master_password    = "foobarbaz"
  node_type          = "dc2.large"
  cluster_type       = "single-node"
  encrypted          = true
  kms_key_id         = var.kms_key_arn
}

# Pass
resource "aws_redshift_cluster" "enhanced_vpc_routing_set_to_false" {
  cluster_identifier   = "my-redshift-cluster"
  database_name        = "mydb"
  master_username      = "admin"
  master_password      = "foobarbaz"
  node_type            = "dc2.large"
  cluster_type         = "single-node"
  encrypted            = true
  kms_key_id           = var.kms_key_arn
  enhanced_vpc_routing = false
}

# Pass
resource "aws_redshift_cluster" "enhanced_vpc_routing_set_to_true" {
  cluster_identifier   = "my-redshift-cluster"
  database_name        = "mydb"
  master_username      = "admin"
  master_password      = "foobarbaz"
  node_type            = "dc2.large"
  cluster_type         = "single-node"
  encrypted            = true
  kms_key_id           = var.kms_key_arn
  enhanced_vpc_routing = true
}
