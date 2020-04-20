# Test that IOT allow statement is not using a wildcard principal
# https://www.terraform.io/docs/providers/aws/r/iot_policy.html#policy

provider "aws" {
  region = "us-east-1"
}

# PASS: Allow with no wildcard principal
resource "aws_iot_policy" "iot_allow_no_wildcard" {
  name = "PubSubToAnyTopic"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "iot:*"
      ],
      "Principal": "arn:aws:iam::1234567890:user/foo",
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

# PASS: Deny with no wildcard principal
resource "aws_iot_policy" "iot_deny_no_wildcard" {
  name = "PubSubToAnyTopic"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "iot:*"
      ],
      "Principal": "arn:aws:iam::1234567890:user/foo",
      "Effect": "Deny",
      "Resource": "*"
    }
  ]
}
EOF
}

# PASS: Deny with a wildcard principal
resource "aws_iot_policy" "iot_deny_with_wildcard" {
  name = "PubSubToAnyTopic"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "iot:*"
      ],
      "Principal": "arn:aws:iam::1234567890:user/*",
      "Effect": "Deny",
      "Resource": "*"
    }
  ]
}
EOF
}

# FAIL: Allow with wildcard principal
resource "aws_iot_policy" "iot_allow_with_wildcard" {
  name = "PubSubToAnyTopic"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "iot:*"
      ],
      "Principal": "arn:aws:iam::1234567890:user/*",
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

# FAIL: Allow with wildcard principal
resource "aws_iot_policy" "iot_allow_principal_is_wildcard" {
  name = "PubSubToAnyTopic"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "iot:*"
      ],
      "Principal": "*",
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}
