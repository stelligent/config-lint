## Setup Helper
resource "aws_kms_key" "test_key" {
}

resource "aws_s3_bucket" "test_bucket" {
  acl = "private"
}

# Pass
resource "aws_s3_bucket_object" "encrypt_with_kms_key_id" {
  key        = "foo"
  bucket     = "${aws_s3_bucket.test_bucket.id}"
  kms_key_id = "${aws_kms_key.test_key.arn}"
}

# Pass
resource "aws_s3_bucket_object" "encrypt_with_server_side_encryption_s3_default_master_key" {
  key                    = "foo"
  bucket                 = "${aws_s3_bucket.test_bucket.id}"
  server_side_encryption = "aws:kms"
}

# Pass
resource "aws_s3_bucket_object" "encrypt_with_server_side_encryption_aws_managed_key" {
  key                    = "foo"
  bucket                 = "${aws_s3_bucket.test_bucket.id}"
  server_side_encryption = "AES256"
}

# Fail
resource "aws_s3_bucket_object" "encryption_method_not_set" {
  key    = "foo"
  bucket = "${aws_s3_bucket.test_bucket.id}"
}
