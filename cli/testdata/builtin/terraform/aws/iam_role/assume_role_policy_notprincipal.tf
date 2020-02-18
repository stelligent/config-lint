# Pass
resource "aws_iam_role" "assume_role_policy_statement_without_NotPrincipal" {
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      }
    }
  ]
}
EOF
}

# Warn
resource "aws_iam_role" "assume_role_policy_statement_with_NotPrincipal" {
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Deny",
      "NotPrincipal": {
        "AWS": [
          "arn:aws:iam::1234567890:user/foo"
        ]
      },
      "Action": "s3:*",
      "Resource": [
        "arn:aws:s3:::fooBucket",
        "arn:aws:s3:::fooBucket/*"
      ]
    }
  ]
}
EOF
}
