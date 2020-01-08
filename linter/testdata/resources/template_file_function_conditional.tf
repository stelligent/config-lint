terraform {
  required_version = ">= 0.12.0"
}

variable "var1" {
  default = "Alpha"
}

variable "var2" {
  default = "Bravo"
}

resource "test_resource" "test" {
  test_value =  templatefile("./testdata/data/template_file_example_conditional", { test_var = var.var1})
}

resource "test_resource" "test2" {
  test_value2 =  templatefile("./testdata/data/template_file_example_conditional", { test_var = var.var2})
}