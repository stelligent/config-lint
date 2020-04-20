# Test that SNS topic policy does not use a wildcard principal for allow statements
# https://www.terraform.io/docs/providers/aws/r/sns_topic_policy.html#policy

provider "aws" {
  region = "us-east-1"
}

# PASS: SNS topic allow policy does not use a wildcard principal
resource "aws_sns_topic_policy" "sns_policy_allow_no_wildcard" {
  arn    = "arn:aws:sns:us-east-1:123456789012:test-topic"
  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sns:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/foo"
                ]
            },
            "Effect": "Allow",
            "Resource": "arn:aws:sns:us-east-1:123456789012:test-topic"
        }
    ]
}
EOF
}

# PASS: SNS topic deny policy does not use a wildcard principal
resource "aws_sns_topic_policy" "sns_policy_deny_no_wildcard" {
  arn    = "arn:aws:sns:us-east-1:123456789012:test-topic"
  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sns:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/foo"
                ]
            },
            "Effect": "Deny",
            "Resource": "arn:aws:sns:us-east-1:123456789012:test-topic"
        }
    ]
}
EOF
}

# PASS: SNS topic deny policy uses a wildcard principal
resource "aws_sns_topic_policy" "sns_policy_deny_with_wildcard" {
  arn    = "arn:aws:sns:us-east-1:123456789012:test-topic"
  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sns:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/*"
                ]
            },
            "Effect": "Deny",
            "Resource": "arn:aws:sns:us-east-1:123456789012:test-topic"
        }
    ]
}
EOF
}

# FAIL: SNS topic allow policy uses a wildcard principal
resource "aws_sns_topic_policy" "sns_policy_allow_with_wildcard" {
  arn    = "arn:aws:sns:us-east-1:123456789012:test-topic"
  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sns:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/*"
                ]
            },
            "Effect": "Allow",
            "Resource": "arn:aws:sns:us-east-1:123456789012:test-topic"
        }
    ]
}
EOF
}

# FAIL: SNS topic allow policy uses a wildcard principal
resource "aws_sns_topic_policy" "sns_policy_allow_principal_is_wildcard" {
  arn    = "arn:aws:sns:us-east-1:123456789012:test-topic"
  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sns:*",
            "Principal": {
                "AWS": [
                    "*"
                ]
            },
            "Effect": "Allow",
            "Resource": "arn:aws:sns:us-east-1:123456789012:test-topic"
        }
    ]
}
EOF
}
