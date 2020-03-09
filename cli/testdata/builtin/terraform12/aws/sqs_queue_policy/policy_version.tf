## Setup Helper
resource "aws_sqs_queue" "test_queue" {
  name = "examplequeue"
  arn = "mocked_arn"
}

# Pass
resource "aws_sqs_queue_policy" "policy_version_set_correctly" {
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

# Fail
resource "aws_sqs_queue_policy" "policy_version_set_incorrectly" {
  queue_url = aws_sqs_queue.test_queue.id

  policy = <<EOF
{
  "Version": "2008-10-17",
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
