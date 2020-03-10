## Setup Helper
resource "aws_s3_bucket" "test_bucket" {
}

# Pass
resource "aws_s3_bucket_policy" "policy_statement_allow_condition_SecureTransport_set_to_true" {
  bucket = aws_s3_bucket.test_bucket.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "s3:GetObject",
      "Resource": [
        "arn:aws:s3:::fooBucket/*"
      ],
      "Condition": {
        "Bool": {
          "aws:SecureTransport": "true"
        }
      },
      "Principal": "*"
    }
  ]
}
EOF
}

# Pass
resource "aws_s3_bucket_policy" "policy_statement_deny_condition_SecureTransport_set_to_true" {
  bucket = aws_s3_bucket.test_bucket.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Deny",
      "Action": "s3:GetObject",
      "Resource": [
        "arn:aws:s3:::fooBucket/*"
      ],
      "Condition": {
        "Bool": {
          "aws:SecureTransport": "true"
        }
      },
      "Principal": "*"
    }
  ]
}
EOF
}

# Pass
resource "aws_s3_bucket_policy" "policy_statement_deny_condition_SecureTransport_set_to_false" {
  bucket = aws_s3_bucket.test_bucket.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Deny",
      "Action": "s3:GetObject",
      "Resource": [
        "arn:aws:s3:::fooBucket/*"
      ],
      "Condition": {
        "Bool": {
          "aws:SecureTransport": "false"
        }
      },
      "Principal": "*"
    }
  ]
}
EOF
}

# Fail
resource "aws_s3_bucket_policy" "policy_statement_allow_condition_SecureTransport_set_to_false" {
  bucket = aws_s3_bucket.test_bucket.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "s3:GetObject",
      "Resource": [
        "arn:aws:s3:::fooBucket/*"
      ],
      "Condition": {
        "Bool": {
          "aws:SecureTransport": "false"
        }
      },
      "Principal": "*"
    }
  ]
}
EOF
}
