resource "aws_iam_role" "test_role" {
    name = "role1"
    assume_role_policy = <<EOF
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

resource "aws_iam_role" "test_role_with_invalid_policy" {
    name = "role1"
    assume_role_policy = [ "foo" ]
}
