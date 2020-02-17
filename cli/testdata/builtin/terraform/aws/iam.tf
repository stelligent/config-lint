resource "aws_iam_policy" "policy_ok" {
    name = "policy_ok"
    path = "/"
    policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:ListObjects"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:s3:::my_corporate_bucket/exampleobject.png"
    }
  ]
}
POLICY
}

resource "aws_iam_policy" "policy_with_wildcard_resource" {
    name = "policy_with_wildcard_resource"
    path = "/"
    policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:DescribeInstance"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
POLICY
}

resource "aws_iam_policy" "policy_with_wildcard_action" {
    name = "policy_with_wildcard_action"
    path = "/"
    policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:*"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:ec2:us-west-2:123456789012:instance/i-123"
    }
  ]
}
POLICY
}

resource "aws_iam_user" "user1" {
    name = "lint"
}

resource "aws_iam_user_policy_attachment" "attachment1" {
    user = "${aws_iam_user.user1.id}"
    policy_arn = "${aws_iam_policy.policy1.arn}"
}

resource "aws_iam_group" "group1" {
    name = "test"
}

resource "aws_iam_group_membership" "membership1" {
    name = "test_membership"
    group = "${aws_iam_group.group1.id}"
    users = [ "${aws_iam_user.user1.id}" ]
}

resource "aws_iam_role" "role1" {
    name = "test_role"
    assume_role_policy = <<POLICY
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
POLICY
}

resource "aws_iam_role_policy" "role_policy" {
    role = "${aws_iam_role.role1.id}"
    policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:ListObjects"
      ],
      "Effect": "Allow",
      "Resource": "arn::aws::blah::blah"
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy" "role_policy_old_version" {
    role = "${aws_iam_role.role1.id}"
    policy = <<POLICY
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Action": [
        "s3:ListObjects"
      ],
      "Effect": "Allow",
      "Resource": "arn::aws::blah::blah"
    }
  ]
}
POLICY
}

resource "aws_iam_role" "role_with_old_version" {
    name = "test_role"
    assume_role_policy = <<POLICY
{
  "Version": "2008-10-17",
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
POLICY
}
