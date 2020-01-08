terraform {
  required_version = ">= 0.12.0"
}

variable "object_1" {
  default = "foo"
}

variable "object_2" {
  default = "bar"
}

resource "aws_s3_bucket" "a_bucket" {
  bucket = templatefile("./testdata/data/template_file_example_basic", { var1 = var.object_1, var2 = var.object_2})
}