resource "aws_kms_key" "test_key" {
  enable_key_rotation = true
}

resource "aws_s3_bucket" "b1" {
  bucket = "test-bucket-1"

  #   server_side_encryption_configuration {
  #     rule {
  #       apply_server_side_encryption_by_default {
  #         kms_master_key_id = aws_kms_key.test_key.arn
  #         sse_algorithm     = "aws:kms"
  #       }
  #     }
  #   }

  tags = {
    "Department" = "invalid"
  }
}

resource "aws_s3_bucket_policy" "b1" {
  bucket = aws_s3_bucket.b1.id

  policy =<<POLICY
{
  "Version": "2012-10-17",
  "Id": "MYBUCKETPOLICY",
  "Statement": [
    {
      "Sid": "IPAllow",
      "Effect": "Deny",
      "Principal": "*",
      "Action": "s3:*",
      "Resource": "arn:aws:s3:::my_tf_test_bucket/*",
      "Condition": {
         "IpAddress": {"aws:SourceIp": "8.8.8.8/32"}
      }
    }
  ]
}
POLICY
}
