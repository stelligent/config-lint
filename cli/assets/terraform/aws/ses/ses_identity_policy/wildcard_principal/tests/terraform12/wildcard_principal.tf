# Test that SES identity policy is not using a wildcard principal for allow statements
# https://www.terraform.io/docs/providers/aws/r/ses_identity_policy.html#policy

provider "aws" {
  region = "us-east-1"
}

# PASS: SES identity allow without using a wildcard principal
resource "aws_ses_identity_policy" "ses_allow_without_wildcard" {
  identity = "arn:aws:ses:us-west-2:123456789012:identity/example.com"
  name     = "example"
  policy   = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "ses:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/foo"
                ]
            },
            "Effect": "Allow",
            "Resource": "arn:aws:ses:us-west-2:123456789012:identity/example.com"
        }
    ]
}
EOF
}

# PASS: Deny without using a wildcard principal
resource "aws_ses_identity_policy" "ses_deny_without_wildcard" {
  identity = "arn:aws:ses:us-west-2:123456789012:identity/example.com"
  name     = "example"
  policy   = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "ses:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/foo"
                ]
            },
            "Effect": "Deny",
            "Resource": "arn:aws:ses:us-west-2:123456789012:identity/example.com"
        }
    ]
}
EOF
}

# PASS: Deny using a wildcard principal
resource "aws_ses_identity_policy" "ses_deny_with_wildcard" {
  identity = "arn:aws:ses:us-west-2:123456789012:identity/example.com"
  name     = "example"
  policy   = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "ses:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/*"
                ]
            },
            "Effect": "Deny",
            "Resource": "arn:aws:ses:us-west-2:123456789012:identity/example.com"
        }
    ]
}
EOF
}

# FAIL: Allow using a wildcard principal
resource "aws_ses_identity_policy" "ses_allow_with_wildcard" {
  identity = "arn:aws:ses:us-west-2:123456789012:identity/example.com"
  name     = "example"
  policy   = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "ses:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890:user/*"
                ]
            },
            "Effect": "Allow",
            "Resource": "arn:aws:ses:us-west-2:123456789012:identity/example.com"
        }
    ]
}
EOF
}

# PASS: Allow where principal is a wildcard
resource "aws_ses_identity_policy" "ses_allow_principal_is_wildcard" {
  identity = "arn:aws:ses:us-west-2:123456789012:identity/example.com"
  name     = "example"
  policy   = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "ses:*",
            "Principal": {
                "AWS": [
                    "*"
                ]
            },
            "Effect": "Allow",
            "Resource": "arn:aws:ses:us-west-2:123456789012:identity/example.com"
        }
    ]
}
EOF
}
