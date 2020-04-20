# Test that a KMS key policy allow is not using a wildcard for the principal
# https://www.terraform.io/docs/providers/aws/r/kms_key.html#policy

provider "aws" {
  region = "us-east-1"
}

# PASS: KMS key policy is an allow not using a wildcard principal
resource "aws_kms_key" "kms_key_allow_no_wildcard" {
  description             = "Example Key"
  deletion_window_in_days = 10

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [{
        "Sid": "Enable IAM User Permissions",
        "Effect": "Allow",
        "Principal": {
            "AWS": "arn:aws:iam::1234567890:user/foo"
        },
        "Action": "kms:*",
        "Resource": "*"
    }]
}
EOF
}


# PASS: KMS key policy is deny not using a wildcard principal
resource "aws_kms_key" "kms_key_deny_without_wildcard" {
  description             = "Example Key"
  deletion_window_in_days = 10

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [{
        "Sid": "Enable IAM User Permissions",
        "Effect": "Deny",
        "Principal": {
            "AWS": "arn:aws:iam::1234567890:user"
        },
        "Action": "kms:*",
        "Resource": "*"
    }]
}
EOF
}

# PASS: KMS key policy is deny using a wildcard principal
resource "aws_kms_key" "kms_key_deny_with_wildcard" {
  description             = "Example Key"
  deletion_window_in_days = 10

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [{
        "Sid": "Enable IAM User Permissions",
        "Effect": "Deny",
        "Principal": {
            "AWS": "arn:aws:iam::1234567890:user/*"
        },
        "Action": "kms:*",
        "Resource": "*"
    }]
}
EOF
}

# FAIL: KMS key policy is an allow using a wildcard principal
resource "aws_kms_key" "kms_key_allow_with_wildcard" {
  description             = "Example Key"
  deletion_window_in_days = 10

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [{
        "Sid": "Enable IAM User Permissions",
        "Effect": "Allow",
        "Principal": {
            "AWS": "arn:aws:iam::1234567890:user/*"
        },
        "Action": "kms:*",
        "Resource": "*"
    }]
}
EOF
}

# FAIL: KMS key policy is an allow using a wildcard principal
resource "aws_kms_key" "kms_key_principal_is_wildcard" {
  description             = "Example Key"
  deletion_window_in_days = 10

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [{
        "Sid": "Enable IAM User Permissions",
        "Effect": "Allow",
        "Principal": {
            "AWS": "*"
        },
        "Action": "kms:*",
        "Resource": "*"
    }]
}
EOF
}
