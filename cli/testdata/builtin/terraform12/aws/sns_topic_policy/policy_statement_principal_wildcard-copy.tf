## Setup Helper
resource "aws_sns_topic" "test_topic" {
}

# Pass
resource "aws_sns_topic_policy" "policy_statement_allow_principal_without_wildcard" {
  arn = aws_sns_topic.test_topic.arn

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sns:Publish",
      "Principal": {
        "AWS": [
          "arn:aws:iam::1234567890:user/foo"
        ]
      },
      "Resource": "arn:aws:sns:us-east-1:1234567890:fooTopic"
    }
  ]
}
EOF
}

# Pass
resource "aws_sns_topic_policy" "policy_statement_deny_principal_without_wildcard" {
  arn = aws_sns_topic.test_topic.arn

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Deny",
      "Action": "sns:Publish",
      "Principal": {
        "AWS": [
          "arn:aws:iam::1234567890:user/foo"
        ]
      },
      "Resource": "arn:aws:sns:us-east-1:1234567890:fooTopic"
    }
  ]
}
EOF
}

# Pass
resource "aws_sns_topic_policy" "policy_statement_deny_principal_with_wildcard" {
  arn = "${aws_sns_topic.test_topic.arn}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Deny",
      "Action": "sns:Publish",
      "Principal": {
        "AWS": [
          "arn:aws:iam::1234567890:user/foo*"
        ]
      },
      "Resource": "arn:aws:sns:us-east-1:1234567890:fooTopic"
    }
  ]
}
EOF
}

# Fail
resource "aws_sns_topic_policy" "policy_statement_allow_principal_with_wildcard" {
  arn = aws_sns_topic.test_topic.arn

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sns:Publish",
      "Principal": {
        "AWS": [
          "arn:aws:iam::1234567890:user/foo*"
        ]
      },
      "Resource": "arn:aws:sns:us-east-1:1234567890:fooTopic"
    }
  ]
}
EOF
}
