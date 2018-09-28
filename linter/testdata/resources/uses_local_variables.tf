locals {
    bucket_name = "myprojectbucket"
}

resource "aws_s3_bucket" "my_bucket" {
    name = "${local.bucket_name}"
}
