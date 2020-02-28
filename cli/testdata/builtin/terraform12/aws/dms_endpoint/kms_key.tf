## Setup Helper
variable "test_id" {
  default = "foo"
}

variable "test_engine_type" {
  default = "source"
}

variable "test_engine_name" {
  default = "aurora"
}

resource "aws_kms_key" "test_key" {
  enable_key_rotation = true
}

# Pass
resource "aws_dms_endpoint" "kms_key_arn_is_set" {
  endpoint_id   = var.test_id
  endpoint_type = var.test_engine_type
  engine_name   = var.test_engine_name
  kms_key_arn   = aws_kms_key.test_key.arn
}

# Warn
resource "aws_dms_endpoint" "kms_key_arn_is_not_set" {
  endpoint_id   = var.test_id
  endpoint_type = var.test_engine_type
  engine_name   = var.test_engine_name
}
