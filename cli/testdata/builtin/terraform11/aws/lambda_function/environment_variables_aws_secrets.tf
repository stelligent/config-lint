## Setup Helper
resource "aws_kms_key" "test_key" {
  enable_key_rotation = true
}

# Pass
resource "aws_lambda_function" "environment_variables_aws_secrets_not_set" {
  filename      = "lambdatest.zip"
  function_name = "foobar"
  role          = "${aws_iam_role.test_role.arn}"
  handler       = "exports.handler"
  runtime       = "nodejs8.10"
  kms_key_arn   = "${aws_kms_key.test_key.arn}"

  environment {
    variables = {
      foo = "bar"
    }
  }
}

# Pass
resource "aws_lambda_function" "environment_variables_aws_secrets_not_set_20_character_capital_string" {
  filename      = "lambdatest.zip"
  function_name = "foobar"
  role          = "${aws_iam_role.test_role.arn}"
  handler       = "exports.handler"
  runtime       = "nodejs8.10"
  kms_key_arn   = "${aws_kms_key.test_key.arn}"

  environment {
    variables = {
      foo = "AXYZIOSFODNN7EXAMPLE"
    }
  }
}

# Pass
resource "aws_lambda_function" "environment_variables_aws_secrets_not_set_21_character_capital_string" {
  filename      = "lambdatest.zip"
  function_name = "foobar"
  role          = "${aws_iam_role.test_role.arn}"
  handler       = "exports.handler"
  runtime       = "nodejs8.10"
  kms_key_arn   = "${aws_kms_key.test_key.arn}"

  environment {
    variables = {
      foo = "AKIAIOSFOODNN7EXAMPLE"
    }
  }
}

# Pass
resource "aws_lambda_function" "environment_variables_aws_secrets_not_set_40_character_string" {
  filename      = "lambdatest.zip"
  function_name = "foobar"
  role          = "${aws_iam_role.test_role.arn}"
  handler       = "exports.handler"
  runtime       = "nodejs8.10"
  kms_key_arn   = "${aws_kms_key.test_key.arn}"

  environment {
    variables = {
      foo = "wJalrXUtnFEMI_K7MDENG=bPxRfiCYEXAMPLEKEY"
    }
  }
}

# Pass
resource "aws_lambda_function" "environment_variables_aws_secrets_not_set_41_character_string" {
  filename      = "lambdatest.zip"
  function_name = "foobar"
  role          = "${aws_iam_role.test_role.arn}"
  handler       = "exports.handler"
  runtime       = "nodejs8.10"
  kms_key_arn   = "${aws_kms_key.test_key.arn}"

  environment {
    variables = {
      foo = "wJalrXUtnFOEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    }
  }
}

# Fail
resource "aws_lambda_function" "environment_variables_aws_secrets_access_key_set" {
  filename      = "lambdatest.zip"
  function_name = "foobar"
  role          = "${aws_iam_role.test_role.arn}"
  handler       = "exports.handler"
  runtime       = "nodejs8.10"
  kms_key_arn   = "${aws_kms_key.test_key.arn}"

  environment {
    variables = {
      foo = "AKIAIOSFODNN7EXAMPLE"
    }
  }
}

# Fail
resource "aws_lambda_function" "environment_variables_aws_secrets_secret_access_key_set" {
  filename      = "lambdatest.zip"
  function_name = "foobar"
  role          = "${aws_iam_role.test_role.arn}"
  handler       = "exports.handler"
  runtime       = "nodejs8.10"
  kms_key_arn   = "${aws_kms_key.test_key.arn}"

  environment {
    variables = {
      foo = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    }
  }
}

# Fail
resource "aws_lambda_function" "environment_variables_aws_secrets_access_key_and_secret_access_key_set" {
  filename      = "lambdatest.zip"
  function_name = "foobar"
  role          = "${aws_iam_role.test_role.arn}"
  handler       = "exports.handler"
  runtime       = "nodejs8.10"
  kms_key_arn   = "${aws_kms_key.test_key.arn}"

  environment {
    variables = {
      foo = "AKIAIOSFODNN7EXAMPLE"
      bar = "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY"
    }
  }
}
