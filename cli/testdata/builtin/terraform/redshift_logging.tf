# Setup
variable "kms_key_arn" {
  default = "arn:aws:kms:us-east-1:1234567890:key/foobar"
}

resource "aws_s3_bucket" "audit_logs" {}


# # Warn
# resource "aws_redshift_cluster" "logging_not_set" {
#   cluster_identifier = "my-redshift-cluster"
#   database_name      = "mydb"
#   master_username    = "admin"
#   master_password    = "foobarbaz"
#   node_type          = "dc2.large"
#   cluster_type       = "single-node"
#   encrypted          = true
#   kms_key_id         = "arn:aws:kms:us-east-1:1234567890:key/foobar"
# }

# # Warn
# resource "aws_redshift_cluster" "logging_is_disabled" {
#   cluster_identifier = "my-redshift-cluster"
#   database_name      = "mydb"
#   master_username    = "admin"
#   master_password    = "foobarbaz"
#   node_type          = "dc2.large"
#   cluster_type       = "single-node"
#   encrypted          = true
#   kms_key_id         = "${var.kms_key_arn}"
#   logging {
#     enable        = false
#     bucket_name   = "${aws_s3_bucket.audit_logs.id}"
#     s3_key_prefix = "aws_redshift_cluster"
#   }
# }

# Pass
resource "aws_redshift_cluster" "logging_is_enabled" {
  cluster_identifier = "my-redshift-cluster"
  database_name      = "mydb"
  master_username    = "admin"
  master_password    = "foobarbaz"
  node_type          = "dc2.large"
  cluster_type       = "single-node"
  encrypted          = true
  kms_key_id         = "${var.kms_key_arn}"
  enable             = false
  logging {
    enable        = true
    bucket_name   = "${aws_s3_bucket.audit_logs.id}"
    s3_key_prefix = "aws_redshift_cluster"
  }
}
