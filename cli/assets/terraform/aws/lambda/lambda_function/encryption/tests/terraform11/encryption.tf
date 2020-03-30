## Setup Helper
resource "aws_kms_key" "test_key" {
}

# Pass
resource "aws_lambda_function" "kms_key_arn_set" {
  filename      = "lambdatest.zip"
  function_name = "foobar"
  role          = "${aws_iam_role.test_role.arn}"
  handler       = "exports.handler"
  runtime       = "nodejs8.10"
  kms_key_arn   = "${aws_kms_key.test_key.arn}"
}

# Warn
resource "aws_lambda_function" "kms_key_arn_not_set" {
  filename      = "lambdatest.zip"
  function_name = "foobar"
  role          = "${aws_iam_role.test_role.arn}"
  handler       = "exports.handler"
  runtime       = "nodejs8.10"
}
