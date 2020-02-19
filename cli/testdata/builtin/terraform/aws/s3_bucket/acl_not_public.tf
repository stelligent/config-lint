# Pass
resource "aws_s3_bucket" "acl_not_set" {
}

# Pass
resource "aws_s3_bucket" "acl_set_to_private" {
  acl = "private"
}

# Pass
resource "aws_s3_bucket" "acl_set_to_aws-exec-read" {
  acl = "aws-exec-read"
}

# Pass
resource "aws_s3_bucket" "acl_set_to_authenticated-read" {
  acl = "authenticated-read"
}

# Pass
resource "aws_s3_bucket" "acl_set_to_bucket-owner-read" {
  acl = "bucket-owner-read"
}

# Pass
resource "aws_s3_bucket" "acl_set_to_bucket-owner-full-control" {
  acl = "bucket-owner-full-control"
}

# Pass
resource "aws_s3_bucket" "acl_set_to_log-delivery-write" {
  acl = "log-delivery-write"
}

# Fail
resource "aws_s3_bucket" "acl_set_to_public-read" {
  acl = "public-read"
}

# Fail
resource "aws_s3_bucket" "acl_set_to_public-read-write" {
  acl = "public-read-write"
}
