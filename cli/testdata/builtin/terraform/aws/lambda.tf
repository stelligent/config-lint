
#resource "aws_api_gateway" "gateway1" {
#}
#
resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_lambda_function" "test_lambda" {
  filename         = "lambda_function_payload.zip"
  function_name    = "example"
  role             = "${aws_iam_role.iam_for_lambda.arn}"
  handler          = "handler.handler"
  source_code_hash = "${base64sha256(file("lambda_function_payload.zip"))}"
  runtime          = "nodejs8.10"

  environment {
    variables = {
      key = "AKIAIOSFODNN7EXAMPLE"
      secret = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    }
  }
}

resource "aws_lambda_alias" "test_alias" {
  name             = "testalias"
  description      = "a sample description"
  function_name    = "${aws_lambda_function.test_lambda.function_name}"
  function_version = "$LATEST"
}

resource "aws_lambda_permission" "permission1" {
  statement_id   = "AllowExecutionFromCloudWatch"
  action         = "lambda:InvokeFunction"
  function_name  = "${aws_lambda_function.test_lambda.function_name}"
  principal      = "events.amazonaws.com"
  source_arn     = "arn:aws:events:us-west-2:111122223333:rule/RunDaily"
  qualifier      = "${aws_lambda_alias.test_alias.name}"
}
