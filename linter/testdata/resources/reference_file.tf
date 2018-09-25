resource "aws_s3_bucket" "a_bucket" {
  bucket = "${file("./testdata/data/bucket_name")}"
}
