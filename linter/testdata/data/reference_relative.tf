resource "aws_s3_bucket" "a_bucket" {
  bucket = "${file("bucket_name")}"
}
