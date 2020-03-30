## Setup Helper
resource "aws_s3_bucket" "test_bucket" {
}

# Pass
resource "aws_s3_bucket_policy" "policy_statement_without_notaction" {
  bucket = aws_s3_bucket.test_bucket.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::fooBucket/*"
    }
  ]
}
EOF
}

# Fail
resource "aws_s3_bucket_policy" "policy_statement_with_notaction" {
  bucket = aws_s3_bucket.test_bucket.id

  policy = <<EOF
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "NotAction": "s3:GetObject",
      "Resource": "arn:aws:s3:::fooBucket/*"
    }
  ]
}
EOF
}
