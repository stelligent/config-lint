# Setup Helper
resource "aws_kms_key" "test_key" {
}

# Pass
resource "aws_s3_bucket" "server_side_encryption_configuration_set" {
  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = aws_kms_key.test_key.arn
        sse_algorithm     = "aws:kms"
      }
    }
  }
}

# Fail
resource "aws_s3_bucket" "server_side_encryption_configuration_not_set" {
}
