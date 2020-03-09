# Pass
resource "aws_iam_role" "assume_role_policy_version_set_correctly" {
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
resource "aws_iam_role" "assume_role_policy_version_set_incorrectly" {
  assume_role_policy = <<EOF
{
  "Version": "2008-10-17",
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
