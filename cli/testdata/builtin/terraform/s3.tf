resource "aws_s3_bucket" "bucket1" {
}

resource "aws_s3_bucket_policy" "policy1" {
    bucket = "${aws_s3_bucket.bucket1.bucket}"
    policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "MYBUCKETPOLICY",
  "Statement": [
    {
      "Sid": "IPAllow",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:*",
      "Resource": "arn:aws:s3:::${aws_s3_bucket.bucket1.bucket}/*",
      "Condition": {
         "IpAddress": {"aws:SourceIp": "10.10.1.10/32"}
      } 
    } 
  ]
}
POLICY
}

resource "aws_s3_bucket_object" "object1" {
    bucket = "${aws_s3_bucket.bucket1.bucket}"
    key = "index.html"
    content = "Hello, world"
    content_type = "text/html"
}
