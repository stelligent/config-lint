## Setup Helper
resource "aws_kms_key" "test_key" {
  enable_key_rotation = true
}

resource "aws_iam_role" "test_role" {
  assume_role_policy = "${data.aws_iam_policy_document.test_assume_role.json}"
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
resource "aws_sagemaker_notebook_instance" "kms_key_id_is_set" {
  name          = "foo"
  role_arn      = "${aws_iam_role.test_role.arn}"
  instance_type = "ml.t2.medium"
  kms_key_id    = "${aws_kms_key.test_key.id}"
}

# Warn
resource "aws_sagemaker_notebook_instance" "kms_key_id_is_not_set" {
  name          = "foo"
  role_arn      = "${aws_iam_role.test_role.arn}"
  instance_type = "ml.t2.medium"
}
