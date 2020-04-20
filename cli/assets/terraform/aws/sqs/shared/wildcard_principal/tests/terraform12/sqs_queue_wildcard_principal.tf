# Test that SQS queue policy does not use a wildcard principal for allow statements
# https://www.terraform.io/docs/providers/aws/r/sqs_queue.html#policy

# PASS: SQS queue allow policy does not use a wildcard principal
resource "aws_sqs_queue" "sqs_policyallow_no_wildcard" {
  name   = "test-queue"
  policy = <<EOF
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
resource "aws_sqs_queue" "sqs_policydeny_no_wildcard" {
  name   = "test-queue"
  policy = <<EOF
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
resource "aws_sqs_queue" "sqs_policydeny_with_wildcard" {
  name   = "test-queue"
  policy = <<EOF
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
resource "aws_sqs_queue" "sqs_policyallow_with_wildcard" {
  name   = "test-queue"
  policy = <<EOF
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
resource "aws_sqs_queue" "sqs_policyallow_principal_is_wildcard" {
  name   = "test-queue"
  policy = <<EOF
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
