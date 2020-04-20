# Test that SQS queue policy does not use a wildcard principal for allow statements
# https://www.terraform.io/docs/providers/aws/r/sqs_queue_policy.html#policy

provider "aws" {
  region = "us-east-1"
}

# PASS: SQS queue allow policy does not use a wildcard principal
resource "aws_sqs_queue_policy" "sqs_policy_allow_no_wildcard" {
  queue_url = "https://queue.amazonaws.com/0123456789012/myqueue"
  policy    = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sqs:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/foo"
                ]
            },
            "Effect": "Allow",
            "Resource": "arn:aws:sqs:us-east-2:444455556666:queue1"
        }
    ]
}
EOF
}

# PASS: SQS queue deny policy does not use a wildcard principal
resource "aws_sqs_queue_policy" "sqs_policy_deny_no_wildcard" {
  queue_url = "https://queue.amazonaws.com/0123456789012/myqueue"
  policy    = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sqs:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/foo"
                ]
            },
            "Effect": "Deny",
            "Resource": "arn:aws:sqs:us-east-2:444455556666:queue1"
        }
    ]
}
EOF
}

# PASS: SQS queue deny policy uses a wildcard principal
resource "aws_sqs_queue_policy" "sqs_policy_deny_with_wildcard" {
  queue_url = "https://queue.amazonaws.com/0123456789012/myqueue"
  policy    = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sqs:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/*"
                ]
            },
            "Effect": "Deny",
            "Resource": "arn:aws:sqs:us-east-2:444455556666:queue1"
        }
    ]
}
EOF
}

# FAIL: SQS queue allow policy uses a wildcard principal
resource "aws_sqs_queue_policy" "sqs_policy_allow_with_wildcard" {
  queue_url = "https://queue.amazonaws.com/0123456789012/myqueue"
  policy    = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sqs:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/*"
                ]
            },
            "Effect": "Allow",
            "Resource": "arn:aws:sqs:us-east-2:444455556666:queue1"
        }
    ]
}
EOF
}

# FAIL: SQS queue allow policy uses a wildcard principal
resource "aws_sqs_queue_policy" "sqs_policy_allow_principal_is_wildcard" {
  queue_url = "https://queue.amazonaws.com/0123456789012/myqueue"
  policy    = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sqs:*",
            "Principal": {
                "AWS": [
                    "*"
                ]
            },
            "Effect": "Allow",
            "Resource": "arn:aws:sqs:us-east-2:444455556666:queue1"
        }
    ]
}
EOF
}
