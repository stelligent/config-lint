## Setup Helper
resource "aws_sns_topic" "test_topic" {
  name = "test_topic"
}

# Pass
resource "aws_sns_topic_policy" "policy_version_set_correctly" {
  arn = "${aws_sns_topic.test_topic.arn}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sns:Subscribe",
      "Resource": "arn:aws:sns:us-east-1:123456789012:foobar"
    }
  ]
}
EOF
}

# Fail
resource "aws_sns_topic_policy" "policy_version_set_incorrectly" {
  arn = "${aws_sns_topic.test_topic.arn}"

  policy = <<EOF
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sns:Subscribe",
      "Resource": "arn:aws:sns:us-east-1:123456789012:foobar"
    }
  ]
}
EOF
}
