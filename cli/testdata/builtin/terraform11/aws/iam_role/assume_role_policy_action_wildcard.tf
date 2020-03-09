# Pass
resource "aws_iam_role" "assume_role_policy_statement_allow_action_without_wildcard" {
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "rds:DeleteDBSnapshot",
      "Resource": "arn:aws:rds:*:*:snapshot:*"
    }
  ]
}
EOF
}

# Fail
resource "aws_iam_role" "assume_role_policy_statement_allow_action_with_wildcard" {
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "rds:*",
      "Resource": "*"
    }
  ]
}
EOF
}
