## Setup Helper
resource "aws_iam_user" "test_user" {
  name = "foobar"
}

# Fail
resource "aws_iam_user_policy" "resource_exists" {
  user = "${aws_iam_user.test_user.name}"

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
