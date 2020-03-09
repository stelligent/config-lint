## Setup Helper
variable "test_cloudtrail_name" {
  default = "foo"
}
resource "aws_kms_key" "test_key" {
  enable_key_rotation = true
}

resource "aws_s3_bucket" "test_bucket" {
  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = "${aws_kms_key.test_key.arn}"
        sse_algorithm     = "aws:kms"
      }
    }
  }
}

# Pass
resource "aws_cloudtrail" "kms_key_id_is_set" {
  name           = var.test_cloudtrail_name
  s3_bucket_name = aws_s3_bucket.test_bucket.bucket
  kms_key_id     = aws_kms_key.test_key.arn
}

# Warn
resource "aws_cloudtrail" "kms_key_id_is_not_set" {
  name           = var.test_cloudtrail_name
  s3_bucket_name = aws_s3_bucket.test_bucket.bucket
}
