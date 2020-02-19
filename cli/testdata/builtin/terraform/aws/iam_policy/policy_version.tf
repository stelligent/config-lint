# Pass
resource "aws_iam_policy" "policy_version_set_correctly" {
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:Describe*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

# Fail
resource "aws_iam_policy" "policy_version_set_incorrectly" {
  policy = <<EOF
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:Describe*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}
