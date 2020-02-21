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
resource "aws_lambda_permission" "action_invokefunction_set" {
  function_name = "${aws_lambda_function.test_lambda.function_name}"
  action        = "lambda:InvokeFunction"
  principal     = "events.amazonaws.com"
}

# Warn
resource "aws_lambda_permission" "action_invokefunction_not_set" {
  function_name = "${aws_lambda_function.test_lambda.function_name}"
  action        = "lambda:GetFunction"
  principal     = "events.amazonaws.com"
}

