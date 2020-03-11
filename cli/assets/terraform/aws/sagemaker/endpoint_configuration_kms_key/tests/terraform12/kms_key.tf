## Setup Helper
resource "aws_kms_key" "test_key" {
  enable_key_rotation = true
}

resource "aws_sagemaker_model" "test_model" {
  execution_role_arn = aws_iam_role.test_role.arn

  primary_container {
    image = "1234567890.dkr.ecr.us-east-1.amazonaws.com/foo:1"
  }
}

resource "aws_iam_role" "test_role" {
  assume_role_policy = data.aws_iam_policy_document.test_assume_role.json
}

data "aws_iam_policy_document" "test_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["sagemaker.amazonaws.com"]
    }
  }
}

# Pass
resource "aws_sagemaker_endpoint_configuration" "kms_key_arn_is_set" {
  kms_key_arn = aws_kms_key.test_key.arn

  production_variants {
    model_name             = aws_sagemaker_model.test_model.name
    initial_instance_count = 1
    instance_type          = "ml.t2.medium"
  }
}

# Warn
resource "aws_sagemaker_endpoint_configuration" "kms_key_arn_is_not_set" {
  production_variants {
    model_name             = aws_sagemaker_model.test_model.name
    initial_instance_count = 1
    instance_type          = "ml.t2.medium"
  }
}
