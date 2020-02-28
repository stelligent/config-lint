resource "aws_s3_bucket" "b" {
  bucket = "a_test_bucket"
}

resource "aws_s3_bucket_policy" "b" {
  bucket = "${aws_s3_bucket.b.id}"

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "MYBUCKETPOLICY",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": "user",
      "Action": "s3:GetBucket",
      "Resource": "arn:aws:s3:::a_test_bucket/*",
      "Condition": {
        "Bool": {
          "aws:SecureTransport": "true"
        }
      }
    }
  ]
}
POLICY
}
