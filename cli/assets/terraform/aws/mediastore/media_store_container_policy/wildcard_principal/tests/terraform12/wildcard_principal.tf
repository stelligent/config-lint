#
# https://www.terraform.io/docs/providers/aws/r/media_store_container_policy.html#policy

provider "aws" {
  region = "us-east-1"
}

# PASS: Allow policy with no wildcard principal
resource "aws_media_store_container_policy" "msc_allow_no_wildcard" {
  container_name = "example"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "mediastore:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890098:user/foo"
                ]
            },
            "Effect": "Allow",
            "Resource": "arn:aws:mediastore:1234567890098:us-east-1:container/example/*"
        }
    ]
}
EOF
}

# PASS: Deny policy with no wildcard principal
resource "aws_media_store_container_policy" "msc_deny_no_wildcard" {
  container_name = "example"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "mediastore:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890098:user/foo"
                ]
            },
            "Effect": "Deny",
            "Resource": "arn:aws:mediastore:1234567890098:us-east-1:container/example/*"
        }
    ]
}
EOF
}

# PASS: Deny policy with wildcard principal
resource "aws_media_store_container_policy" "msc_deny_wildcard" {
  container_name = "example"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "mediastore:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890098:user/*"
                ]
            },
            "Effect": "Deny",
            "Resource": "arn:aws:mediastore:1234567890098:us-east-1:container/example/*"
        }
    ]
}
EOF
}

# FAIL: Allow policy with wildcard in principal
resource "aws_media_store_container_policy" "msc_allow_with_wildcard" {
  container_name = "example"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "mediastore:*",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::1234567890098:user/*"
                ]
            },
            "Effect": "Allow",
            "Resource": "arn:aws:mediastore:1234567890098:us-east-1:container/example/*"
        }
    ]
}
EOF
}

# FAIL: Allow policy principal is a wildcard
resource "aws_media_store_container_policy" "msc_allow_principal_is_wildcard" {
  container_name = "example"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "mediastore:*",
            "Principal": {
                "AWS": [
                    "*"
                ]
            },
            "Effect": "Allow",
            "Resource": "arn:aws:mediastore:1234567890098:us-east-1:container/example/*"
        }
    ]
}
EOF
}
