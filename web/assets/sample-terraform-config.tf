resource "aws_s3_bucket" "bucket_example_1" {
  bucket = "my-bucket-1"
  acl = "public-read"
}

resource "aws_s3_bucket" "bucket_example_2" {
  bucket = "my-bucket-2"
  acl = "public-read-write"
  encrypted = false
}
