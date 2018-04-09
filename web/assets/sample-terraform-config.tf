resource "aws_kms_key" "key_for_s3_encryption" {
  description = "key for S3 bucket encryption"
}

resource "aws_s3_bucket" "bucket_example_1" {
  acl = "public-read"
}

resource "aws_s3_bucket" "bucket_example_2" {
  acl = "public-read-write"
  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = "${aws_kms_key.key_for_s3_encryption.arn}"
        sse_algorithm     = "aws:kms"
      }
    }
  }
}
