# Test that encryption is enabled
# https://www.terraform.io/docs/providers/aws/r/elastictranscoder_pipeline.html#aws_kms_key_arn

provider "aws" {
  region = "us-east-1"
}

# PASS: KMS key is defined
resource "aws_elastictranscoder_pipeline" "transcoder_kms_key_defined" {
  input_bucket = "MyBucket"
  name         = "aws_elastictranscoder_pipeline_tf_test_"
  role         = "arn:aws:iam::123456789012:role/example-role"

  aws_kms_key_arn = "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"

  content_config {
    bucket        = "MyBucket"
    storage_class = "Standard"
  }

  thumbnail_config {
    bucket        = "MyBucket"
    storage_class = "Standard"
  }
}

# FAIL: KMS key is not defined
resource "aws_elastictranscoder_pipeline" "transcoder_kms_key_not_defined" {
  input_bucket = "MyBucket"
  name         = "aws_elastictranscoder_pipeline_tf_test_"
  role         = "arn:aws:iam::123456789012:role/example-role"

  content_config {
    bucket        = "MyBucket"
    storage_class = "Standard"
  }

  thumbnail_config {
    bucket        = "MyBucket"
    storage_class = "Standard"
  }
}
