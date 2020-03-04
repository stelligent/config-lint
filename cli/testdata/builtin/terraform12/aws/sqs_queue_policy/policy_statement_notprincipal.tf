## Setup Helper
resource "aws_sqs_queue" "test_queue" {
  name                              = "test_queue"
  kms_master_key_id                 = "alias/foo/bar"
  kms_data_key_reuse_period_seconds = 60
}

# Pass
resource "aws_sqs_queue_policy" "policy_statement_without_notprincipal" {
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
      "Resource": ${aws_sqs_queue.test_queue.arn}
    }
  ]
}
EOF
}

# Warn
resource "aws_sqs_queue_policy" "policy_statement_with_notprincipal" {
  queue_url = aws_sqs_queue.test_queue.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sqs:SendMessage",
      "NotPrincipal": {
        "AWS": [
          "arn:aws:iam::1234567890:user/foo"
        ]
      },
      "Resource": ${aws_sqs_queue.test_queue.arn}
    }
  ]
}
EOF
}
