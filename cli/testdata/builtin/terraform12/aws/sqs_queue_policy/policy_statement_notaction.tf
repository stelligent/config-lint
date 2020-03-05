## Setup Helper
resource "aws_sqs_queue" "test_queue" {
  name                              = "test_queue"
  kms_master_key_id                 = "alias/foo/bar"
  kms_data_key_reuse_period_seconds = 60
  arn                               = "mocked_arn"
}

# Pass
resource "aws_sqs_queue_policy" "policy_statement_without_notaction" {
  queue_url = aws_sqs_queue.test_queue.id

  policy = <<EOF
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
EOF
}

# Warn
resource "aws_sqs_queue_policy" "policy_statement_with_notaction" {
  queue_url = aws_sqs_queue.test_queue.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "NotAction": "sqs:SendMessage",
      "Resource": "${aws_sqs_queue.test_queue.arn}"
    }
  ]
}
EOF
}
