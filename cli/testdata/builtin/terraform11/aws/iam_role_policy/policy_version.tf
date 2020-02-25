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
resource "aws_iam_role_policy" "policy_version_set_correctly" {
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
resource "aws_iam_role_policy" "policy_version_set_incorrectly" {
  role = "${aws_iam_role.test_role.id}"

  policy = <<EOF
{
  "Version": "2008-10-17",
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
