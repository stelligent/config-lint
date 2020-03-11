## Setup Helper
resource "aws_sqs_queue" "test_queue" {
  name                              = "test_queue"
  kms_master_key_id                 = "alias/foo/bar"
  kms_data_key_reuse_period_seconds = 60
  arn                               = "mockedarn"
}

# Pass
resource "aws_sqs_queue_policy" "policy_statement_allow_action_without_wildcard" {
  queue_url = aws_sqs_queue.test_queue.id

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sqs:SendMessage",
      "Resource": "${aws_sqs_queue.test_queue.arn}"
    }
  ]
}
POLICY
}

# Pass
resource "aws_sqs_queue_policy" "policy_statement_deny_action_without_wildcard" {
  queue_url = aws_sqs_queue.test_queue.id

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Deny",
      "Action": "sqs:SendMessage",
      "Resource": "${aws_sqs_queue.test_queue.arn}"
    }
  ]
}
POLICY
}

# Pass
resource "aws_sqs_queue_policy" "policy_statement_deny_action_with_wildcard" {
  queue_url = aws_sqs_queue.test_queue.id

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Deny",
      "Action": "sqs:*",
      "Resource": "${aws_sqs_queue.test_queue.arn}"
    }
  ]
}
POLICY
}

# Fail
resource "aws_sqs_queue_policy" "policy_statement_allow_action_with_wildcard" {
  queue_url = aws_sqs_queue.test_queue.id

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sqs:*",
      "Resource": "${aws_sqs_queue.test_queue.arn}"
    }
  ]
}
POLICY
}
