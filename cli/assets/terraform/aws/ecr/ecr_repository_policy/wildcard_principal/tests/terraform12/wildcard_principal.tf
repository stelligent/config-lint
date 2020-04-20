# Test that ECR allow policy is not using a wildcard principal
# https://www.terraform.io/docs/providers/aws/r/ecr_repository_policy.html#policy

provider "aws" {
  region = "us-east-1"
}

# PASS: Allow policy not using wildcard principal
resource "aws_ecr_repository_policy" "ecr_allow_no_wildcard" {
  repository = "ecr-repo"

  policy = <<EOF
{
    "Version": "2008-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": "arn:aws:iam::1234567890:user/foo",
            "Action": [
                "ecr:*"
            ],
            "Resource": "*"
        }
    ]
}
EOF
}


# PASS: Deny policy using wildcard principal
resource "aws_ecr_repository_policy" "ecr_deny_wildcard" {
  repository = "ecr-repo"

  policy = <<EOF
{
    "Version": "2008-10-17",
    "Statement": [
        {
            "Effect": "Deny",
            "Principal": "arn:aws:iam::1234567890:user/*",
            "Action": [
                "ecr:*"
            ],
            "Resource": "*"
        }
    ]
}
EOF
}

# FAIL Allow policy using wildcard principal
resource "aws_ecr_repository_policy" "ecr_allow_with_wildcard" {
  repository = "ecr-repo"

  policy = <<EOF
{
    "Version": "2008-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": "arn:aws:iam::1234567890:user/*",
            "Action": [
                "ecr:*"
            ],
            "Resource": "*"
        }
    ]
}
EOF
}

# FAIL: Allow policy where principal is a wildcard
resource "aws_ecr_repository_policy" "ecr_allow_principal_is_wildcard" {
  repository = "ecr-repo"

  policy = <<EOF
{
    "Version": "2008-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": "*",
            "Action": [
                "ecr:*"
            ],
            "Resource": "*"
        }
    ]
}
EOF
}
