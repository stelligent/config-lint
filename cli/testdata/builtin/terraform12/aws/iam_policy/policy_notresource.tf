# Pass
resource "aws_iam_policy" "policy_statement_without_notresource" {
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

# Warn
resource "aws_iam_policy" "policy_statement_with_notresource" {
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Deny",
      "Action": "s3:GetObject",
      "NotResource": [
        "arn:aws:s3:::HRBucket/Payroll",
        "arn:aws:s3:::HRBucket/Payroll/*"
      ]
    }
  ]
}
EOF
}
