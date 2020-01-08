//resource "aws_s3_bucket" "a_bucket" {
//  bucket = "${file("./testdata/data/multi_line_content")}"
//}

terraform {
  required_version = ">= 0.12.0"
}

//data "template_file" "multi_line" {
//  template = file("./testdata/data/multi_line_content")
//}

resource "test_resource" "test" {
  test_value = file("./testdata/data/multi_line_content")
}

resource "test_resource2" "test2" {
  test_value2 = <<EOT
multi
line
example
EOT
}
