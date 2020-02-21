## Setup Helper
resource "aws_sns_topic" "test_topic" {
}

# Pass
resource "aws_sns_topic_policy" "policy_statement_without_notprincipal" {
  arn = "${aws_sns_topic.test_topic.arn}"

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

# Warn
resource "aws_sns_topic_policy" "policy_statement_with_notprincipal" {
  arn = "${aws_sns_topic.test_topic.arn}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Deny",
      "Action": "sns:Publish",
      "NotPrincipal": {
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
