## Setup Helper
resource "aws_iam_user" "test_user" {
  name = "test-user"
}

resource "aws_iam_policy" "test_policy" {
  name   = "test-policy"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "rds:AddRoleToDBCluster",
      "Effect": "Allow",
      "Resource": "arn:aws:rds:us-east-1:1234567890:cluster:foo_cluster"
    }
  ]
}
EOF
}

# Fail
resource "aws_iam_user_policy_attachment" "resource_exists" {
  user       = "${aws_iam_user.user.name}"
  policy_arn = "${aws_iam_policy.policy.arn}"
}
