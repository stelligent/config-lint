# Test that server side encryption is used with a KMS key and aws:kms algorithm
# https://www.terraform.io/docs/providers/aws/r/s3_bucket.html#kms_master_key_id
# https://www.terraform.io/docs/providers/aws/r/s3_bucket.html#sse_algorithm

provider "aws" {
  region = "us-east-1"
}

# Setup Helper
resource "aws_kms_key" "test_key" {
}

# PASS: A KMS key is provided and sse_algorithm is aws:kms
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

# WARN: No KMS key is specified
resource "aws_s3_bucket" "server_side_encryption_configuration_set_no_key" {
  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "aws:kms"
      }
    }
  }
}

# WARN: No KMS key is specified
resource "aws_s3_bucket" "server_side_encryption_configuration_set_no_key_aes256" {
  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }
}

# FAIL: No server_side_encryption_configuration provided
# WARN: kms_master_key_id and sse_algorithm not set
resource "aws_s3_bucket" "server_side_encryption_configuration_not_set" {
}
