## Setup Helper
resource "aws_kms_key" "test_key" {
  enable_key_rotation = true
}

# Pass
resource "aws_dms_endpoint" "kms_key_arn_is_set" {
  endpoint_id   = "foo"
  endpoint_type = "source"
  engine_name   = "aurora"
  kms_key_arn   = "${aws_kms_key.test_key.arn}"
}

# Warn
resource "aws_dms_endpoint" "kms_key_arn_is_not_set" {
  endpoint_id   = "foo"
  endpoint_type = "source"
  engine_name   = "aurora"
}
