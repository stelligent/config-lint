## Setup Helper
resource "aws_s3_bucket" "test_bucket" {
}

# Pass
resource "aws_emr_cluster" "log_uri_is_set" {
  name          = "foo"
  release_label = "emr-4.6.0"
  service_role  = "arn:aws:iam::1234567890:role/EMR_DefaultRole"
  log_uri       = "s3://${aws_s3_bucket.test_bucket.bucket}/"
}

# Fail
resource "aws_emr_cluster" "log_uri_is_set" {
  name          = "foo"
  release_label = "emr-4.6.0"
  service_role  = "arn:aws:iam::1234567890:role/EMR_DefaultRole"
}
