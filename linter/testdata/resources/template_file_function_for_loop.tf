terraform {
  required_version = ">= 0.12.0"
}

variable "words" {
  default = ["foo", "bar"]
}

resource "test" "test_resource" {
  test_value = templatefile("./testdata/data/template_file_example_for_loop", { words = var.words})
}