## Setup Helper
resource "aws_sqs_queue" "test_queue" {
  name                              = "test_queue"
  kms_master_key_id                 = "alias/foo/bar"
  kms_data_key_reuse_period_seconds = 60
}

# Pass
resource "aws_sqs_queue_policy" "policy_statement_allow_principal_without_wildcard" {
  queue_url = aws_sqs_queue.test_queue.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sqs:SendMessage",
      "Principal": {
        "AWS": [
          "arn:aws:iam::1234567890:user/foo"
        ]
      },
      "Resource": "${aws_sqs_queue.test_queue.arn}"
    }
  ]
}
EOF
}

# Pass
resource "aws_sqs_queue_policy" "policy_statement_deny_principal_without_wildcard" {
  queue_url = aws_sqs_queue.test_queue.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Deny",
      "Action": "sqs:SendMessage",
      "Principal": {
        "AWS": [
          "arn:aws:iam::1234567890:user/foo"
        ]
      },
      "Resource": "${aws_sqs_queue.test_queue.arn}"
    }
  ]
}
EOF
}

# Pass
resource "aws_sqs_queue_policy" "policy_statement_deny_principal_with_wildcard" {
  queue_url = aws_sqs_queue.test_queue.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Deny",
      "Action": "sqs:SendMessage",
      "Principal": "*",
      "Resource": "${aws_sqs_queue.test_queue.arn}"
    }
  ]
}
EOF
}

# Fail
resource "aws_sqs_queue_policy" "policy_statement_allow_principal_with_wildcard" {
  queue_url = aws_sqs_queue.test_queue.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sqs:SendMessage",
      "Principal": "*",
      "Resource": "${aws_sqs_queue.test_queue.arn}"
    }
  ]
}
EOF
}
