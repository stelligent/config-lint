## Setup Helper
resource "aws_kms_key" "test_key" {
}

resource "aws_lambda_function" "test_lambda" {
  filename      = "lambdatest.zip"
  function_name = "foobar"
  role          = "${aws_iam_role.test_role.arn}"
  handler       = "exports.handler"
  runtime       = "nodejs8.10"
  kms_key_arn   = "${aws_kms_key.test_key.arn}"
}

resource "aws_iam_role" "test_role" {
  name = "test_role"

  assume_role_policy = <<EOF
 {
   "Version": "2012-10-17",
   "Statement": [
     {
       "Effect": "Allow",
       "Action": "sts:AssumeRole",
       "Principal": {
         "Service": "lambda.amazonaws.com"
       }
     }
   ]
 }
 EOF
}

# Pass
resource "aws_lambda_permission" "principal_without_wildcard" {
  function_name = "${aws_lambda_function.test_lambda.function_name}"
  action        = "lambda:InvokeFunction"
  principal     = "events.amazonaws.com"
}

# Fail
resource "aws_lambda_permission" "principal_with_wildcard" {
  function_name = "${aws_lambda_function.test_lambda.function_name}"
  action        = "lambda:InvokeFunction"
  principal     = "*"
}

# Fail
resource "aws_lambda_permission" "principal_with_wildcard_prefix" {
  function_name = "${aws_lambda_function.test_lambda.function_name}"
  action        = "lambda:InvokeFunction"
  principal     = "*.amazonaws.com"
}
