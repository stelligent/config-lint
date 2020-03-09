## Setup Helper
resource "aws_iam_role" "test_role" {
  assume_role_policy = <<-EOF
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Principal": {
          "Service": "ec2.amazonaws.com"
        },
        "Effect": "Allow",
        "Sid": ""
      }
    ]
  }
  EOF
}

# Pass
resource "aws_iam_role_policy" "policy_statement_allow_resource_without_wildcard" {
  role = "${aws_iam_role.test_role.id}"

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
resource "aws_iam_role_policy" "policy_statement_allow_resource_with_wildcard" {
  role = "${aws_iam_role.test_role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "rds:Add*",
      "Effect": "Allow",
      "Resource": "arn:aws:rds:*:*:cluster:*"
    }
  ]
}

EOF
}

# Pass
resource "aws_iam_role_policy" "policy_statement_deny_resource_with_wildcard" {
  role = "${aws_iam_role.test_role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "rds:Add*",
      "Effect": "Deny",
      "Resource": "arn:aws:rds:*:*:cluster:*"
    }
  ]
}

EOF
}
