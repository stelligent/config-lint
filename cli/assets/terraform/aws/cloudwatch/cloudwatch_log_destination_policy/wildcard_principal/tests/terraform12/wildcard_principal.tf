# Test that CloudWatch log destination policy is not using a wildcard principal
# https://www.terraform.io/docs/providers/aws/r/cloudwatch_log_destination_policy.html#access_policy

provider "aws" {
  region = "us-east-1"
}

# PASS: Allow statement does not use a wildcard principal
resource "aws_cloudwatch_log_destination_policy" "cw_destination_no_wildcard" {
  destination_name = "cloudwatch_destination"
  access_policy    = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "cloudwatch:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/foo"
                ]
            },
            "Effect": "Allow",
            "Resource": "arn:aws:logs:us-west-1:123456789012:log-group:/mystack-testgroup-12ABC1AB12A1:*"
        }
    ]
}
EOF
}

# PASS: Deny statement does not use a wildcard principal
resource "aws_cloudwatch_log_destination_policy" "cw_destination_deny_no_wildcard" {
  destination_name = "cloudwatch_destination"
  access_policy    = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "cloudwatch:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/foo"
                ]
            },
            "Effect": "Deny",
            "Resource": "arn:aws:logs:us-west-1:123456789012:log-group:/mystack-testgroup-12ABC1AB12A1:*"
        }
    ]
}
EOF
}

# PASS: Deny statement uses a wildcard principal
resource "aws_cloudwatch_log_destination_policy" "cw_destination_deny_with_wildcard" {
  destination_name = "cloudwatch_destination"
  access_policy    = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "cloudwatch:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/*"
                ]
            },
            "Effect": "Deny",
            "Resource": "arn:aws:logs:us-west-1:123456789012:log-group:/mystack-testgroup-12ABC1AB12A1:*"
        }
    ]
}
EOF
}

# FAIL: Allow statement uses a wildcard principal
resource "aws_cloudwatch_log_destination_policy" "cw_destination_allow_with_wildcard" {
  destination_name = "cloudwatch_destination"
  access_policy    = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "cloudwatch:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/*"
                ]
            },
            "Effect": "Allow",
            "Resource": "arn:aws:logs:us-west-1:123456789012:log-group:/mystack-testgroup-12ABC1AB12A1:*"
        }
    ]
}
EOF
}

# FAIL: Allow statement uses a wildcard principal
resource "aws_cloudwatch_log_destination_policy" "cw_destination_principal_is_wildcard" {
  destination_name = "cloudwatch_destination"
  access_policy    = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "cloudwatch:*",
            "Principal": {
                "AWS": [
                    "*"
                ]
            },
            "Effect": "Allow",
            "Resource": "arn:aws:logs:us-west-1:123456789012:log-group:/mystack-testgroup-12ABC1AB12A1:*"
        }
    ]
}
EOF
}
